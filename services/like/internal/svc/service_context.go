package svc

import (
	"GoFlix/common/infra/hotkey-go/hotkey"
	"GoFlix/common/infra/lua"
	"GoFlix/common/util"
	"GoFlix/services/like/internal/config"
	"GoFlix/services/like/internal/script"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/groupcache/singleflight"
	"github.com/redis/go-redis/v9"
	etcd "go.etcd.io/etcd/client/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	Producer sarama.SyncProducer
	Logger   *slog.Logger
	Client   *redis.Client
	Cache    *hotkey.Core
	DB       *gorm.DB
	Group    *singleflight.Group
	Executor *lua.Executor
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		"root", "root", "127.0.0.1", "4000", "goflix",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	rClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6378",
		DB:   1,
	})
	if err := rClient.Ping(context.Background()).Err(); err != nil {
		panic(err.Error())
	}
	eClient, err := etcd.New(etcd.Config{
		Endpoints:   []string{"127.0.0.1:4379"},
		DialTimeout: 3 * time.Second,
	})
	cache, err := hotkey.NewCore("like.rpc", eClient,
		hotkey.WithCacheSize(1024*1024*1024),
		hotkey.WithChannelSize(1024*32),
	)
	if err != nil {
		panic(err.Error())
	}

	executor := lua.NewExecutor(rClient)
	_, err = executor.Load(context.Background(), []*lua.Script{
		script.List,
		script.Set,
		script.BuildList,
	})
	if err != nil {
		panic(err.Error())
	}
	logger, err := util.InitLog("like.rpc", slog.LevelDebug)
	if err != nil {
		panic(err.Error())
	}

	return &ServiceContext{
		Config:   c,
		Producer: nil,
		Logger:   logger,
		Client:   rClient,
		Cache:    cache,
		DB:       db,
		Executor: executor,
	}
}
