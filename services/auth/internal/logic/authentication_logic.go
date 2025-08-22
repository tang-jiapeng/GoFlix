package logic

import (
	"auth/authRpc"
	"context"
	"encoding/json"
	"time"

	"auth/internal/svc"

	"github.com/golang-jwt/jwt/v5"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthenticationLogic) Authentication(in *authRpc.AuthenticationReq) (*authRpc.AuthenticationResp, error) {
	token, err := jwt.ParseWithClaims(in.Token, &svc.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return l.svcCtx.Secret, nil
	})
	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	res, err := l.svcCtx.RDB.Get(timeout, in.SessionId).Result()
	cancel()
	if err != nil {
		return nil, err
	}
	s := &svc.Session{}
	err = json.Unmarshal([]byte(res), s)
	if err != nil {
		return nil, err
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, svc.JwtClaims{
		UserId: s.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	tokenStr, err := token.SignedString([]byte(l.svcCtx.Secret))
	if err != nil {
		return nil, err
	}

	return &authRpc.AuthenticationResp{
		Pass:  true,
		Token: tokenStr,
	}, nil
}
