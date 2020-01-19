package options

import "time"

const (
	// BuffSize 缓冲区大小
	BuffSize = 256
	// ReadSize 单次读大小
	ReadSize = 128
	// DeadLine io接口超时时间
	DeadLine = 15 * time.Second
	// TimeOut 整体超时时间
	TimeOut = 60 * time.Second
)
