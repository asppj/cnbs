package tunnel

import (
	"sync"

	"github.com/asppj/cnbs/net-bridge/options"

	"github.com/gogf/gf/net/gtcp"
)

var globalChatIDMap = newGlobalBridgeRoom() // 保存会话映射关系

// SetChat 保存临时通信channel
func SetChat(conn *gtcp.Conn, chatID string, ch options.BuffIterator) {
	globalChatIDMap.setChat(conn, chatID, ch)
}

// GetChat 获取通信channel
func GetChat(conn *gtcp.Conn, chatID string) options.BuffIterator {
	return globalChatIDMap.getChat(conn, chatID)
}

// DeleteChat 删除临时通信channel
func DeleteChat(conn *gtcp.Conn, chatID string) {
	globalChatIDMap.deleteCh(conn, chatID)
}

// DeleteConn 客户端断开，删除全部临时通信channel
func DeleteConn(conn *gtcp.Conn) {
	globalChatIDMap.deleteConn(conn)
}

type globalBridgeRoom struct {
	sync.RWMutex
	Chats map[*gtcp.Conn]*bridgeChat
}

func newGlobalBridgeRoom() *globalBridgeRoom {
	chats := make(map[*gtcp.Conn]*bridgeChat)
	return &globalBridgeRoom{
		RWMutex: sync.RWMutex{},
		Chats:   chats,
	}
}

func (g *globalBridgeRoom) setChat(conn *gtcp.Conn, chatID string, ch options.BuffIterator) {
	g.Lock()
	defer g.Unlock()
	if room, ok := g.Chats[conn]; ok {
		room.setChat(chatID, ch)
	} else {
		bc := newBridgeChat()
		bc.setChat(chatID, ch)
		g.Chats[conn] = bc
	}
}
func (g *globalBridgeRoom) getChat(conn *gtcp.Conn, chatID string) options.BuffIterator {
	g.Lock()
	defer g.Unlock()
	if room, ok := g.Chats[conn]; ok {
		return room.getChat(chatID)
	}
	return nil
}

func (g *globalBridgeRoom) deleteCh(conn *gtcp.Conn, chatID string) {
	g.Lock()
	defer g.Unlock()
	if room, ok := g.Chats[conn]; ok {
		room.deleteChat(chatID)
	}
}
func (g *globalBridgeRoom) deleteConn(conn *gtcp.Conn) {
	g.Lock()
	defer g.Unlock()
	delete(g.Chats, conn)
}

type bridgeChat struct {
	sync.RWMutex
	Chs map[string]options.BuffIterator
}

func newBridgeChat() *bridgeChat {
	chs := make(map[string]options.BuffIterator)
	return &bridgeChat{
		RWMutex: sync.RWMutex{},
		Chs:     chs,
	}
}
func (b *bridgeChat) setChat(chatID string, ch options.BuffIterator) {
	b.Lock()
	defer b.Unlock()
	b.Chs[chatID] = ch
}

func (b *bridgeChat) getChat(chatID string) (ch options.BuffIterator) {
	b.Lock()
	defer b.Unlock()
	if ch, ok := b.Chs[chatID]; ok {
		return ch
	}
	return nil
}

func (b *bridgeChat) deleteChat(chatID string) {
	b.Lock()
	defer b.Unlock()
	delete(b.Chs, chatID)
}
