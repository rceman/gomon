package app

import "time"

type dataPoint struct {
	ts  time.Time
	cpu float32
	mem uint32
}

type Stats struct {
	CPU [4]float32 `json:"cpu"`
	Mem [4]uint32  `json:"mem"`
}

type NodeStats struct {
	Name string `json:"name"`
	Stats
}

func computeSummary(points []dataPoint, now time.Time) Stats {
	var s Stats
	for _, p := range points {
		delta := now.Sub(p.ts)
		if delta <= time.Minute {
			if p.cpu > s.CPU[0] {
				s.CPU[0] = p.cpu
			}
			if p.mem > s.Mem[0] {
				s.Mem[0] = p.mem
			}
		}
		if delta <= 5*time.Minute {
			if p.cpu > s.CPU[1] {
				s.CPU[1] = p.cpu
			}
			if p.mem > s.Mem[1] {
				s.Mem[1] = p.mem
			}
		}
		if delta <= time.Hour {
			if p.cpu > s.CPU[2] {
				s.CPU[2] = p.cpu
			}
			if p.mem > s.Mem[2] {
				s.Mem[2] = p.mem
			}
		}
		if delta <= 24*time.Hour {
			if p.cpu > s.CPU[3] {
				s.CPU[3] = p.cpu
			}
			if p.mem > s.Mem[3] {
				s.Mem[3] = p.mem
			}
		}
	}
	return s
}
