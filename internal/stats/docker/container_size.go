// internal/stats/docker/container_size.go

package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/therceman/gomon/internal/helpers"
)

type ContainerInfo struct {
	ID   string `json:"ID"`
	Size string `json:"Size"`
}

func GetContainerSize(containerID string, virtual bool) (float32, error) {
	cmd := exec.Command("docker", "ps", "--format", "{{json .}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("error executing docker ps --format='{{json .}}' command: %v", err)
	}

	output := strings.TrimSpace(out.String())
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "'") // Remove single quotes
		var containerInfo ContainerInfo
		if err := json.Unmarshal([]byte(line), &containerInfo); err != nil {
			return 0, fmt.Errorf("error parsing JSON from docker ps: %v", err)
		}

		if containerInfo.ID == containerID {
			sizeStr := containerInfo.Size
			if virtual {
				re := regexp.MustCompile(`(virtual ([^)]+))`)
				matches := re.FindStringSubmatch(sizeStr)
				if len(matches) > 2 {
					sizeStr = matches[2]
				} else {
					return 0, fmt.Errorf("virtual size not found for container %s", containerID)
				}
			} else {
				parts := strings.Fields(sizeStr)
				if len(parts) > 0 {
					sizeStr = parts[0]
				} else {
					return 0, fmt.Errorf("size not found for container %s", containerID)
				}
			}

			size, err := helpers.ConvertSizeToMB(sizeStr)
			if err != nil {
				return 0, err
			}

			return size, nil
		}
	}

	return 0, fmt.Errorf("container %s not found", containerID)
}
