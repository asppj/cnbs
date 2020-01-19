package tunnel

import (
	"fmt"

	"github.com/asppj/cnbs/net-bridge/options"
)

func init() {
	// 校验
	if p, _ := NewBuffWithPrefix(options.HeartbeatNet, 0); len(p) != options.PrefixLen {
		panic(fmt.Errorf("前缀长度与规定不匹配:PrefixLen(%d)!=NewBuffWithPrefix(%d)", options.PrefixLen, len(p)))
	}
}
