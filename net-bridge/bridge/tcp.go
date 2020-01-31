package bridge

//
// import (
// 	"fmt"
//
// 	"github.com/asppj/cnbs/net-bridge/options"
// 	"github.com/gogf/gf/net/gtcp"
// )
//
// // 封装
// // 转发流程
// /*
//  添加前缀->转发->解析->回复
// */
//
// // TCPConn TCP网桥
// type TCPConn struct {
// 	Conn *gtcp.Conn
// 	Type options.NetType
// 	UID  []byte // 每个连接唯一
// }
//
// // NewTCPConn 新建
// func NewTCPConn(conn *gtcp.Conn, netType options.NetType) *TCPConn {
// 	return &TCPConn{Conn: conn, Type: netType, UID: NewUID()}
// }
//
// // ReadIterator 缓存迭代器
// func (c *TCPConn) ReadIterator() (ch options.BuffIterator, err error) {
// 	if c.Conn == nil {
// 		err = fmt.Errorf("conn  is nil")
// 		return
// 	}
//
// 	ch = options.NewBuffGIterator()
// 	buf := options.NewBuffWithPrefix(c.Type, c.UID)
// 	prefixLength := len(c.UID) + 1
// 	var n int
// 	loop := func() {
// 		for {
// 			n, err = c.Conn.Read(buf[prefixLength:])
// 			if err != nil {
// 				close(ch)
// 				return
// 			}
// 			ch <- buf[:prefixLength+n]
// 		}
// 	}
// 	go loop()
// 	return
// }
//
// // SendFromIterator 发送
// func (c *TCPConn) SendFromIterator(conn *TCPConn) (n int, err error) {
// 	ch, err := c.ReadIterator()
// 	if err != nil {
// 		return
// 	}
// 	for chEle := range ch {
// 		nt, err := conn.Send(chEle)
// 		if err != nil {
// 			return
// 		}
// 		n += nt
// 	}
// 	return
// }
//
// // Send 发送
// func (c *TCPConn) Send(buf options.BuffType) (int, error) {
// 	return c.Conn.Write(buf)
// }
