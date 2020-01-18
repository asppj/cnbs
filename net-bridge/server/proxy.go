package server

import (
	"github.com/asppj/cnbs/log"
	"github.com/asppj/cnbs/net-bridge/tunnel"
	"github.com/gogf/gf/net/gtcp"
)

// 读取代理端口请求
func (s *Server) proxyHTTPHandle(conn *gtcp.Conn) {
	bs, err := s.httpBalance.Balance("")
	if err != nil {
		log.Error("没有可用的隧道")
		return
	}
	recv, send, err := tunnel.ReadHTTP(conn, bs.BridgeConn)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(recv, send)
}
