package logic

import (
	"context"
	"encoding/json"
	"time"

	"auth/authRpc"
	"auth/internal/svc"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateVoucherLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateVoucherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVoucherLogic {
	return &CreateVoucherLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateVoucherLogic) CreateVoucher(in *authRpc.CreateVoucherReq) (*authRpc.CreateVoucherResp, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, svc.JwtClaims{
		UserId: in.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenStr, err := token.SignedString([]byte(l.svcCtx.Secret))
	if err != nil {
		return nil, err
	}

	s := svc.Session{
		UserId: in.UserId,
	}
	sessionId := uuid.New().String()
	js, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	err = l.svcCtx.RDB.Set(timeout, sessionId, string(js), time.Hour*24*7).Err()
	cancel()
	if err != nil {
		return nil, err
	}

	return &authRpc.CreateVoucherResp{
		Ok:        true,
		SessionId: sessionId,
		Token:     tokenStr,
	}, nil
}
