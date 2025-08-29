package server

import (
	"GoFlix/common/infra/hotkey-go/model"
	"GoFlix/common/infra/hotkey-go/worker/connection"
	"encoding/json"
	"log/slog"

	"github.com/panjf2000/gnet"
)

func (h *Handler) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(connection.NewConn(c))
	return
}

// React tcp报文处理
func (h *Handler) React(packet []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	msg := &model.ClientMessage{}
	err := json.Unmarshal(packet, msg)
	if err != nil {
		slog.Warn("json unmarshal:" + err.Error())
		return nil, gnet.None
	}
	s := GetStrategy(msg.Type)
	if s == nil {
		slog.Warn("get strategy fail unknow type:" + msg.Type)
		return nil, gnet.Close
	}
	ctx := c.Context()
	// 协程池处理
	_ = h.pool.Submit(func() {
		s.Handle(msg, ctx.(*connection.Conn))
	})

	return nil, gnet.None
}

func (h *Handler) OnClosed(_ gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("connect closed")
	}

	return gnet.None
}
