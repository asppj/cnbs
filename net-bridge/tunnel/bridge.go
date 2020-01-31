package tunnel

import (
	"context"
	"net"

	"github.com/asppj/cnbs/net-bridge/bridge"

	"github.com/asppj/cnbs/net-bridge/options"
)

type tcpBridge struct {
	proxy          net.Conn        // 代理链接
	bridgeReadChan chan []byte     // 读取response通道，ctx超时结束
	bridge         net.Conn        // 发送请求
	proxyFn        func(int64)     // 统计数据
	bridgeFn       func(int64)     // 统计数据
	ctx            context.Context // 超时context
}

func newTCPBridge(ctx context.Context, proxy, bridge net.Conn, proxyFn, bridgeFn func(int64)) *tcpBridge {
	context.Background().Done()
	return &tcpBridge{
		proxy:          proxy,
		bridge:         bridge,
		bridgeReadChan: make(chan []byte),
		proxyFn:        proxyFn,
		bridgeFn:       bridgeFn,
		ctx:            ctx,
	}
}

func (b *tcpBridge) pack(buff []byte) []byte {
	bl := len(buff)
	prefix, _ := bridge.NewBuffWithPrefix(options.HTTPNet, bl)
	return append(prefix, buff...)
}
func (b *tcpBridge) unpack(buff []byte) []byte {
	return nil
}

// OnceRequest 代理一次请求
func (b *tcpBridge) OnceRequest() (err error) {
	// buf, err := readAll(b.proxy)
	// if err != nil {
	// 	log.Error("接收http代理request失败")
	// 	return
	// }
	//
	// n, err := b.bridge.Send(b.pack(buf))
	// if err != nil {
	// 	log.Error("发送http代理request失败")
	// 	return
	// }
	// b.proxyFn(int64(n))
	// buf, err = readAll(b.bridge)
	// if err != nil {
	// 	log.Error("接受http代理Response失败")
	// 	return
	// }
	// n, err = b.proxy.Send(b.unpack(buf))
	// if err != nil {
	// 	log.Error("发送http代理Response失败")
	// 	return
	// }
	// b.bridgeFn(int64(n))
	return
}

func (b *tcpBridge) OnceRequestWithTimeOut() (err error) {
	return
}

// LongExchange 长连接交换数据
func (b *tcpBridge) LongExchange() (err error) {
	return
}
