package server

import (
	"GoFlix/common/infra/hotkey-go/model"
	"GoFlix/common/infra/hotkey-go/worker/connection"
	"GoFlix/common/infra/hotkey-go/worker/group"
)

var (
	strategies map[string]MessageStrategy
)

func init() {
	strategies = make(map[string]MessageStrategy)
	Register(model.AddKey, &AddStrategy{})
	Register(model.Ping, &PingStrategy{})
	Register(model.Pong, &PongStrategy{})
}

func GetStrategy(msgType string) MessageStrategy {
	return strategies[msgType]
}

func Register(msgType string, strategy MessageStrategy) {
	strategies[msgType] = strategy
}

func (as *AddStrategy) Handle(msg *model.ClientMessage, conn *connection.Conn) {
	keys := make([]string, len(msg.Key))
	times := make([]int64, len(msg.Key))
	i := 0
	for k, v := range msg.Key {
		keys[i] = k
		times[i] = int64(v)
		i++
	}
	err := group.GetGroupMap().AddKey(msg.GroupName, conn, keys, times)
	if err != nil {
		conn.Close()
	}
}

func (ps *PingStrategy) Handle(msg *model.ClientMessage, conn *connection.Conn) {
	conn.Pong()
}

func (ps *PongStrategy) Handle(msg *model.ClientMessage, conn *connection.Conn) {
	conn.ReSetTime()
}
