package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/os/glog"

	"github.com/gogf/gf/net/gtcp"

	"github.com/asppj/cnbs/log"
	"github.com/asppj/cnbs/net-bridge/tunnel"
)

var _lock sync.Mutex

// Server 网桥服务端
type Server struct {
	Name string
	ip   string
	// 代理套接字
	proxyServer *gtcp.Server
	// 代理端口
	proxyPort int
	// 桥接套接字 与client建立通道
	bridgeServer *gtcp.Server
	// 桥接端口
	bridgePort  int
	httpBalance *tunnel.BalanceHTTP
	tcpBalance  *tunnel.BalanceTCP
	udpBalance  *tunnel.BalanceUDP
	// 所有通道
	channels map[string]*tunnel.TCPChannel
	ctx      context.Context
	cancel   context.CancelFunc
	running  chan os.Signal
}

// PrintInfo 打印信息
func (s *Server) PrintInfo() {
	log.InfoF("代理端口:%d;桥接端口:%d;clients Count:%d", s.proxyPort, s.bridgePort, len(s.channels))
}

// NewServer 创建一个服务器
func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	// 读取配置
	g.Cfg().SetFileName("config.server.toml")
	serviceName := g.Cfg().GetString("serviceName")
	ip := g.Cfg().GetString("serviceIp")
	bPort := g.Cfg().GetInt("servicePort")
	pPort := g.Cfg().GetInt("server.httpPort")
	return &Server{
		ctx:        ctx,
		cancel:     cancel,
		Name:       serviceName,
		running:    make(chan os.Signal),
		ip:         ip,
		proxyPort:  pPort,
		bridgePort: bPort,
		channels:   make(map[string]*tunnel.TCPChannel, 8),
		httpBalance: &tunnel.BalanceHTTP{
			BridgeConns: make(map[string]*tunnel.HTTPConn),
		},
		tcpBalance: &tunnel.BalanceTCP{
			BridgeConns: make(map[string]*tunnel.TCPConn),
		},
		udpBalance: &tunnel.BalanceUDP{
			BridgeConns: make(map[string]*tunnel.UDPConn),
		},
	}
}

// Run 启动服务
// pPort 代理端口
// bPort 桥接端口-与client创建隧道
func (s *Server) Run() (err error) {
	s.proxyServer = gtcp.NewServer(fmt.Sprintf("%s:%d", s.ip, s.proxyPort), s.proxyHTTPHandle)
	go func() {
		if err := s.proxyServer.Run(); err != nil {
			if s.ctx.Err() == nil {
				glog.Error("监听代理端口失败", err)
				s.Stop()
			}
		}
	}()

	s.bridgeServer = gtcp.NewServer(fmt.Sprintf("%s:%d", s.ip, s.bridgePort), s.bridgeHandle)
	go func() {
		if err := s.bridgeServer.Run(); err != nil {
			if s.ctx.Err() == nil {
				glog.Error("监听代理端口失败", err)
				s.Stop()
			}
		}
	}()
	s.PrintInfo()
	s.Wait()
	return
}

// Stop 停止
func (s *Server) Stop() {
	s.cancel()
	if s.bridgeServer != nil {
		if err := s.bridgeServer.Close(); err != nil {
			log.ErrorF("关闭隧道端口:%d失败", s.bridgePort)
		} else {
			log.InfoF("关闭隧道端口:%d成功", s.bridgePort)
		}
	}
	if s.proxyServer != nil {
		if err := s.proxyServer.Close(); err != nil {
			log.ErrorF("关闭代理端口%d失败", s.proxyPort)
		} else {
			log.InfoF("关闭代理端口%d成功", s.proxyPort)
		}
	}
}

// Wait 阻塞等待关闭
func (s *Server) Wait() {
	log.InfoF("服务运行中...:隧道Count:%d", len(s.channels))
	signal.Notify(s.running, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	pTick := time.NewTicker(120 * time.Second)
	defer pTick.Stop()
loop:
	for {
		select {
		case r := <-s.running:
			log.InfoF("正在服务关闭...:%+v。", r)
			break loop
		case <-pTick.C:
			s.PrintInfo()
		}
	}
	s.Stop()
}
