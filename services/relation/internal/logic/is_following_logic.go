package logic

import (
	"context"

	"GoFlix/services/relation/internal/svc"
	"GoFlix/services/relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFollowingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowingLogic {
	return &IsFollowingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowingLogic) IsFollowing(in *relationRpc.IsFollowingReq) (*relationRpc.IsFollowingResp, error) {
	// todo: add your logic here and delete this line

	return &relationRpc.IsFollowingResp{}, nil
}
