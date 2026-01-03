package users

import (
	"context"
	"errors"
	"fmt"

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

// Insert 创建用户
func (m *gormDao) Insert(ctx context.Context, data *User) (*User, error) {
	if err := m.db.WithContext(ctx).Create(data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}
	return data, nil
}

// FindOne 根据ID查询用户
func (m *gormDao) FindOne(ctx context.Context, id string) (*User, error) {
	var user User
	if err := m.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// FindByMobile 根据手机号查询用户
func (m *gormDao) FindByMobile(ctx context.Context, mobile string) (*User, error) {
	var user User
	if err := m.db.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// Update 更新用户信息
func (m *gormDao) Update(ctx context.Context, data *User) error {
	if err := m.db.WithContext(ctx).Save(data).Error; err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

// UpdatePassword 更新密码
func (m *gormDao) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	if err := m.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		Update("password_hash", passwordHash).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}
	return nil
}

// UpdateStatus 更新账户状态
func (m *gormDao) UpdateStatus(ctx context.Context, id string, status int) error {
	if err := m.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("更新账户状态失败: %w", err)
	}
	return nil
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

