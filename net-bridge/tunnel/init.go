package tunnel

import (
	"fmt"

	"github.com/asppj/cnbs/net-bridge/options"
)

func init() {
	// 校验
	if l := len(NewUIDByPrefix(options.HeartbeatNet, 0)); l != options.PrefixLen {
		panic(fmt.Errorf("前缀长度与规定不匹配:PrefixLen(%d)!=NewUIDByPrefix(%d)", options.PrefixLen, l))
	}
}
