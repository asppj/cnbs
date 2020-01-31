package server

import (
	"context"
	"time"

	"github.com/asppj/cnbs/net-bridge/tunnel"
	"github.com/gogf/gf/os/glog"

	"github.com/asppj/cnbs/log"
	"github.com/asppj/cnbs/net-bridge/auth"
	"github.com/asppj/cnbs/net-bridge/options"
	"github.com/gogf/gf/net/gtcp"
)

// 建立隧道
func (s *Server) bridgeHandle(conn *gtcp.Conn) {
	buff, err := conn.RecvPkg()
	if err != nil {
		return
	}
	identity, err := auth.IdentAuth(buff)
	if err != nil {
		log.Info("身份验证失败", err)
		return
	}
	err = identity.Login()
	if err != nil {
		// identity.PassWord = ""
		log.InfoF("client登录失败:%+v", identity)
		return
	}
	loop := func() {
		// 登录
		_lock.Lock()
		defer _lock.Unlock()
		for _, netType := range identity.NetType {
			switch netType {
			case options.HTTPNet:
				{
					// 允许多个桥接，相当于负载均衡
					if err := s.httpBalance.Add(identity, conn); err != nil {
						log.Error("http添加失败", err)
						return
					}
					httpBridgeHandle(s.ctx, conn)
				}
			case options.TCPNet:
				{
					if err := s.tcpBalance.Add(identity, conn); err != nil {
						log.Error("tcp添加失败", err)
						return
					}
				}
			case options.UDPNet:
				{
					if err := s.udpBalance.Add(identity, conn); err != nil {
						log.Error("udp添加失败", err)
						return
					}
				}
			default:
				// 一定不会发生，身份验证的时候就会去除
				log.ErrorF("不允许的代理类型:%d", identity.NetType)
				return
			}
		}
	}
	// 循环处理
	for {
		loop()
		s.PrintInfo()
	}
}

// http 隧道监听
func httpBridgeHandle(ctx context.Context, conn *gtcp.Conn) {
	log.Info("启动隧道监听。。。")
	ticker := time.NewTicker(time.Hour * 24)
	rFn := func() {
		buf, err := conn.RecvPkg()
		if err != nil {
			return
		}
		chatID, data, ok := tunnel.UnpackBuff(buf)
		if ok {
			ch := tunnel.GetChat(conn, chatID)
			ch <- data
		} else {
			glog.Warning("无效数据", string(buf))
		}
	}
	loop := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				return
			default:
				rFn()
			}
		}
	}
	go loop()
	log.Info("http隧道监听成功")
}
