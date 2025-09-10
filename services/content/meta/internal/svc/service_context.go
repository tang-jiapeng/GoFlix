package svc

import (
	leaf "GoFlix/common/infra/leaf-go"
	"GoFlix/common/util"
	"GoFlix/services/content/meta/internal/config"
	"fmt"
	"log/slog"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Logger  *slog.Logger
	Creator leaf.Core
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "root", "127.0.0.1", "3306", "goflix",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	logger, err := util.InitLog("metaContent.rpc", slog.LevelDebug)
	if err != nil {
		panic(err.Error())
	}

	creator, err := leaf.NewCore(leaf.Config{
		Model: leaf.Snowflake,
		SnowflakeConfig: &leaf.SnowflakeConfig{
			CreatorName: "metaContent.rpc",
			Addr:        "127.0.0.1:8082",
			EtcdAddr:    []string{"127.0.0.1:2379"},
		},
	})
	if err != nil {
		panic(err.Error())
	}
	svc := &ServiceContext{
		Config:  c,
		DB:      db,
		Logger:  logger,
		Creator: creator,
	}
	logx.DisableStat()
	return svc
}
