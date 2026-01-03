package sms

import "context"

// Service 短信验证码服务接口
type Service interface {
	// SendCode 发送验证码
	// mobile: 手机号
	// codeType: 验证码类型（register/reset）
	SendCode(ctx context.Context, mobile string, codeType string) error

	// VerifyCode 验证验证码
	// mobile: 手机号
	// codeType: 验证码类型（register/reset）
	// code: 验证码
	VerifyCode(ctx context.Context, mobile string, codeType string, code string) error

	// CheckRateLimit 检查发送频率限制
	// mobile: 手机号
	CheckRateLimit(ctx context.Context, mobile string) error
}
