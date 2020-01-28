package auth

import (
	"fmt"
	"time"

	"github.com/asppj/cnbs/net-bridge/options"

	"github.com/asppj/cnbs/log"
)

// Identity 身份识别
type Identity struct {
	NetType   []options.NetType `json:"net_type"`   // 代理类型 http,tcp,udp,kcp,ssh等
	AuthKey   string            `json:"auth_key"`   // 身份标识 全局唯一
	HTTPPort  int               `json:"http_port"`  // 本地程序端口，显示详情用
	HTTPSPort int               `json:"https_port"` // 本地程序端口，显示详情用
	TCPPort   int               `json:"tcp_port"`   // 本地程序端口，显示详情用
	UDPPort   int               `json:"udp_port"`   // 本地程序端口，显示详情用
	SSHPort   int               `json:"ssh_port"`   // 本地程序端口，显示详情用
	VNCPort   int               `json:"vnc_port"`   // 本地程序端口，显示详情用
	ProxyPort int               `json:"proxy_port"` // 服务器代理端口
	UserName  string            `json:"user_name"`  // 帐号
	PassWord  string            `json:"pass_word"`  // 密码
	StartTime time.Time         `json:"start_time"` // 创建时间
}

// Login 登录
func (i *Identity) Login() error {
	_lockAuth.Lock()
	defer _lockAuth.Unlock()
	if _, ok := allClientMap[i.AuthKey]; ok {
		return fmt.Errorf("AuthKey已登录：%s", i.AuthKey)
	}
	i.StartTime = time.Now()
	allClientMap[i.AuthKey] = i
	log.Info("登录新Client:\n\t%s\n", i.string())
	return nil
}

// string 概览
func (i *Identity) string() string {
	return fmt.Sprintf("[ AuthKey：%s ] [ UserName：%s ] [ ProxyPort：%d ] [ HTTPPort：%d ]\n\t[StartTime:%s] [running:%vHour]",
		i.AuthKey, i.UserName, i.ProxyPort, i.HTTPPort, i.StartTime.String(), time.Since(i.StartTime).Hours())
}

// NewIdentity 身份
func NewIdentity() *Identity {
	return &Identity{
		NetType:   []options.NetType{options.HTTPNet},
		AuthKey:   "123456kegel",
		HTTPPort:  80,
		HTTPSPort: 443,
		TCPPort:   8000,
		UDPPort:   8001,
		SSHPort:   22,
		VNCPort:   43,
		ProxyPort: 808,
		UserName:  "kegel",
		PassWord:  "lsp",
		StartTime: time.Time{},
	}
}
