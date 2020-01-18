package tunnel

import (
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/asppj/cnbs/log"

	"github.com/go-errors/errors"
)

// HTTPChannel Http 连接
type HTTPChannel struct {
	unionKeyAuth string
	// 代理
	proxyLis net.TCPListener
	// 网桥
	bridgeConn net.Conn
	// 接收字节数
	recvByte *int64
	// 发送字节数
	sendByte  *int64
	startTime time.Time
}

// NewHTTPChanel 新建HTTP代理服务
func NewHTTPChanel(unionID string, proxy net.TCPListener, bridge net.Conn) *HTTPChannel {
	rb, sb := int64(0), int64(0)
	return &HTTPChannel{
		unionKeyAuth: unionID,
		proxyLis:     proxy,
		bridgeConn:   bridge,
		recvByte:     &rb,
		sendByte:     &sb,
		startTime:    time.Now(),
	}
}

// PrintInfo 打印
func (c *HTTPChannel) PrintInfo() {
	log.Info(strings.Join([]string{
		fmt.Sprintf("创建时间%+v", c.startTime.String()),
		fmt.Sprintf("运行：%+v Hours", time.Now().Sub(c.startTime).Hours()),
		fmt.Sprintf("recv:%d B", c.recvByte),
		fmt.Sprintf("send:%d B", c.sendByte),
		fmt.Sprintf("校验key：%s", c.unionKeyAuth),
	}, "\n"))
}

// Close 关闭通道
func (c *HTTPChannel) Close() {
	if err := c.bridgeConn.Close(); err != nil {
		log.Error("关闭bridgeConn-TCPChannel", errors.New(err))
	}
	if err := c.proxyLis.Close(); err != nil {
		log.Error("关闭proxyLis-TCPChannel", errors.New(err))
	}
	log.Info("关闭TCPChannel")
}

func (c *HTTPChannel) addRecvB(b int64) {
	atomic.AddInt64(c.recvByte, b)
}
func (c *HTTPChannel) addSendB(b int64) {
	atomic.AddInt64(c.sendByte, b)
}

// ListenProxy 监听代理端口
func (c *HTTPChannel) ListenProxy() {
	_, err := c.proxyLis.Accept()
	if err != nil {
		log.Error(errors.New(err))
		return
	}
}
