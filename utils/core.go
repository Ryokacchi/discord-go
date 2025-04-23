package utils

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
)

var (
	Uptime         time.Time = time.Now()
	ActiveHandlers int       = 0
)

type MemoryStats struct {
	Runtime uint64
	System  *mem.VirtualMemoryStat
}

// TrackHandler increments the ActiveHandlers counter and returns the original handler.
// It is useful for tracking how many handlers have been registered.
func TrackHandler(handler interface{}) interface{} {
	ActiveHandlers++
	return handler
}

func ReadCPUUsage() float64 {
	const sampleInterval = 500 * time.Millisecond
	usage, err := cpu.Percent(sampleInterval, false)

	if err != nil {
		return 0
	}

	if len(usage) == 0 {
		return 0
	}

	return usage[0]
}

func ReadMemoryStats() *MemoryStats {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &MemoryStats{
		Runtime: memStats.Alloc,
		System:  vm,
	}
}
