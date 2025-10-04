package svc

import (
	"GoFlix/common/infra/hotkey-go/hotkey"
	leaf "GoFlix/common/infra/leaf-go"
	"GoFlix/common/infra/lua"
	"GoFlix/common/util"
	"GoFlix/services/relation/internal/config"
	"GoFlix/services/relation/internal/script"
	"context"
	"fmt"
	"log/slog"

	"github.com/golang/groupcache/singleflight"
	"github.com/redis/go-redis/v9"
	etcd "go.etcd.io/etcd/client/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	RClient  *redis.Client
	Creator  leaf.Core
	Logger   *slog.Logger
	Single   *singleflight.Group
	Executor *lua.Executor
	HotKey   *hotkey.Core
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		"root", "root", "127.0.0.1", "4000", "goflix",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	r := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6378",
		DB:   0,
	})
	creator, err := leaf.NewCore(leaf.Config{
		Model: leaf.Snowflake,
		SnowflakeConfig: &leaf.SnowflakeConfig{
			CreatorName: "relation.rpc",
			Addr:        "127.0.0.1:8081",
			EtcdAddr:    []string{"127.0.0.1:4379"},
		},
	})
	if err != nil {
		panic(err.Error())
	}
	logger, err := util.InitLog("relation.rpc", slog.LevelDebug)
	if err != nil {
		panic(err.Error())
	}
	e := lua.NewExecutor(r)
	_, err = e.Load(context.Background(), []*lua.Script{
		script.BuildZSet,
		script.RevRangeZSet,
		script.GetFiled,
	})
	if err != nil {
		panic(err.Error())
	}
	eClient, err := etcd.New(etcd.Config{
		Endpoints: []string{"127.0.0.1:4379"},
	})
	if err != nil {
		panic(err.Error())
	}
	core, err := hotkey.NewCore("relation.rpc", eClient, hotkey.WithCacheSize(1024*1024*1024))
	if err != nil {
		panic(err.Error())
	}
	svc := &ServiceContext{
		Config:   c,
		DB:       db,
		RClient:  r,
		Creator:  creator,
		Logger:   logger,
		Executor: e,
		Single:   &singleflight.Group{},
		HotKey:   core,
	}
	return svc
}
