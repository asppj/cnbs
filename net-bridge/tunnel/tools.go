package tunnel

import (
	"io"
	"io/ioutil"
	"net"

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

// NewUIDByPrefix NewUIDByPrefix
// p 前缀
// length 预分配buff内存
func NewUIDByPrefix(p options.NetType, length int) []byte {
	buff := make([]byte, 0, 17+length)
	buff = append(buff, byte(p))
	buff = append(buff, NewUID()...)
	return buff
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
func ReadHTTP(proxy io.ReadWriter, bridge gtcp.Conn) (recvByte, sendByte int, err error) {
	buf := make([]byte, 1024)
	for {
		n, err := proxy.Read(buf)
		if err == io.EOF {
			log.Info("proxy连接断开")
			return
		}
		if err != nil {
			continue
		}
		log.Info(string(buf[:n]))
		recvByte += n
		data, err := bridge.SendRecvPkg(buf[:n])
		if err != nil {
			log.Info("bridge连接断开")
			return
		}
		n, err = proxy.Write(data)
		if err != nil {
			log.Info(err)
			return
		}
		sendByte += n
	}

}
