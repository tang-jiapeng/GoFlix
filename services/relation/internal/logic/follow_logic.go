package logic

import (
	"context"

	"GoFlix/services/relation/internal/svc"
	"GoFlix/services/relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowLogic) Follow(in *relationRpc.FollowReq) (*relationRpc.Empty, error) {
	// todo: add your logic here and delete this line

	return &relationRpc.Empty{}, nil
}
