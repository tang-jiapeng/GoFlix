package logic

import (
	"context"

	"GoFlix/services/auth/authRpc"
	"GoFlix/services/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsActiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsActiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsActiveLogic {
	return &IsActiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsActiveLogic) IsActive(in *authRpc.IsActiveReq) (*authRpc.IsActiveResp, error) {

	return &authRpc.IsActiveResp{}, nil
}
