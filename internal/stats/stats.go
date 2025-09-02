// internal/stats/stats.go

package stats

import (
	"log"

	"github.com/therceman/gomon/internal/helpers"
	"github.com/therceman/gomon/internal/sender/grafana"
	"github.com/therceman/gomon/internal/stats/docker"
	"github.com/therceman/gomon/internal/stats/system"
	"github.com/therceman/gomon/internal/stats/worker"
	"github.com/therceman/gomon/internal/types"
)

func FlushStats(statsMap map[string]*types.Stats, config types.Config) {
	log.Println("Flushing stats map")

	for _, stat := range statsMap {
		data := grafana.PrepareInfluxData(config.MetricKeys, config.ContainerName, *stat)
		log.Printf("data: %v\n", data)

		//err := grafana.SendToInflux(types.GrafanaInfluxURL, types.GrafanaUsername, types.GrafanaAPIKey, data)
		//if err != nil {
		//	log.Printf("Error sending data to InfluxDB: %v", err)
		//} else {
		//	log.Println("Metrics pushed successfully to InfluxDB")
		//}
	}
}

// FetchDockerStats fetches and updates Docker stats
func FetchDockerStats(statsMap map[string]*types.Stats) error {
	dockerStats, err := docker.GetStats()
	if err != nil {
		return err
	}

	for _, stat := range dockerStats {
		if existing, found := statsMap[stat.ID]; found {
			// Update CPUPerc percentages
			if stat.CPU < existing.CPUMinPerc {
				existing.CPUMinPerc = stat.CPU
			}
			if stat.CPU > existing.CPUMaxPerc {
				existing.CPUMaxPerc = stat.CPU
			}

			// Update CPU average
			existing.CPUPercSum += stat.CPU
			existing.CPUCount++
			existing.CPUAvgPerc = helpers.RoundToTwoDecimal(existing.CPUPercSum / float32(existing.CPUCount))

			// Update Memory usage in MB
			if stat.MemMB < existing.MemMinMB {
				existing.MemMinMB = stat.MemMB
			}
			if stat.MemMB > existing.MemMaxMB {
				existing.MemMaxMB = stat.MemMB
			}

			// Update Memory average
			existing.MemMBPercSum += stat.MemMB
			existing.MemCount++
			existing.MemAvgMB = helpers.RoundToTwoDecimal(existing.MemMBPercSum / float32(existing.MemCount))

			// Update Memory percentages
			if stat.MemPerc < existing.MemMinPerc {
				existing.MemMinPerc = stat.MemPerc
			}
			if stat.MemPerc > existing.MemMaxPerc {
				existing.MemMaxPerc = stat.MemPerc
			}

			// Update Memory percentage average
			existing.MemPercSum += stat.MemPerc
			existing.MemPercCount++
			existing.MemAvgPerc = helpers.RoundToTwoDecimal(existing.MemPercSum / float32(existing.MemPercCount))

			// Update Disk usage
			existing.DiskMB = stat.SizeMB
		} else {
			statsMap[stat.ID] = &types.Stats{
				ID:           stat.ID,
				Name:         stat.Name,
				Group:        "docker",
				CPUMinPerc:   stat.CPU,
				CPUMaxPerc:   stat.CPU,
				CPUPercSum:   stat.CPU,
				CPUCount:     1,
				CPUAvgPerc:   stat.CPU,
				MemMinMB:     stat.MemMB,
				MemMaxMB:     stat.MemMB,
				MemMBPercSum: stat.MemMB,
				MemCount:     1,
				MemAvgMB:     stat.MemMB,
				MemMinPerc:   stat.MemPerc,
				MemMaxPerc:   stat.MemPerc,
				MemPercSum:   stat.MemPerc,
				MemPercCount: 1,
				MemAvgPerc:   stat.MemPerc,
				DiskMB:       stat.SizeMB,
			}
		}
	}

	return nil
}

