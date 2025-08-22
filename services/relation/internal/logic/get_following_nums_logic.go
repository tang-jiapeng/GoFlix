package logic

import (
	"context"

	"relation/internal/svc"
	"relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowingNumsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowingNumsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowingNumsLogic {
	return &GetFollowingNumsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowingNumsLogic) GetFollowingNums(in *relationRpc.GetFollowingNumsReq) (*relationRpc.GetFollowingNumsResp, error) {
	// todo: add your logic here and delete this line

	return &relationRpc.GetFollowingNumsResp{}, nil
}
