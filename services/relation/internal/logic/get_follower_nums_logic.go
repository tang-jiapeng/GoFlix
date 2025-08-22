package logic

import (
	"context"

	"relation/internal/svc"
	"relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerNumsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerNumsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerNumsLogic {
	return &GetFollowerNumsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerNumsLogic) GetFollowerNums(in *relationRpc.GetFollowerNumsReq) (*relationRpc.GetFollowerNumsResp, error) {
	// todo: add your logic here and delete this line

	return &relationRpc.GetFollowerNumsResp{}, nil
}
