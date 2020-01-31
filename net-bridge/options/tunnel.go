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
	TimeOut = 60 * 60 * time.Second
	// HeartBeat 心跳
	HeartBeat = 15 * time.Second
)

// BuffIterator 读取缓存
type BuffIterator chan BuffType

// BuffType 缓存类型
type BuffType []byte

// NewBuff 新建缓存
func NewBuff() BuffType {
	return make(BuffType, BuffSize)
}

// NewBuffWithPrefix 新建缓存
func NewBuffWithPrefix(netType NetType, uid []byte) BuffType {
	buf := make(BuffType, BuffSize+len(uid)+1)
	buf[0] = byte(netType)
	for i, v := range uid {
		buf[i+1] = v
	}
	return buf
}

// NewBuffGIterator 新建缓存迭代器
func NewBuffGIterator() BuffIterator {
	return make(BuffIterator)
}

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
