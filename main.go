package main

import (
	"github.com/asppj/cnbs/log"
	"github.com/asppj/cnbs/net-bridge/client"

	"github.com/asppj/cnbs/net-bridge/server"
)

const (
	cp      = 5008 // http 请求
	tp      = 5009 // net bridge
	localIP = "127.0.0.1"
)

func newServer() {
	s := server.NewServer()
	if err := s.Run(); err != nil {
		log.Error(err)
		return
	}
}

func main() {
	c := client.NewClient()
	if err := c.Run(); err != nil {
		log.Error(err)
		return
	}
}
