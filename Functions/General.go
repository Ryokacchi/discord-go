package functions

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func FormatBytes(b uint64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	bytes := float64(b)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", bytes/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", bytes/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", bytes/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", bytes/KB)
	default:
		return fmt.Sprintf("%.2f B", bytes)
	}
}

func GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		return 0, err
	}

	if len(percentages) > 0 {
		return percentages[0], nil
	}

	return 0, nil
}
