package tunnel

import (
	"github.com/gogf/gf/net/gtcp"

	"github.com/asppj/cnbs/net-bridge/auth"
)

// HTTPConn clients
type HTTPConn struct {
	Identity   *auth.Identity
	BridgeConn *gtcp.Conn
}

// NewHTTPConn 新建
func NewHTTPConn(ident *auth.Identity, cc *gtcp.Conn) *HTTPConn {
	return &HTTPConn{
		Identity:   ident,
		BridgeConn: cc,
	}
}

// TCPConn clients
type TCPConn struct {
	Identity   *auth.Identity
	BridgeConn *gtcp.Conn
}

// NewTCPConn 新建
func NewTCPConn(ident *auth.Identity, cc *gtcp.Conn) *TCPConn {
	return &TCPConn{
		Identity:   ident,
		BridgeConn: cc,
	}
}

// UDPConn clients
type UDPConn struct {
	Identity   *auth.Identity
	BridgeConn *gtcp.Conn
}

// NewUDPConn 新建
func NewUDPConn(ident *auth.Identity, cc *gtcp.Conn) *UDPConn {
	return &UDPConn{
		Identity:   ident,
		BridgeConn: cc,
	}
}
