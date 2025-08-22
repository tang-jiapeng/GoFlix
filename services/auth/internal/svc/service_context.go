package svc

import (
	"auth/internal/config"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	RDB    *redis.Client
	Secret string
}

type JwtClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type Session struct {
	UserId int64 `json:"user_id"`
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6378",
		DB:   1,
	})
	timeout, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	if err := rdb.Ping(timeout).Err(); err != nil {
		panic(err.Error())
	}
	cancel()

	return &ServiceContext{
		Config: c,
		RDB:    rdb,
	}
}
