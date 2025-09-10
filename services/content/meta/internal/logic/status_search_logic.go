package logic

import (
	"GoFlix/common/model/database"
	"GoFlix/common/util"
	"context"
	"errors"
	"time"

	"GoFlix/services/content/meta/internal/svc"
	"GoFlix/services/content/meta/metaContentRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatusSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatusSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusSearchLogic {
	return &StatusSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatusSearchLogic) StatusSearch(in *metaContentRpc.StatusSearchReq) (*metaContentRpc.StatusSearchResp, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db := l.svcCtx.DB.WithContext(timeout)
	logger := util.SetTrace(l.ctx, l.svcCtx.Logger)
	logger.Info("user search content status", "userId", in.UserId, "contentId", in.ContentId)

	record := &database.InvisibleContentInfo{}
	err := db.Take(record, in.ContentId).Error
	if err != nil {
		logger.Error("search status err:" + err.Error())
		return nil, err
	}
	if record.UserId != in.UserId {
		logger.Error("user update content info:is not the publish user")
		return nil, errors.New("you can not do this it is not your content")
	}
	return &metaContentRpc.StatusSearchResp{
		Status: int32(record.Status),
		Desc:   record.Desc,
	}, nil
}
