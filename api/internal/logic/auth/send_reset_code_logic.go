// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/svc"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/types"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendResetCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendResetCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetCodeLogic {
	return &SendResetCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendResetCodeLogic) SendResetCode(req *types.SendResetCodeReq) (resp *types.SendResetCodeResp, err error) {
	// 检查用户是否存在
	_, err = l.svcCtx.UserModel.FindByMobile(l.ctx, req.Mobile)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeNotFound, "该手机号未注册")
	}

	// 发送密码重置验证码
	if err := l.svcCtx.SMSService.SendCode(l.ctx, req.Mobile, "reset"); err != nil {
		return nil, err
	}

	return &types.SendResetCodeResp{
		Message: "验证码已发送",
	}, nil
}
