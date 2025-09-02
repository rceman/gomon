// internal/app/app.go

package app

import (
	"log"
	"runtime"
	"time"

	"github.com/therceman/gomon/internal/helpers"
	"github.com/therceman/gomon/internal/stats"
	"github.com/therceman/gomon/internal/types"
)

func Run(config types.Config) {
	log.Println("Running Go Monitor for Container:", config.ContainerName)
	log.Println("Metric Keys:", config.MetricKeys)
	log.Printf("Read Ticker Time: %ds, Flush Ticker Time: %ds, Sleep Between Fetches: %dms",
		config.ReadTickerTimeSec, config.FlushTickerTimeSec, config.SleepBetweenFetchesMs,
	)

	pid := helpers.GetCurrentPID()
	pidStr := helpers.ConvertUint32ToString(pid)

	readTickerTime := time.Duration(config.ReadTickerTimeSec) * time.Second
	flushTickerTime := time.Duration(config.FlushTickerTimeSec) * time.Second

	ticker := time.NewTicker(readTickerTime)
	flushTicker := time.NewTicker(flushTickerTime)

	defer ticker.Stop()
	defer flushTicker.Stop()

	statsMap := make(map[string]*types.Stats)

	for {
		select {
		case <-ticker.C:
			systemFetchError := stats.FetchSystemStats(statsMap)
			if systemFetchError != nil {
				log.Printf("Error fetching system stats: %v", systemFetchError)
			}
			time.Sleep(time.Millisecond * 250)
			dockerFetchError := stats.FetchDockerStats(statsMap)
			if dockerFetchError != nil {
				log.Printf("Error fetching docker stats: %v", dockerFetchError)
			}
			time.Sleep(time.Millisecond * 250)
			workerFetchError := stats.FetchWorkerStats(statsMap, pidStr, pid, "gomon")
			if workerFetchError != nil {
				log.Printf("Error fetching worker stats: %v", workerFetchError)
			}
			runtime.GC()
		case <-flushTicker.C:
			stats.FlushStats(statsMap, config)
			statsMap = make(map[string]*types.Stats)
			runtime.GC()
		}
	}
}
