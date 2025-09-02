// internal/stats/system/stats.go

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/therceman/gomon/internal/helpers"
)

// Stats holds combined metrics for system resources
type Stats struct {
	MemMB    uint32  `json:"mem"`       // Used memory in MB
	MemPerc  float32 `json:"mem_perc"`  // Used memory percentage
	CPUPerc  float32 `json:"cpu_perc"`  // CPU usage percentage
	DiskMB   uint32  `json:"disk_mb"`   // Used disk space in MB
	DiskPerc float32 `json:"disk_perc"` // Used disk space percentage
}

// GetStats retrieves system statistics including memory, CPUPerc, and disk usage
func GetStats() (Stats, error) {
	memStats, err := getMemStats()
	if err != nil {
		return Stats{}, err
	}

	cpuStats, err := getCPUStats()
	if err != nil {
		return Stats{}, err
	}

	diskStats, err := getDiskStats("/")
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		MemMB:    memStats.Used,
		MemPerc:  memStats.UsedPercent,
		CPUPerc:  cpuStats.CPUUsagePercent,
		DiskMB:   diskStats.Used,
		DiskPerc: diskStats.UsedPerc,
	}, nil
}

type memStats struct {
	Used        uint32  `json:"used"`
	UsedPercent float32 `json:"used_percent"`
}

type cpuStats struct {
	CPUUsagePercent float32 `json:"cpu_usage_percent"`
}

type diskStats struct {
	Used     uint32  `json:"used"`
	UsedPerc float32 `json:"used_perc"`
}

func getMemStats() (memStats, error) {
	var result memStats

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return result, err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var total, free uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			total, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return result, err
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			free, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return result, err
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	used := total - free
	usedPercent := (float32(used) / float32(total)) * 100

	return memStats{
		Used:        uint32(used / 1024), // Convert KB to MB
		UsedPercent: helpers.RoundToTwoDecimal(usedPercent),
	}, nil
}

// getCPUStats retrieves CPUPerc usage percentage
func getCPUStats() (cpuStats, error) {
	idle0, total0, _, _, err := getCPUSample()
	if err != nil {
		return cpuStats{}, err
	}

	// Sleep for a short duration to get a new sample
	time.Sleep(200 * time.Millisecond)

	idle1, total1, _, _, err := getCPUSample()
	if err != nil {
		return cpuStats{}, err
	}

	idleTicks := float32(idle1 - idle0)
	totalTicks := float32(total1 - total0)

	// Avoid division by zero
	if totalTicks == 0 {
		return cpuStats{CPUUsagePercent: 0}, nil
	}

	cpuUsage := (1.0 - (idleTicks / totalTicks)) * 100

	return cpuStats{
		CPUUsagePercent: helpers.RoundToTwoDecimal(cpuUsage),
	}, nil
}

// getCPUSample reads CPUPerc usage statistics from /proc/stat
// and returns per-core idle and total times
func getCPUSample() (idle, total uint64, idleTimes, totalTimes []uint64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, nil, nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 && strings.HasPrefix(fields[0], "cpu") {
			idleTime, totalTime := uint64(0), uint64(0)
			for i, v := range fields[1:] {
				val, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return 0, 0, nil, nil, err
				}
				totalTime += val
				if i == 3 { // idle time is the 4th field
					idleTime = val
				}
			}
			if fields[0] == "cpu" {
				total += totalTime
				idle += idleTime
			} else {
				idleTimes = append(idleTimes, idleTime)
				totalTimes = append(totalTimes, totalTime)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, nil, nil, fmt.Errorf("error reading /proc/stat: %v", err)
	}

	return idle, total, idleTimes, totalTimes, nil
}

func getDiskStats(path string) (diskStats, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		return diskStats{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free
	usedPerc := (float32(used) / float32(total)) * 100

	return diskStats{
		Used:     uint32(used / 1024 / 1024), // Convert KB to MB
		UsedPerc: helpers.RoundToTwoDecimal(usedPerc),
	}, nil
}
