package tunnel

import (
	"net"
	"time"
)

// TCPChannel TCP 连接
type TCPChannel struct {
	unionKeyAuth string
	// 代理
	proxyConn net.Conn
	// 网桥
	bridgeConn net.Conn
	// 接收字节数
	recvByte int64
	// 发送字节数
	sendByte  int64
	startTime time.Time
}
