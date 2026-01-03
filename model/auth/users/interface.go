package users

import "context"

// Model 用户数据访问接口
type Model interface {
	// Insert 创建用户
	Insert(ctx context.Context, data *User) (*User, error)

	// FindOne 根据ID查询用户
	FindOne(ctx context.Context, id string) (*User, error)

	// FindByMobile 根据手机号查询用户
	FindByMobile(ctx context.Context, mobile string) (*User, error)

	// Update 更新用户信息
	Update(ctx context.Context, data *User) error

	// UpdatePassword 更新密码
	UpdatePassword(ctx context.Context, id string, passwordHash string) error

	// UpdateStatus 更新账户状态
	UpdateStatus(ctx context.Context, id string, status int) error

	// WithTx 支持事务
	WithTx(tx interface{}) Model

	// Trans 事务执行
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}

