package users

import "time"

// User 用户表结构
type User struct {
	Id          string     `gorm:"primaryKey;column:id;type:char(36)"` // UUID v7
	Mobile      string     `gorm:"column:mobile;size:11;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null"`
	Status      int        `gorm:"column:status;default:1;index"` // 1-启用，2-禁用，3-锁定
	LockedAt    *time.Time `gorm:"column:locked_at"`
	LockReason  string     `gorm:"column:lock_reason;size:255"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

