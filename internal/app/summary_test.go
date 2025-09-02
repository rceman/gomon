package app

import (
	"reflect"
	"testing"
	"time"
)

func TestComputeSummary(t *testing.T) {
	now := time.Now()
	points := []dataPoint{
		{ts: now.Add(-2 * time.Minute), cpu: 10, mem: 100},
		{ts: now.Add(-30 * time.Minute), cpu: 50, mem: 300},
		{ts: now.Add(-2 * time.Hour), cpu: 90, mem: 400},
	}
	got := computeSummary(points, now)
	want := Stats{
		CPU: [4]float32{0, 10, 50, 90},
		Mem: [4]uint32{0, 100, 300, 400},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("computeSummary() = %v, want %v", got, want)
	}
}
