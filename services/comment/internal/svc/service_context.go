package svc

import (
	"GoFlix/common/infra/hotkey-go/hotkey"
	leaf "GoFlix/common/infra/leaf-go"
	"GoFlix/common/infra/lua"
	syncx "GoFlix/common/infra/sync"
	"GoFlix/common/model/database"
	"GoFlix/common/util"
	"GoFlix/services/comment/commentRpc"
	"GoFlix/services/comment/internal/config"
	"GoFlix/services/comment/internal/script"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
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
	Client   *redis.Client
	Executor *lua.Executor
	DB       *gorm.DB
	Producer sarama.SyncProducer
	Logger   *slog.Logger
	Creator  leaf.Core
	Cache    *hotkey.Core
	Sync     *syncx.Sync
	Group    *singleflight.Group
	ch       chan string
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config: c,
		ch:     make(chan string, 1024*64),
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		"root", "root", "127.0.0.1", "", "goflix",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	svc.DB = db

	rClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6378",
		DB:   0,
	})
	if err := rClient.Ping(context.Background()).Err(); err != nil {
		panic(err.Error())
	}
	svc.Client = rClient

	eClient, err := etcd.New(etcd.Config{
		Endpoints:   []string{"127.0.0.1:4379"},
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		panic(err.Error())
	}

	cache, err := hotkey.NewCore("", eClient,
		hotkey.WithCacheSize(1024*1024*1024),
		hotkey.WithChannelSize(1024*32),
		hotkey.WithObserver(svc),
	)
	if err != nil {
		panic(err.Error())
	}
	svc.Cache = cache

	logger, err := util.InitLog("comment.rpc", slog.LevelDebug)
	if err != nil {
		panic(err.Error())
	}
	svc.Logger = logger

	sync, err := syncx.NewSync(rClient)
	if err != nil {
		panic(err.Error())
	}
	svc.Sync = sync

	executor := lua.NewExecutor(rClient)
	_, err = executor.Load(context.Background(), []*lua.Script{
		script.Build,
		script.GetCountScript,
		script.GetByTime,
		script.GetByHot,
	})
	if err != nil {
		panic(err.Error())
	}
	svc.Executor = executor

	creator, err := leaf.NewCore(leaf.Config{
		Model: leaf.Snowflake,
		SnowflakeConfig: &leaf.SnowflakeConfig{
			CreatorName: "comment.rpc",
			Addr:        "127.0.0.1:8084",
			EtcdAddr:    []string{"127.0.0.1:4379"},
		},
	})
	if err != nil {
		panic(err.Error())
	}
	svc.Creator = creator
	kafkaConfig := sarama.NewConfig()

	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:"}, kafkaConfig)
	if err != nil {
		panic(err.Error())
	}
	svc.Producer = producer

	go svc.addCache()

	return svc
}

func (svc *ServiceContext) addCache() {
	db := svc.DB
	cache := svc.Cache
	for key := range svc.ch {
		if strings.HasPrefix(key, "CommentListByHot:") {
			svc.addHot(key, db, cache)
		} else {
			svc.addTime(key, db, cache)
		}
	}
}

func (svc *ServiceContext) addHot(key string, db *gorm.DB, cache *hotkey.Core) {
	str, _ := strings.CutPrefix("CommentListByHot:", key)
	contentId, _ := strconv.ParseInt(str, 10, 64)
	records := make([]database.Comment, 0)
	err := db.Where("content_id = ? and root_id = ? and status= ?", contentId, 0, database.CommentStatusCommon).
		Limit(1000).Order("hot desc").Find(records).Error
	if err != nil {
		return
	}
	resp := commentRpc.CommentListResp{
		CommentId:   make([]int64, len(records)),
		UserId:      make([]int64, len(records)),
		ContentId:   make([]int64, len(records)),
		RootId:      make([]int64, len(records)),
		ParentId:    make([]int64, len(records)),
		CreatedAt:   make([]int64, len(records)),
		ShortText:   make([]string, len(records)),
		LongTextUri: make([]string, len(records)),
	}
	for i, v := range records {
		resp.CommentId[i] = v.Id
		resp.UserId[i] = v.UserId
		resp.RootId[i] = v.RootId
		resp.ParentId[i] = v.ParentId
		resp.CreatedAt[i] = v.CreatedAt
		resp.ShortText[i] = v.ShortText
		resp.LongTextUri[i] = v.LongTextUri
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		return
	}
	cache.Set(key, b, 15)
	return
}

func (svc *ServiceContext) addTime(key string, db *gorm.DB, cache *hotkey.Core) {
	str, _ := strings.CutPrefix("CommentListByTime:", key)
	contentId, _ := strconv.ParseInt(str, 10, 64)
	records := make([]database.Comment, 0)
	err := db.Where("content_id = ? and root_id = ? and status= ?", contentId, 0, database.CommentStatusCommon).
		Limit(1000).Order("created_at desc").Find(records).Error
	if err != nil {
		return
	}
	resp := commentRpc.CommentListResp{
		CommentId:   make([]int64, len(records)),
		UserId:      make([]int64, len(records)),
		ContentId:   make([]int64, len(records)),
		RootId:      make([]int64, len(records)),
		ParentId:    make([]int64, len(records)),
		CreatedAt:   make([]int64, len(records)),
		ShortText:   make([]string, len(records)),
		LongTextUri: make([]string, len(records)),
	}
	for i, v := range records {
		resp.CommentId[i] = v.Id
		resp.UserId[i] = v.UserId
		resp.RootId[i] = v.RootId
		resp.ParentId[i] = v.ParentId
		resp.CreatedAt[i] = v.CreatedAt
		resp.ShortText[i] = v.ShortText
		resp.LongTextUri[i] = v.LongTextUri
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		return
	}
	cache.Set(key, b, 15)
	return

}

func (svc *ServiceContext) Do(key string) {
	svc.ch <- key
}
