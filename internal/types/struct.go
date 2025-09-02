// internal/types/struct.go

package types

type Config struct {
	ContainerName         string
	GrafanaInfluxURL      string
	GrafanaAPIKey         string
	GrafanaUsername       string
	ReadTickerTimeSec     uint16
	FlushTickerTimeSec    uint16
	SleepBetweenFetchesMs uint16
	MetricKeys            []string
}

type Stats struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Group        string  `json:"group"`
	CPUMaxPerc   float32 `json:"cpu_max_perc"`
	CPUMinPerc   float32 `json:"cpu_min_perc"`
	CPUAvgPerc   float32 `json:"cpu_avg_perc"`
	CPUPercSum   float32 `json:"-"` // Used for calculating average
	CPUCount     int     `json:"-"` // Used for calculating average
	MemMaxMB     float32 `json:"mem_max_mb"`
	MemMinMB     float32 `json:"mem_min_mb"`
	MemAvgMB     float32 `json:"mem_avg_mb"`
	MemMBPercSum float32 `json:"-"` // Used for calculating average
	MemCount     int     `json:"-"` // Used for calculating average
	MemMaxPerc   float32 `json:"mem_max_perc"`
	MemMinPerc   float32 `json:"mem_min_perc"`
	MemAvgPerc   float32 `json:"mem_avg_perc"`
	MemPercSum   float32 `json:"-"` // Used for calculating average
	MemPercCount int     `json:"-"` // Used for calculating average
	DiskMB       float32 `json:"disk_mb"`
}
