package leaf_go

import (
	"GoFlix/common/infra/leaf-go/segment"
	"GoFlix/common/infra/leaf-go/snowflake"
	"context"
	"errors"
)

var ErrNoModel = errors.New("no such model")

// NewCore 省略factory的简单工厂模式
func NewCore(config Config) (Core, error) {
	switch config.Model {
	case Segment:
		return segment.NewCreator(&segment.Config{
			Name:     config.SegmentConfig.Name,
			UserName: config.SegmentConfig.UserName,
			Password: config.SegmentConfig.Password,
			Address:  config.SegmentConfig.Address,
		})
	case Snowflake:
		return snowflake.NewCreator(context.Background(), &snowflake.Config{
			CreatorName: config.SnowflakeConfig.CreatorName,
			Addr:        config.SnowflakeConfig.Addr,
			EtcdAddr:    config.SnowflakeConfig.EtcdAddr,
		})
	default:
		return nil, ErrNoModel
	}
}

func init() {

}
