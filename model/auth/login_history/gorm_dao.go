package login_history

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func init() {
	RegisterGormFactory(newGormDao)
}

// gormDao GORM实现
type gormDao struct {
	db *gorm.DB
}

// newGormDao 创建GORM DAO
func newGormDao(db *gorm.DB) Model {
	return &gormDao{db: db}
}

// Insert 插入登录历史
func (m *gormDao) Insert(ctx context.Context, data *LoginHistory) (*LoginHistory, error) {
	if err := m.db.WithContext(ctx).Create(data).Error; err != nil {
		return nil, fmt.Errorf("插入登录历史失败: %w", err)
	}
	return data, nil
}

// FindOne 根据ID查询
func (m *gormDao) FindOne(ctx context.Context, id string) (*LoginHistory, error) {
	var history LoginHistory
	if err := m.db.WithContext(ctx).Where("id = ?", id).First(&history).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLoginHistoryNotFound
		}
		return nil, fmt.Errorf("查询登录历史失败: %w", err)
	}
	return &history, nil
}

// FindByUserID 根据用户ID分页查询
func (m *gormDao) FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]*LoginHistory, int64, error) {
	var histories []*LoginHistory
	var total int64

	offset := (page - 1) * pageSize
	query := m.db.WithContext(ctx).Where("user_id = ?", userID)

	// 统计总数
	if err := query.Model(&LoginHistory{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计登录历史失败: %w", err)
	}

	// 分页查询
	if err := query.Order("login_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		return nil, 0, fmt.Errorf("查询登录历史失败: %w", err)
	}

	return histories, total, nil
}

// DeleteOldRecords 删除过期记录（90天前）
func (m *gormDao) DeleteOldRecords(ctx context.Context, beforeTime time.Time) error {
	if err := m.db.WithContext(ctx).Where("login_at < ?", beforeTime).
		Delete(&LoginHistory{}).Error; err != nil {
		return fmt.Errorf("删除过期记录失败: %w", err)
	}
	return nil
}

// CountByUserID 统计用户登录历史数量
func (m *gormDao) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := m.db.WithContext(ctx).Model(&LoginHistory{}).
		Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计登录历史失败: %w", err)
	}
	return count, nil
}

// WithTx 支持事务
func (m *gormDao) WithTx(tx interface{}) Model {
	if gormTx, ok := tx.(*gorm.DB); ok {
		return &gormDao{db: gormTx}
	}
	return m
}

// Trans 事务执行
func (m *gormDao) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txModel := &gormDao{db: tx}
		return fn(ctx, txModel)
	})
}

