// internal/sender/grafana/influx.go

package grafana

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/therceman/gomon/internal/types"
)

// PrepareInfluxData formats the stats data according to the specified metric keys.
func PrepareInfluxData(metricKeys []string, cont string, stats types.Stats) string {
	prefix := "gomon"

	var values []string
	for _, key := range metricKeys {
		switch key {
		case "cpu_max_perc":
			values = append(values, fmt.Sprintf("cpu_max_perc=%.2f", stats.CPUMaxPerc))
		case "cpu_avg_perc":
			values = append(values, fmt.Sprintf("cpu_avg_perc=%.2f", stats.CPUAvgPerc))
		case "mem_max_mb":
			values = append(values, fmt.Sprintf("mem_max_mb=%.2f", stats.MemMaxMB))
		case "mem_avg_mb":
			values = append(values, fmt.Sprintf("mem_avg_mb=%.2f", stats.MemAvgMB))
		case "mem_max_perc":
			values = append(values, fmt.Sprintf("mem_max_perc=%.2f", stats.MemMaxPerc))
		case "mem_avg_perc":
			values = append(values, fmt.Sprintf("mem_avg_perc=%.2f", stats.MemAvgPerc))
		case "disk_mb":
			values = append(values, fmt.Sprintf("disk_mb=%.2f", stats.DiskMB))
		}
	}

	dataLine := fmt.Sprintf(
		"%s,cont=%s,group=%s,id=%s,name=%s %s",
		prefix, cont, stats.Group, stats.ID, stats.Name, strings.Join(values, ","),
	)

	return dataLine
}

// SendToInflux sends the prepared data to InfluxDB.
func _(url string, username string, apiKey string, data string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain")
	req.SetBasicAuth(username, apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing resp.Body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send data to InfluxDB, status code: %d", resp.StatusCode)
	}

	log.Println("Metrics pushed successfully to Grafana->Influx")
	return nil
}
