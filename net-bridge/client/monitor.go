package client

import (
	"context"
	"time"

	"github.com/asppj/cnbs/net-bridge/options"
	"github.com/asppj/cnbs/net-bridge/tunnel"

	"github.com/gogf/gf/net/gtcp"
)

func readHTTP(ctx context.Context, src *gtcp.Conn, ticker *time.Ticker) {
	defer func() {
		_ = src.Close()
	}()
	recvCh := make(chan [][]byte)
	buff := make([]byte, options.BuffSize)
	recvCh = tunnel.ReadConn(ctx, src, buff, recvCh, ticker)
	fn := func() {
		for _, buf := range <-recvCh {
			// TODO 转发
			print(buf)

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
	return nil
}
