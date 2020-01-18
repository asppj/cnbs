package options

// NetType 代理类型
type NetType byte

// PrefixLen 前缀长度
const PrefixLen = 17
const (
	// 占位
	unknownNet NetType = iota
	// HTTPNet http代理
	HTTPNet
	// HTTPSNet https代理
	HTTPSNet

	// TCPNet TCP代理

	TCPNet
	UDPNet
	// kcpNet
	// sshNet
	// vncNet

	// HeartbeatNet 心跳
	HeartbeatNet
)
