package util

import (
	"fmt"
	"time"
)

// GetDate 用于返回 yyyy-mm-dd HH:MM:SS 格式的字符串
func GetDate() string {
	t := time.Now()
	return fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
}
