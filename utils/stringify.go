package utils

import (
	"fmt"
	"strings"
)

var (
	FormatsTimeMappings = map[string]string{
		"years ago":   "yıl önce",
		"year ago":    "yıl önce",
		"months ago":  "ay önce",
		"month ago":   "ay önce",
		"days ago":    "gün önce",
		"day ago":     "gün önce",
		"hours ago":   "saat önce",
		"hour ago":    "saat önce",
		"minutes ago": "dakika önce",
		"minute ago":  "dakika önce",
		"seconds ago": "saniye önce",
		"second ago":  "saniye önce",
		" ago":        " önce",
	}
)

// StringPtr returns a pointer to the given string.
//
// Useful when you need to assign string literals or values
// to fields that expect *string, such as in API structs.
func StringPtr(s string) *string {
	return &s
}

func FormatTime(result string) string {
	for k, v := range FormatsTimeMappings {
		if strings.Contains(result, k) {
			result = strings.ReplaceAll(result, k, v)
		}
	}

	return result
}

func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := float64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	return fmt.Sprintf("%.2f %s", float64(b)/div, units[exp])
}
