package tunnel

import (
	"fmt"
	"sync"

	"github.com/gogf/gf/net/gtcp"

	"github.com/asppj/cnbs/net-bridge/auth"

	"github.com/asppj/cnbs/log"
)

// BalanceHTTP 负载均衡
type BalanceHTTP struct {
	BridgeConns map[string]*HTTPConn
}

var _lockBHTTP sync.Mutex

// Add 新增连接
func (b *BalanceHTTP) Add(ident *auth.Identity, cc *gtcp.Conn) error {
	_lockBHTTP.Lock()
	defer _lockBHTTP.Unlock()
	if _, ok := b.BridgeConns[ident.AuthKey]; ok {
		return fmt.Errorf("已存在：%s", ident.AuthKey)
	}
	b.BridgeConns[ident.AuthKey] = NewHTTPConn(ident, cc)
	return nil
}

// Balance 负载均衡
func (b *BalanceHTTP) Balance(c string) (cc *HTTPConn, err error) {
	// 选择
	if c != "" {
		// 指定
		if cc, ok := b.BridgeConns[c]; ok {
			return cc, err
		}
	}
	// 假随机
	for i, cc := range b.BridgeConns {
		log.InfoF("随机选择%s", i)
		return cc, err

	}
	// TODO 去掉
	if c == "" {
		cc = &HTTPConn{
			Identity:   nil,
			BridgeConn: nil,
		}
		return
	}
	err = fmt.Errorf("未找到可供使用的负载")
	return cc, err
}

// Send 转发请求
func (b *BalanceHTTP) Write(c string, buf []byte) (bs *HTTPConn, n int, err error) {
	bs, err = b.Balance(c)
	if err != nil {
		return
	}
	n, err = bs.BridgeConn.Write(buf)
	return
}

// BalanceTCP 负载均衡
type BalanceTCP struct {
	BridgeConns    map[string]*TCPConn
	BridgeReadChan map[string]chan []byte
}

// Add 添加
func (b *BalanceTCP) Add(ident *auth.Identity, cc *gtcp.Conn) error {
	_lockBHTTP.Lock()
	defer _lockBHTTP.Unlock()
	if _, ok := b.BridgeConns[ident.AuthKey]; ok {
		return fmt.Errorf("已存在：%s", ident.AuthKey)
	}
	b.BridgeConns[ident.AuthKey] = NewTCPConn(ident, cc)
	return nil
}

// BalanceUDP 负载均衡
type BalanceUDP struct {
	BridgeConns    map[string]*UDPConn
	BridgeReadChan map[string]chan []byte
}

// Add 添加
func (b *BalanceUDP) Add(ident *auth.Identity, cc *gtcp.Conn) error {
	_lockBHTTP.Lock()
	defer _lockBHTTP.Unlock()
	if _, ok := b.BridgeConns[ident.AuthKey]; ok {
		return fmt.Errorf("已存在：%s", ident.AuthKey)
	}
	b.BridgeConns[ident.AuthKey] = NewUDPConn(ident, cc)
	return nil
}
