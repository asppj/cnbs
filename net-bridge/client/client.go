package client

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asppj/cnbs/net-bridge/options"

	"github.com/asppj/cnbs/log"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
)

// Client 客户端
type Client struct {
	running        chan os.Signal
	host           string
	remoteHTTPConn *gtcp.Conn
	remoteTCPConn  *gtcp.Conn
	remoteUDPConn  *gtcp.Conn
	localHTTPPort  int
	localTCPPort   int
	localUDPPort   int
	ctx            context.Context
	cancel         context.CancelFunc
	Name           string
}

// NewClient 客户端
// host ip:port
func NewClient() *Client {
	ctx, cancel := context.WithCancel(context.Background())
	// 读取配置
	g.Cfg().SetFileName("config.client.toml")
	localHTTPPort := g.Cfg().GetInt("client.httpPort")
	localTCPConn := g.Cfg().GetInt("client.tcpPort")
	localUDPConn := g.Cfg().GetInt("client.udpPort")
	ip := g.Cfg().GetString("ip")
	port := g.Cfg().GetString("port")
	serviceName := g.Cfg().GetString("serviceName")
	return &Client{
		Name:          serviceName,
		running:       make(chan os.Signal),
		host:          ip + ":" + port,
		localHTTPPort: localHTTPPort,
		localTCPPort:  localTCPConn,
		localUDPPort:  localUDPConn,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Run 开始
func (c *Client) Run() error {
	conn, err := gtcp.NewConn(c.host, options.DeadLine)
	if err != nil {
		return err
	}
	c.remoteHTTPConn = conn
	readHTTP(c.ctx, conn)
	c.Wait()
	return nil
}

// Stop 停止
func (c *Client) Stop() {
	c.cancel()
	if c.remoteHTTPConn != nil {
		if err := c.remoteHTTPConn.Close(); err != nil {
			log.ErrorF("关闭隧道端口:%d失败", c.localHTTPPort)
		} else {
			log.InfoF("关闭隧道端口:%d成功", c.localHTTPPort)
		}
	}
	if c.remoteTCPConn != nil {
		if err := c.remoteTCPConn.Close(); err != nil {
			log.ErrorF("关闭代理端口%d失败", c.localTCPPort)
		} else {
			log.InfoF("关闭代理端口%d成功", c.localTCPPort)
		}
	}
	if c.remoteUDPConn != nil {
		if err := c.remoteUDPConn.Close(); err != nil {
			log.ErrorF("关闭代理端口%d失败", c.localUDPPort)
		} else {
			log.InfoF("关闭代理端口%d成功", c.localUDPPort)
		}
	}
}

// Wait 阻塞等待关闭
func (c *Client) Wait() {
	log.Info("服务运行中...")
	signal.Notify(c.running, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	pTick := time.NewTicker(120 * time.Second)
	defer pTick.Stop()
loop:
	for {
		select {
		case r := <-c.running:
			log.InfoF("正在服务关闭...:%+v。", r)
			break loop
		case <-pTick.C:
			c.PrintInfo()
		}
	}
	c.Stop()
}

// PrintInfo 打印信息
func (c *Client) PrintInfo() {
	log.InfoF("Server Host:%s\n\t httpPort:%d\n\ttcpPort:%d\n\t udpPort:%d",
		c.host, c.localHTTPPort, c.localTCPPort, c.localUDPPort)
}
