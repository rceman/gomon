// internal/stats/worker/stats.go

package worker

import (
	"bufio"
	"fmt"
	"github.com/therceman/gomon/internal/helpers"
	"os/exec"
	"strings"
)

type Stats struct {
	MemKB   uint32  `json:"mem_kb"`   // Used memory in KB
	CPUPerc float32 `json:"cpu_perc"` // CPU usage percentage
	MemPerc float32 `json:"mem_perc"` // Memory usage percentage
	PID     uint32  `json:"pid"`      // Process ID
}

func GetStats(pidStr string, pid uint32) (Stats, error) {
	// Get the CPU and memory usage of the process by PID
	psCmd := exec.Command("ps", "-p", pidStr, "-o", "pid,pcpu,pmem,rss")
	psOutput, err := psCmd.Output()
	if err != nil {
		return Stats{}, fmt.Errorf("error executing ps command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(psOutput)))
	// Skip the header line
	scanner.Scan()

	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			return Stats{}, fmt.Errorf("unexpected format in ps output")
		}

		cpuPerc, err := helpers.ConvertStringToFloat32(fields[1])
		if err != nil {
			return Stats{}, err
		}

		memPerc, err := helpers.ConvertStringToFloat32(fields[2])
		if err != nil {
			return Stats{}, err
		}

		memKB, err := helpers.ConvertStringToUint32(fields[3])
		if err != nil {
			return Stats{}, err
		}

		return Stats{
			PID:     pid,
			CPUPerc: cpuPerc,
			MemPerc: memPerc,
			MemKB:   memKB,
		}, nil
	}

	if err := scanner.Err(); err != nil {
		return Stats{}, err
	}

	return Stats{}, fmt.Errorf("failed to read process stats")
}
