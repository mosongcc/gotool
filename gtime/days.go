package gtime

import "time"

// Days 计算时间差（天数）
func Days(t1, t2 time.Time) (day int) {
	day = int(t1.Sub(t2).Hours() / 24)
	return
}
