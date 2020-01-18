package auth

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/asppj/cnbs/net-bridge/options"

	"github.com/asppj/cnbs/log"
	"github.com/go-errors/errors"
)

// 保存在线client
var allClientMap = make(map[string]*Identity)
var _lockAuth = sync.Mutex{}

// IdentAuth 身份识别
func IdentAuth(ident []byte) (u *Identity, err error) {
	u = &Identity{}
	if err = json.Unmarshal(ident, u); err != nil {
		log.Info("身份格式错误", errors.New(err))
		return
	}
	// 验证身份逻辑 TODO
	if len(u.NetType) == 0 {
		err = fmt.Errorf("没有指定业务类型")
	}
	for _, t := range u.NetType {
		if t != options.HTTPNet && t != options.TCPNet {
			err = fmt.Errorf("不允许的业务类型：%d", u.NetType)
			return
		}
	}
	if _, ok := allClientMap[u.AuthKey]; ok {
		err = fmt.Errorf("AuthKey已登录：%s", u.AuthKey)
		return
	}
	//  userName
	return
}
