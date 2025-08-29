package connection

import (
	"sync"

	"github.com/panjf2000/gnet"
)

type Conn struct {
	// 连接id
	id string
	// 上次接收报文时间
	last int64
	conn gnet.Conn
}

var idMutex sync.Mutex
var nextId int64
