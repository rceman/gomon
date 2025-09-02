// internal/stats/docker/stats.go

package docker

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/therceman/gomon/internal/helpers"
)

type Stats struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	CPU     float32 `json:"cpu"`
	MemMB   float32 `json:"mem"`
	MemPerc float32 `json:"mem_perc"`
	NetI    float32 `json:"net_i"`
	NetO    float32 `json:"net_o"`
	BlockI  float32 `json:"block_i"`
	BlockO  float32 `json:"block_o"`
	PIDs    int     `json:"pids"`
	SizeMB  float32 `json:"size"`
}

func GetStats() ([]Stats, error) {
	cmd := exec.Command("docker", "stats", "--no-stream")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	output := out.String()

	stats, err := parseDockerStatsOutput(output)
	if err != nil {
		return nil, err
	}

	// Ensure we do not hold on to memory longer than needed
	out.Reset()
	return stats, nil
}

func parseDockerStatsOutput(output string) ([]Stats, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var stats []Stats
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "CONTAINER ID") {
			continue // Skip the header line
		}

		fields := strings.Fields(line)
		if len(fields) < 14 {
			log.Printf("Skipping incomplete line: %s", line)
			continue // Skip incomplete lines
		}

		netI, err := helpers.ConvertSizeToMB(fields[7])
		if err != nil {
			log.Printf("Error converting Net I: %v", err)
			continue
		}

		netO, err := helpers.ConvertSizeToMB(fields[9])
		if err != nil {
			log.Printf("Error converting Net O: %v", err)
			continue
		}

		blockI, err := helpers.ConvertSizeToMB(fields[10])
		if err != nil {
			log.Printf("Error converting Block I: %v", err)
			continue
		}

		blockO, err := helpers.ConvertSizeToMB(fields[12])
		if err != nil {
			log.Printf("Error converting Block O: %v", err)
			continue
		}

		cpuUsage, err := helpers.ConvertToPerc(fields[2])
		if err != nil {
			log.Printf("Error parsing CPUPerc usage: %v", err)
			continue
		}

		memUsage, err := helpers.ConvertMemoryToMB(fields[3])
		if err != nil {
			log.Printf("Error parsing memory usage: %v", err)
			continue
		}

		memPerc, err := helpers.ConvertToPerc(fields[6])
		if err != nil {
			log.Printf("Error parsing memory percent: %v", err)
			continue
		}

		pids, err := strconv.Atoi(fields[13])
		if err != nil {
			log.Printf("Error parsing PIDs: %v", err)
			continue
		}

		containerSize, err := GetContainerSize(fields[0], true)
		if err != nil {
			log.Printf("Error getting docker container size: %v", err)
			continue
		}

		stat := Stats{
			ID:      fields[0],
			Name:    fields[1],
			CPU:     helpers.RoundToTwoDecimal(cpuUsage),
			MemMB:   helpers.RoundToTwoDecimal(memUsage),
			MemPerc: helpers.RoundToTwoDecimal(memPerc),
			NetI:    helpers.RoundToTwoDecimal(netI),
			NetO:    helpers.RoundToTwoDecimal(netO),
			BlockI:  helpers.RoundToTwoDecimal(blockI),
			BlockO:  helpers.RoundToTwoDecimal(blockO),
			PIDs:    pids,
			SizeMB:  helpers.RoundToTwoDecimal(containerSize),
		}
		stats = append(stats, stat)
	}
	return stats, nil
}
