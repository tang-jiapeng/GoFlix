package main

import (
	"GoFlix/common/infra/hotkey-go/worker/server"
	"GoFlix/common/infra/hotkey-go/worker/service"
	"strconv"
	"time"
)

func main() {
	err := service.RegisterService([]string{"127.0.0.1:2379"}, "127.0.0.1:23030", "worker/"+strconv.FormatInt(time.Now().UnixNano(), 10))

	if err != nil {
		panic(err.Error())
	}

	err = server.Serve("tcp://0.0.0.0:23030")
	if err != nil {
		panic(err.Error())
	}
}