// FetchSystemStats fetches and updates system stats
func FetchSystemStats(statsMap map[string]*types.Stats) error {
	sysStats, err := system.GetStats()
	if err != nil {
		return err
	}

	ID := "system"
	NAME := helpers.GetOperatingSystem()

	if existing, found := statsMap[ID]; found {
		// Update CPUPerc percentages
		if sysStats.CPUPerc < existing.CPUMinPerc {
			existing.CPUMinPerc = sysStats.CPUPerc
		}
		if sysStats.CPUPerc > existing.CPUMaxPerc {
			existing.CPUMaxPerc = sysStats.CPUPerc
		}

		// Update CPU average
		existing.CPUPercSum += sysStats.CPUPerc
		existing.CPUCount++
		existing.CPUAvgPerc = helpers.RoundToTwoDecimal(existing.CPUPercSum / float32(existing.CPUCount))

		// Update Memory usage in MB
		if float32(sysStats.MemMB) < existing.MemMinMB {
			existing.MemMinMB = float32(sysStats.MemMB)
		}
		if float32(sysStats.MemMB) > existing.MemMaxMB {
			existing.MemMaxMB = float32(sysStats.MemMB)
		}

		// Update Memory average
		existing.MemMBPercSum += float32(sysStats.MemMB)
		existing.MemCount++
		existing.MemAvgMB = helpers.RoundToTwoDecimal(existing.MemMBPercSum / float32(existing.MemCount))

		// Update Memory percentages
		if sysStats.MemPerc < existing.MemMinPerc {
			existing.MemMinPerc = sysStats.MemPerc
		}
		if sysStats.MemPerc > existing.MemMaxPerc {
			existing.MemMaxPerc = sysStats.MemPerc
		}

		// Update Memory percentage average
		existing.MemPercSum += sysStats.MemPerc
		existing.MemPercCount++
		existing.MemAvgPerc = helpers.RoundToTwoDecimal(existing.MemPercSum / float32(existing.MemPercCount))

		// Update Disk usage
		existing.DiskMB = float32(sysStats.DiskMB)
	} else {
		statsMap[ID] = &types.Stats{
			ID:           ID,
			Name:         NAME,
			Group:        "system",
			CPUMinPerc:   sysStats.CPUPerc,
			CPUMaxPerc:   sysStats.CPUPerc,
			CPUPercSum:   sysStats.CPUPerc,
			CPUCount:     1,
			CPUAvgPerc:   sysStats.CPUPerc,
			MemMinMB:     float32(sysStats.MemMB),
			MemMaxMB:     float32(sysStats.MemMB),
			MemMBPercSum: float32(sysStats.MemMB),
			MemCount:     1,
			MemAvgMB:     float32(sysStats.MemMB),
			MemMinPerc:   sysStats.MemPerc,
			MemMaxPerc:   sysStats.MemPerc,
			MemPercSum:   sysStats.MemPerc,
			MemPercCount: 1,
			MemAvgPerc:   sysStats.MemPerc,
			DiskMB:       float32(sysStats.DiskMB),
		}
	}

	return nil
}

// FetchWorkerStats fetches and updates worker stats
func FetchWorkerStats(statsMap map[string]*types.Stats, pidStr string, pid uint32, processName string) error {
	workerStats, err := worker.GetStats(pidStr, pid)
	if err != nil {
		return err
	}

	ID := pidStr
	NAME := processName

	if existing, found := statsMap[ID]; found {
		// Update CPU percentages
		if workerStats.CPUPerc < existing.CPUMinPerc {
			existing.CPUMinPerc = workerStats.CPUPerc
		}
		if workerStats.CPUPerc > existing.CPUMaxPerc {
			existing.CPUMaxPerc = workerStats.CPUPerc
		}

		// Update CPU average
		existing.CPUPercSum += workerStats.CPUPerc
		existing.CPUCount++
		existing.CPUAvgPerc = helpers.RoundToTwoDecimal(existing.CPUPercSum / float32(existing.CPUCount))

		// Update Memory usage in MB
		memMB := helpers.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024)
		if memMB < existing.MemMinMB {
			existing.MemMinMB = memMB
		}
		if memMB > existing.MemMaxMB {
			existing.MemMaxMB = memMB
		}

		// Update Memory average
		existing.MemMBPercSum += memMB
		existing.MemCount++
		existing.MemAvgMB = helpers.RoundToTwoDecimal(existing.MemMBPercSum / float32(existing.MemCount))

		// Update Memory percentages
		if workerStats.MemPerc < existing.MemMinPerc {
			existing.MemMinPerc = workerStats.MemPerc
		}
		if workerStats.MemPerc > existing.MemMaxPerc {
			existing.MemMaxPerc = workerStats.MemPerc
		}

		// Update Memory percentage average
		existing.MemPercSum += workerStats.MemPerc
		existing.MemPercCount++
		existing.MemAvgPerc = helpers.RoundToTwoDecimal(existing.MemPercSum / float32(existing.MemPercCount))
	} else {
		statsMap[ID] = &types.Stats{
			ID:           ID,
			Name:         NAME,
			Group:        "worker",
			CPUMinPerc:   workerStats.CPUPerc,
			CPUMaxPerc:   workerStats.CPUPerc,
			CPUPercSum:   workerStats.CPUPerc,
			CPUCount:     1,
			CPUAvgPerc:   workerStats.CPUPerc,
			MemMinMB:     helpers.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemMaxMB:     helpers.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemMBPercSum: helpers.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemCount:     1,
			MemAvgMB:     helpers.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemMinPerc:   workerStats.MemPerc,
			MemMaxPerc:   workerStats.MemPerc,
			MemPercSum:   workerStats.MemPerc,
			MemPercCount: 1,
			MemAvgPerc:   workerStats.MemPerc,
			DiskMB:       0,
		}
	}

	return nil
}
