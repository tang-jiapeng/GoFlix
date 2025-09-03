package logic

import (
	"GoFlix/common/model/database"
	"GoFlix/common/util"
	"context"
	"errors"
	"time"

	"GoFlix/services/relation/internal/svc"
	"GoFlix/services/relation/relationRpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	db := l.svcCtx.DB
	creator := l.svcCtx.Creator
	logger := util.SetTrace(l.ctx, l.svcCtx.Logger)
	logger.Info("user following", "userId", in.UserId, "followingId", in.FollowId)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	tx := db.WithContext(timeout).Begin()

	nums := &database.FollowerNums{}
	//  锁关注计数
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(nums, in.UserId).Error
	if err != nil {
		logger.Error("lock table-following_nums:" + err.Error())
		tx.Commit()
		return nil, err
	} else if nums.Nums == 2000 {
		logger.Info("following nums is not enough", "nums", nums.Nums)
		tx.Commit()
		return nil, errors.New("following nums is not enough")
	}
	logger.Debug("lock table-following_nums")

	record := &database.Following{}
	// 关系查询
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("follower_id = ? and type in (0,1) and following_id = ?", in.UserId, in.FollowId).
		Take(record).Error
	// 出错
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("search from table-followings:" + err.Error())
		tx.Commit()
		return nil, err
		// 有记录
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 已经关注
		if record.Type == database.Followed {
			logger.Info("also following")
			return &relationRpc.Empty{}, nil
		}
		// 未关注
		err = tx.Take(record).Update("type", database.Followed).Error
		if err != nil {
			logger.Error("update table-followings:" + err.Error())
			tx.Rollback()
			return nil, err
		}
		// 无记录(关系)
	} else {
		id, ok := creator.GetId()
		if !ok {
			logger.Error("get unique id failed")
			return nil, errors.New("id create failed")
		}
		err := tx.Create(&database.Following{
			Id:          id,
			FollowerId:  in.UserId,
			FollowingId: in.FollowId,
			Type:        database.Followed,
		}).Error
		if err != nil {
			logger.Error("create record from table-followings:" + err.Error())
			tx.Rollback()
			return nil, err
		}
	}
	logger.Debug("update table-followings")
	// 关注计数更新
	err = tx.Take(&database.FollowingNums{}, in.UserId).
		Update("nums", gorm.Expr("nums + 1")).Error
	if err != nil {
		logger.Error("update table-following_nums:" + err.Error())
		tx.Rollback()
		return nil, err
	}
	logger.Debug("update table-following_nums")

	tx.Commit()
	return &relationRpc.Empty{}, nil
}
