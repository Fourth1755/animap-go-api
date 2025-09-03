package utils

import (
	"fmt"
	"time"
)

func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	years := int(diff.Hours() / 24 / 365)
	if years > 0 {
		if years == 1 {
			return "1 ปีที่แล้ว"
		}
		return fmt.Sprintf("%d ปีที่แล้ว", years)
	}

	months := int(diff.Hours() / 24 / 30)
	if months > 0 {
		if months == 1 {
			return "1 เดือนที่แล้ว"
		}
		return fmt.Sprintf("%d เดือนที่แล้ว", months)
	}

	days := int(diff.Hours() / 24)
	if days > 0 {
		if days == 1 {
			return "1 วันที่แล้ว"
		}
		return fmt.Sprintf("%d วันที่แล้ว", days)
	}

	hours := int(diff.Hours())
	if hours > 0 {
		if hours == 1 {
			return "1 ชั่วโมงที่แล้ว"
		}
		return fmt.Sprintf("%d ชั่วโมงที่แล้ว", hours)
	}

	minutes := int(diff.Minutes())
	if minutes > 0 {
		if minutes == 1 {
			return "1 นาทีที่แล้ว"
		}
		return fmt.Sprintf("%d นาทีที่แล้ว", minutes)
	}

	seconds := int(diff.Seconds())
	if seconds < 5 {
		return "เมื่อสักครู่"
	}
	return fmt.Sprintf("%d วินาทีที่แล้ว", seconds)
}
