package app

// styledStats represents the default JSON output structure.
type styledStats struct {
	CPU  map[string]float32 `json:"cpu"`
	Mem  map[string]uint32  `json:"mem"`
	Disk map[string]float32 `json:"disk"`
}

var windows = []string{"1m", "5m", "1h", "24h"}

// formatStatsDefault converts a map of Stats into the default JSON format
// where the first key is the VM name and each metric is keyed by time window.
func formatStatsDefault(statsMap map[string]Stats) map[string]styledStats {
	out := make(map[string]styledStats, len(statsMap))
	for name, s := range statsMap {
		cpu := make(map[string]float32, len(windows))
		mem := make(map[string]uint32, len(windows))
		disk := make(map[string]float32, len(windows))
		for i, w := range windows {
			cpu[w] = s.CPU[i]
			mem[w] = s.Mem[i]
			disk[w] = s.Disk[i]
		}
		out[name] = styledStats{CPU: cpu, Mem: mem, Disk: disk}
	}
	return out
}
