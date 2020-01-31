package client

import (
	"context"

	"github.com/asppj/cnbs/net-bridge/tunnel"

	"github.com/asppj/cnbs/log"

	"github.com/asppj/cnbs/net-bridge/options"
	"github.com/gogf/gf/net/gtcp"
)

func monitorHTTPTunnel(ctx context.Context, src *gtcp.Conn) {
	defer func() {
		_ = src.Close()
	}()
	ticker := options.NewTickerSecond()
	defer ticker.Stop()
	uid := tunnel.NewUID()
	buff := options.NewBuffWithPrefix(options.HTTPNet, uid)
	recvCh := tunnel.ReadConn(ctx, src, buff, ticker)

	fn := func() {
		for buf := range recvCh {
			// TODO 转发
			print(string(buf))
			err := src.SendPkg(buf)
			if err != nil {
				log.Info("转发：")
			}
		}
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				return
			default:
				fn()
			}
		}
	}()
	return
}
