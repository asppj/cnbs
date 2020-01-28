package options

import "time"

const (
	// BuffSize 缓冲区大小
	BuffSize = 256
	// ReadSize 单次读大小
	ReadSize = 128
	// DeadLine io接口超时时间
	DeadLine = 60 * 10 * time.Second
	// TimeOut 整体超时时间
	TimeOut   = 60 * 60 * time.Second
	HeartBeat = 15 * time.Second
)

// NewTickerSystem 24Hour
func NewTickerSystem() *time.Ticker {
	return time.NewTicker(TimeOut)
}

// NewTickerSecond 接口级别
func NewTickerSecond() *time.Ticker {
	return time.NewTicker(DeadLine)
}

// NewTickerHeart 心跳
func NewTickerHeart() *time.Ticker {
	return time.NewTicker(HeartBeat)
}
