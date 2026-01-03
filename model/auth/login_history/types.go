package login_history

import "time"

// LoginHistory 登录历史表结构
type LoginHistory struct {
	Id        string    `gorm:"primaryKey;column:id;type:char(36)"` // UUID v7
	UserID    string    `gorm:"column:user_id;type:char(36);index;not null"` // UUID v7
	IP        string    `gorm:"column:ip;size:45;not null"`
	DeviceType string   `gorm:"column:device_type;size:20"`
	DeviceID   string   `gorm:"column:device_id;size:255"`
	UserAgent  string   `gorm:"column:user_agent;size:500"`
	LoginAt    time.Time `gorm:"column:login_at;index"`
}

// TableName 指定表名
func (LoginHistory) TableName() string {
	return "login_history"
}

