package util

import (
	"time"
)

// NowNano returns UnixNano in UTC
func NowNano() int64 {
	return time.Now().UTC().UnixNano()
}
