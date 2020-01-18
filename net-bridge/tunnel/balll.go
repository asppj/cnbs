package tunnel

import (
	"net"

	"github.com/gogf/gf/net/gtcp"

	"github.com/asppj/cnbs/net-bridge/auth"
)

// http clients
type HTTPConn struct {
	Identity   *auth.Identity
	BridgeConn gtcp.PoolConn
}

// NewHTTPConn 新建
func NewHTTPConn(ident *auth.Identity, cc gtcp.PoolConn) *HTTPConn {
	return &HTTPConn{
		Identity:   ident,
		BridgeConn: cc,
	}
}

// TCPConn clients
type TCPConn struct {
	Identity   *auth.Identity
	BridgeConn net.Conn
}

// NewTCPConn 新建
func NewTCPConn(ident *auth.Identity, cc net.Conn) *TCPConn {
	return &TCPConn{
		Identity:   ident,
		BridgeConn: cc,
	}
}

// UDPConn clients
type UDPConn struct {
	Identity   *auth.Identity
	BridgeConn net.Conn
}

// NewUDPConn 新建
func NewUDPConn(ident *auth.Identity, cc net.Conn) *UDPConn {
	return &UDPConn{
		Identity:   ident,
		BridgeConn: cc,
	}
}
