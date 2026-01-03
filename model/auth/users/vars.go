package users

import "github.com/jinguoxing/idrm-cursor-demo/pkg/errorx"

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errorx.NewWithMsg(errorx.ErrCodeNotFound, "用户不存在")

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errorx.NewWithMsg(errorx.ErrCodeAlreadyExists, "该手机号已注册")

	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = errorx.NewWithMsg(errorx.ErrCodeAuth, "手机号或密码错误")

	// ErrAccountDisabled 账户已禁用
	ErrAccountDisabled = errorx.NewWithMsg(errorx.ErrCodeForbidden, "当前用户存在异常，请联系管理员")

	// ErrAccountLocked 账户已锁定
	ErrAccountLocked = errorx.NewWithMsg(errorx.ErrCodeForbidden, "当前用户存在异常，请联系管理员")
)

