package svc

import (
	"GoFlix/common/infra/hotkey-go/hotkey"
	"GoFlix/common/infra/lua"
	"GoFlix/common/util"
	"GoFlix/services/content/public/internal/config"
	"GoFlix/services/content/public/internal/script"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	etcd "go.etcd.io/etcd/client/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	Core     *hotkey.Core
	DB       *gorm.DB
	RClient  *redis.Client
	Executor *lua.Executor
	Logger   *slog.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "root", "127.0.0.1", "3306", "goflix",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err.Error())
	}

	logger, err := util.InitLog("publicContent.rpc", slog.LevelDebug)
	if err != nil {
		panic(err.Error())
	}

	e := lua.NewExecutor(client)
	_, err = e.Load(context.Background(), []*lua.Script{
		script.BuildZSet,
		script.RevByScore,
	})
	if err != nil {
		panic(err.Error())
	}

	eClient, err := etcd.New(etcd.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		panic(err.Error())
	}

	core, err := hotkey.NewCore("publicContent.rpc", eClient,
		hotkey.WithCacheSize(1024*1024*1024),
		hotkey.WithChannelSize(1024*64),
	)
	if err != nil {
		panic(err.Error())
	}
	return &ServiceContext{
		Config:   c,
		Core:     core,
		DB:       db,
		RClient:  client,
		Executor: e,
		Logger:   logger,
	}
}
