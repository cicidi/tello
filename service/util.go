package service

import "time"

func GetTimestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}
