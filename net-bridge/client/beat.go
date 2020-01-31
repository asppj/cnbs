package client

import (
	"github.com/asppj/cnbs/net-bridge/bridge"
	"github.com/asppj/cnbs/net-bridge/options"
	"github.com/gogf/gf/net/gtcp"
)

// 发送心跳
func sendHeart(conn *gtcp.Conn) error {
	heartBeat, _ := bridge.NewBuffWithPrefix(options.HeartbeatNet, 0)
	return conn.SendPkg(heartBeat)
}
