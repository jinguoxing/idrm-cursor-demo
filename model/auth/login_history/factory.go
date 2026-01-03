package login_history

import (
	"database/sql"

	"gorm.io/gorm"
)

var (
	gormFactory func(*gorm.DB) Model
)

// RegisterGormFactory 注册GORM工厂函数
func RegisterGormFactory(fn func(*gorm.DB) Model) {
	gormFactory = fn
}

// NewModel 创建Model实例（优先GORM，降级SQLx）
func NewModel(sqlConn *sql.DB, gormDB *gorm.DB) Model {
	if gormDB != nil && gormFactory != nil {
		return gormFactory(gormDB)
	}
	// TODO: 实现SQLx版本
	panic("GORM连接不可用，且SQLx未实现")
}

