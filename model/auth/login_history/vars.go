package login_history

import "github.com/jinguoxing/idrm-cursor-demo/pkg/errorx"

var (
	// ErrLoginHistoryNotFound 登录历史不存在
	ErrLoginHistoryNotFound = errorx.NewWithMsg(errorx.ErrCodeNotFound, "登录历史不存在")
)

