package main

import (
	"flag"

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
func newClient() {
	c := client.NewClient()
	if err := c.Run(); err != nil {
		log.Error(err)
		return
	}
}
func main() {
	model := flag.String("m", "server", "go  run main.go -m client/server")
	flag.Parse()
	flag.PrintDefaults()
	log.Info("运行模式:", model)
	if *model == "client" {
		newClient()
	} else {
		newServer()
	}

}
