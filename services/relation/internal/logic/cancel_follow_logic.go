package logic

import (
	"GoFlix/common/model/database"
	"GoFlix/common/util"
	"context"
	"time"

	"GoFlix/services/relation/internal/svc"
	"GoFlix/services/relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	db := l.svcCtx.DB
	logger := util.SetTrace(l.ctx, l.svcCtx.Logger)
	logger.Info("user cancelFollowed", "userId", in.UserId, "followedId", in.FollowId)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx := db.WithContext(timeout).Begin()
	// 直接更新
	res := tx.Model(&database.Following{}).
		Where("follower_id = ? and type = ? and following_id = ?", in.UserId, database.Followed, in.FollowId).
		Update("type", database.UnFollowed)
	if res.Error != nil {
		logger.Error("update table-following:" + res.Error.Error())
		tx.Rollback()
		return nil, res.Error
	}
	// 没有记录(关系)
	if res.RowsAffected == 0 {
		logger.Info("also cancel following relation")
		tx.Commit()
		return &relationRpc.Empty{}, nil
	}

	logger.Debug("update table-following")
	// 关注数量更新
	err := tx.Take(&database.FollowingNums{}, in.UserId).
		Update("nums", gorm.Expr("nums - 1")).Error
	if err != nil {
		logger.Error("update table-following_nums:" + err.Error())
		tx.Rollback()
		return nil, err
	}
	logger.Debug("update table-following_nums")
	tx.Commit()
	return &relationRpc.Empty{}, nil
}
