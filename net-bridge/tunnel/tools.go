package tunnel

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/gogf/gf/net/gtcp"

	"github.com/asppj/cnbs/log"

	"github.com/asppj/cnbs/net-bridge/options"

	uuid "github.com/satori/go.uuid"
)

/*
	http请求前缀：一B类型前缀 16B UUID长度 共17位
	其他类型类似
*/

// NewUID uuid
func NewUID() []byte {
	return uuid.NewV1().Bytes()
}

// NewBuffWithPrefix NewBuffWithPrefix
// p 前缀
// length 预分配buff内存
func NewBuffWithPrefix(p options.NetType, length int) (buff []byte, uuid string) {
	buff = make([]byte, 0, 17+length)
	buff = append(buff, byte(p))
	uid := NewUID()
	buff = append(buff, NewUID()...)
	return buff, string(uid)
}

// SetDeadLine 设置超时时间点
func SetDeadLine(conn *gtcp.Conn) error {
	return conn.SetDeadline(time.Now().Add(options.DeadLine))

}

// readAll 读取所有
func readAll(conn net.Conn) (buff []byte, err error) {
	return ioutil.ReadAll(conn)
}

func exchangePair(proxy net.Conn, bridge net.Conn) (err error) {
	return
}

func authHTTPBridge(conn io.Reader) (key string, err error) {
	// buff := make([]byte, 24)
	// n, err := conn.Read(buff)
	// if err != nil {
	// 	return
	// }
	return
}

// ReadHTTP 读取http
func ReadHTTP(ctx context.Context, proxy *gtcp.Conn, bridge *gtcp.Conn) (recvByte, sendByte int, err error) {
	ticker := time.NewTicker(options.TimeOut)
	// 分片转发-request
	chatID, n, err := exchangeRequest(ctx, proxy, bridge, ticker)
	if err != nil {
		return
	}
	recvByte += n
	// 分片转发-response
	n, err = exchangeResponse(ctx, chatID, bridge, proxy, ticker)
	if err != nil {
		return
	}
	sendByte += n
	return
}

func exchangeRequest(ctx context.Context, src *gtcp.Conn, dst *gtcp.Conn, ticker *time.Ticker) (chatID string, n int, err error) {
	buf, chatID := NewBuffWithPrefix(options.HTTPNet, options.BuffSize)
	ch := make(chan [][]byte)
	ch = ReadConn(ctx, src, buf, ch, ticker)
	// 分片转发-request
	for _, res := range <-ch {
		buf = append(buf, res...)
		err = dst.Send(buf)
		if err != nil {
			return
		}
		n += len(buf)
		buf = buf[:0]
	}
	return
}

// ReadConn 读取
func ReadConn(ctx context.Context, src *gtcp.Conn, buf []byte, ch chan [][]byte, ticker *time.Ticker) chan [][]byte {
	rFn := func() error {
		recv, err := src.Read(buf)
		if err != nil {
			return err
		}
		// 重新定义超时时间
		err = SetDeadLine(src)
		if err != nil {
			return err
		}
		ch <- [][]byte{buf[:recv]}
		return nil
	}
	go func() {
		// 读取完毕关闭通道
		defer close(ch)
		for {
			select {
			case <-ticker.C: // 整体超时
				return
			case <-ctx.Done(): // 主动取消
				return
			default:
				err := rFn()
				if err != nil {
					log.Error(err) // 单次超时或其他错误 TODO io.EOF不应该打印
					return
				}
			}
		}
	}()
	return ch
}

func exchangeResponse(ctx context.Context, chatID string, src *gtcp.Conn, dst *gtcp.Conn, ticker *time.Ticker) (n int, err error) {
	buf := make([]byte, options.ReadSize)
	ch := make(chan [][]byte)

	rFn := func() error {
		recv, err := src.Read(buf)
		if err != nil {
			return err
		}
		// 重新定义超时时间
		err = SetDeadLine(src)
		if err != nil {
			return err
		}
		ch <- [][]byte{buf[:recv]}
		return nil
	}
	go func() {
		// 读取完毕关闭通道
		defer close(ch)
		for {
			select {
			case <-ticker.C: // 整体超时
				return
			case <-ctx.Done(): // 主动取消
				return
			default:
				err := rFn()
				if err != nil {
					log.Error(err) // 单次超时或其他错误 TODO io.EOF不应该打印
					return
				}
			}
		}
	}()
	// 分片转发-request
	for _, res := range <-ch {
		buf = append(buf, res...)
		err = dst.Send(buf)
		if err != nil {
			return
		}
		n += len(buf)
		buf = buf[:0]
	}
	return
}

func readResponse(ctx context.Context, bridge *gtcp.Conn, ticker time.Ticker) error {
	return nil
}
