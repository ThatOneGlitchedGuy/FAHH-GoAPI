package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func PaginationParams(pageStr, sizeStr string) (int, int) {
	page := 1
	size := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
		size = s
	}

	return page, size
}

func CalculateOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}

func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func StringInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if !seen[s] {
			result = append(result, s)
			seen[s] = true
		}
	}
	return result
}

func FormatPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)
}

func RoundPrice(price float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(price*multiplier) / multiplier
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	unit := "ms"
	val := float64(d.Milliseconds())

	if val >= 1000 {
		val = val / 1000
		unit = "s"
	}

	return fmt.Sprintf("%.2f%s", val, unit)
}

func IsValidUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}
	parts := strings.Split(uuid, "-")
	if len(parts) != 5 {
		return false
	}
	return len(parts[0]) == 8 && len(parts[1]) == 4 && len(parts[2]) == 4 && len(parts[3]) == 4 && len(parts[4]) == 12
}

func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func IsTimeAfter(t1, t2 time.Time) bool {
	return t1.After(t2)
}

func IsTimeEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2)
}

func PointerOf(v interface{}) interface{} {
	return &v
}

func ValueOf(p interface{}) interface{} {
	switch ptr := p.(type) {
	case *string:
		if ptr == nil {
			return ""
		}
		return *ptr
	case *int:
		if ptr == nil {
			return 0
		}
		return *ptr
	case *bool:
		if ptr == nil {
			return false
		}
		return *ptr
	default:
		return nil
	}
}

func CalculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (part / total) * 100
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
