package server

import (
	"io"

	"github.com/asppj/cnbs/log"
	"github.com/asppj/cnbs/net-bridge/tunnel"
	"github.com/gogf/gf/net/gtcp"
)

// 读取代理端口请求
func (s *Server) proxyHTTPHandle(conn *gtcp.Conn) {
	err := tunnel.SetDeadLine(conn)
	if err != nil {
		log.Error("设置超时时间失败")
		return
	}
	defer func() {
		if err != io.EOF {
			if err2 := conn.Close(); err != nil {
				log.Error(err2)
			}
		}

	}()
	bs, err := s.httpBalance.Balance("")
	if err != nil {
		log.Error("没有可用的隧道")
		return
	}
	recv, send, err := tunnel.ReadHTTP(s.ctx, conn, bs.BridgeConn)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(recv, send)
}
