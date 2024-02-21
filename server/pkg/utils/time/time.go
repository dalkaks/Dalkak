package timeutil

import (
	"time"
)

func GetTimestamp() int64 {
	return time.Now().UTC().Unix()
}
// time.Unix(time, 0).UTC().Format("2006-01-02 15:04:05")
