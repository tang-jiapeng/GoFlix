package logic

import (
	"context"

	"relation/internal/svc"
	"relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelFollowLogic {
	return &CancelFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelFollowLogic) CancelFollow(in *relationRpc.CancelFollowReq) (*relationRpc.Empty, error) {
	// todo: add your logic here and delete this line

	return &relationRpc.Empty{}, nil
}
