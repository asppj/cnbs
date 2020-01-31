package bridge

import (
	"context"
	"time"

	"github.com/asppj/cnbs/net-bridge/tunnel"

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

// UnpackBuff 解隧道包
func UnpackBuff(buf []byte) (chatID string, data []byte, ok bool) {
	if len(buf) < options.PrefixLen {
		return
	}
	chatID = string(buf[:options.PrefixLen])
	data = buf[options.PrefixLen:]
	return
}

// PackBuff 打包隧道数据
func PackBuff(chatID string, buf []byte) (buff []byte) {
	buff = make([]byte, options.PrefixLen+len(buf))
	buff = append(buff, []byte(chatID)...)
	buff = append(buff, buf...)
	return
}

// SetDeadLine 设置超时时间点
func SetDeadLine(conn *gtcp.Conn) error {
	return conn.SetDeadline(time.Now().Add(options.DeadLine))

}

// ProxyHTTP 读取http
func ProxyHTTP(ctx context.Context, proxy *gtcp.Conn, bridge *gtcp.Conn) (recvByte, sendByte int, err error) {
	ticker := time.NewTicker(options.TimeOut)
	// 分片转发-request-不阻塞
	chatID, err := exchangeRequest(ctx, proxy, bridge, ticker, recvByte)
	if err != nil {
		return
	}
	// 阻塞
	sendByte = exchangeResponse(proxy, bridge, chatID)
	return
}

func exchangeResponse(proxy *gtcp.Conn, bridge *gtcp.Conn, chatID string) (sendB int) {
	// 保存chatID 设置回显通道
	chatCh := make(chan [][]byte)
	tunnel.SetChat(bridge, chatID, chatCh)
	for _, buf := range <-chatCh {
		n, err := proxy.Write(buf)
		if err != nil {
			return
		}
		sendB += n
	}
	tunnel.DeleteChat(bridge, chatID)
	return
}

// 不阻塞读取
func exchangeRequest(ctx context.Context, src *gtcp.Conn, dst *gtcp.Conn, ticker *time.Ticker, recvB int) (chatID string, err error) {
	buf, chatID := NewBuffWithPrefix(options.HTTPNet, options.BuffSize)
	ch := ReadConn(ctx, src, buf, ticker)
	// 分片转发-request
	go func() {
		if err := WriteConn(dst, ch, recvB); err != nil {
			log.Error(err)
		}
	}()
	return
}

// ReadConn 读取
func ReadConn(ctx context.Context, src *gtcp.Conn, buf []byte, ticker *time.Ticker) chan [][]byte {
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
	rLoop := func() {
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
	}
	go rLoop()
	return ch
}

// WriteConn 不带超时，read关闭则关闭
func WriteConn(conn *gtcp.Conn, ch chan [][]byte, sendB int) (err error) {
	wLoop := func(buf []byte) error {
		n, wErr := conn.Write(buf)
		if err != nil {
			return wErr
		}
		sendB += n
		return nil
	}

	for _, buf := range <-ch {
		if err = wLoop(buf); err != nil {
			return
		}
	}
	return
}
