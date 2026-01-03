package login_history

import (
	"context"
	"time"
)

// Model 登录历史数据访问接口
type Model interface {
	// Insert 插入登录历史
	Insert(ctx context.Context, data *LoginHistory) (*LoginHistory, error)

	// FindOne 根据ID查询
	FindOne(ctx context.Context, id string) (*LoginHistory, error)

	// FindByUserID 根据用户ID分页查询
	FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]*LoginHistory, int64, error)

	// DeleteOldRecords 删除过期记录（90天前）
	DeleteOldRecords(ctx context.Context, beforeTime time.Time) error

	// CountByUserID 统计用户登录历史数量
	CountByUserID(ctx context.Context, userID string) (int64, error)

	// WithTx 支持事务
	WithTx(tx interface{}) Model

	// Trans 事务执行
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}

