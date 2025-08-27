package logic

import (
	"context"

	"GoFlix/services/relation/internal/svc"
	"GoFlix/services/relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllFollowerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAllFollowerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllFollowerLogic {
	return &ListAllFollowerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAllFollowerLogic) ListAllFollower(in *relationRpc.ListAllFollowerReq) (*relationRpc.ListFollowerResp, error) {
	// todo: add your logic here and delete this line

	return &relationRpc.ListFollowerResp{}, nil
}
