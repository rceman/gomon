package app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/therceman/gomon/internal/stats/system"
	"github.com/therceman/gomon/internal/types"
)

func Run(config types.Config) {
	log.Println("Running Go Monitor")
	log.Printf("Read Ticker Time: %.3fs", config.ReadTickerTimeSec)

	var (
		mu          sync.RWMutex
		history     []dataPoint
		remoteStats = make(map[string]Stats)
	)

	go func() {
		ticker := time.NewTicker(time.Duration(float64(config.ReadTickerTimeSec) * float64(time.Second)))
		defer ticker.Stop()

		sendInterval := time.Duration(float64(config.MasterSendIntervalMin) * float64(time.Minute))
		var lastSend time.Time

		for range ticker.C {
			stats, err := system.GetStats()
			if err != nil {
				log.Printf("Error fetching system stats: %v", err)
				continue
			}
			mu.Lock()
			now := time.Now()
			history = append(history, dataPoint{ts: now, cpu: stats.CPUPerc, mem: stats.MemMB, disk: float32(stats.DiskMB) / 1024})
			cutoff := now.Add(-24 * time.Hour)
			for len(history) > 0 && history[0].ts.Before(cutoff) {
				history = history[1:]
			}
			mu.Unlock()

			if config.MasterSend && time.Since(lastSend) >= sendInterval {
				mu.RLock()
				s := computeSummary(history, time.Now())
				mu.RUnlock()
				sendToMaster(config, NodeStats{Name: config.Name, Stats: s})
				lastSend = time.Now()
			}
			runtime.GC()
		}
	}()

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		outputFormat := r.URL.Query().Get("output_format")
		outputStyle := r.URL.Query().Get("output_style")
		if outputStyle == "" {
			outputStyle = "default"
		}
		if outputFormat == "html" {
			w.Header().Set("Content-Type", "text/html")
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
		switch r.Method {
		case http.MethodPost:
			if !config.MasterNode {
				http.NotFound(w, r)
				return
			}
			if !checkAuth(r.Header.Get("Authorization"), config.MasterKey) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			var ns NodeStats
			if err := json.NewDecoder(r.Body).Decode(&ns); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			mu.Lock()
			remoteStats[ns.Name] = ns.Stats
			mu.Unlock()
			return
		default:
			mu.RLock()
			local := computeSummary(history, time.Now())
			if config.MasterNode {
				statsMap := make(map[string]Stats, len(remoteStats)+1)
				for k, v := range remoteStats {
					statsMap[k] = v
				}
				mu.RUnlock()
				statsMap[config.Name] = local
				if outputFormat == "html" {
					renderHTMLStats(w, statsMap)
				} else {
					if outputStyle == "short" {
						if err := json.NewEncoder(w).Encode(statsMap); err != nil {
							log.Printf("error encoding stats: %v", err)
						}
					} else {
						out := formatStatsDefault(statsMap)
						if err := json.NewEncoder(w).Encode(out); err != nil {
							log.Printf("error encoding stats: %v", err)
						}
					}
				}
				return
			}
			mu.RUnlock()
			if outputFormat == "html" {
				renderHTMLStats(w, map[string]Stats{config.Name: local})
			} else {
				if outputStyle == "short" {
					if err := json.NewEncoder(w).Encode(local); err != nil {
						log.Printf("error encoding stats: %v", err)
					}
				} else {
					out := formatStatsDefault(map[string]Stats{config.Name: local})
					if err := json.NewEncoder(w).Encode(out); err != nil {
						log.Printf("error encoding stats: %v", err)
					}
				}
			}
		}
	})
	addr := fmt.Sprintf(":%d", config.StatsPort)
	log.Printf("Stats endpoint available at %s/stats", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}

func checkAuth(header, key string) bool {
	const prefix = "Basic "
	if !strings.HasPrefix(header, prefix) {
		return false
	}
	payload, err := base64.StdEncoding.DecodeString(header[len(prefix):])
	if err != nil {
		return false
	}
	return string(payload) == key
}

func sendToMaster(cfg types.Config, ns NodeStats) {
	if cfg.MasterIP == "" || cfg.MasterPort == 0 || cfg.MasterKey == "" {
		return
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(ns); err != nil {
		log.Printf("error encoding stats for master: %v", err)
		return
	}
	url := fmt.Sprintf("http://%s:%d/stats", cfg.MasterIP, cfg.MasterPort)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}
	auth := base64.StdEncoding.EncodeToString([]byte(cfg.MasterKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending stats to master: %v", err)
		return
	}
	resp.Body.Close()
}

func renderHTMLStats(w http.ResponseWriter, statsMap map[string]Stats) {
	fmt.Fprintln(w, "<html><body>")
	for name, s := range statsMap {
		fmt.Fprintf(w, "<h2>%s</h2>", name)
		fmt.Fprintln(w, "<table border=\"1\">")
		fmt.Fprintln(w, "<tr><th></th><th>1m</th><th>5m</th><th>1h</th><th>24h</th></tr>")
		fmt.Fprintf(w, "<tr><td>RAM</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>", s.Mem[0], s.Mem[1], s.Mem[2], s.Mem[3])
		fmt.Fprintf(w, "<tr><td>CPU</td><td>%.2f</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>", s.CPU[0], s.CPU[1], s.CPU[2], s.CPU[3])
		fmt.Fprintf(w, "<tr><td>Disk (GB)</td><td>%.2f</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>", s.Disk[0], s.Disk[1], s.Disk[2], s.Disk[3])
		fmt.Fprintln(w, "</table>")
	}
	fmt.Fprintln(w, "</body></html>")
}
