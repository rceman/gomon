package app

import "testing"

func TestFormatStatsDefault(t *testing.T) {
	statsMap := map[string]Stats{
		"api": {
			CPU:  [4]float32{1, 2, 3, 4},
			Mem:  [4]uint32{10, 20, 30, 40},
			Disk: [4]float32{0.1, 0.2, 0.3, 0.4},
		},
	}
	got := formatStatsDefault(statsMap)
	want := map[string]styledStats{
		"api": {
			CPU:  map[string]float32{"1m": 1, "5m": 2, "1h": 3, "24h": 4},
			Mem:  map[string]uint32{"1m": 10, "5m": 20, "1h": 30, "24h": 40},
			Disk: map[string]float32{"1m": 0.1, "5m": 0.2, "1h": 0.3, "24h": 0.4},
		},
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d entries, got %d", len(want), len(got))
	}
	g := got["api"]
	w := want["api"]
	for k, v := range w.CPU {
		if g.CPU[k] != v {
			t.Fatalf("cpu %s: expected %v got %v", k, v, g.CPU[k])
		}
	}
	for k, v := range w.Mem {
		if g.Mem[k] != v {
			t.Fatalf("mem %s: expected %v got %v", k, v, g.Mem[k])
		}
	}
	for k, v := range w.Disk {
		if g.Disk[k] != v {
			t.Fatalf("disk %s: expected %v got %v", k, v, g.Disk[k])
		}
	}
}
