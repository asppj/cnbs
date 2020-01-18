package main

import (
	"github.com/asppj/cnbs/log"

	"github.com/asppj/cnbs/net-bridge/server"
)

const (
	cp      = 5008 // http 请求
	tp      = 5009 // net bridge
	localIP = "127.0.0.1"
)

var s *server.Server

func main() {
	s = server.NewServer(localIP, cp, tp)
	if err := s.Start(); err != nil {
		log.Error(err)
		return
	}
	s.Wait()
}
