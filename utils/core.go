package utils

import (
	"runtime"
	"time"

	"log"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"

	"gopkg.in/ini.v1"
)

var (
	Uptime         time.Time = time.Now()
	ActiveHandlers int       = 0
	configPath               = "config.ini"
)

type MemoryStats struct {
	Runtime uint64
	System  *mem.VirtualMemoryStat
}

// configLoader loads the configuration from config.ini file.
// If the config file cannot be loaded, it logs the error and exits the program.
func ConfigLoader() *ini.File {
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config.ini file: %v", err)
	}

	return cfg
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
