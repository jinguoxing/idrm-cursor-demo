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

type ResetPasswordRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordRequestLogic {
	return &ResetPasswordRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordRequestLogic) ResetPasswordRequest(req *types.ResetPasswordRequestReq) (resp *types.ResetPasswordRequestResp, err error) {
	// 检查用户是否存在
	_, err = l.svcCtx.UserModel.FindByMobile(l.ctx, req.Mobile)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeNotFound, "该手机号未注册")
	}

	// 发送密码重置验证码（复用SendResetCode逻辑）
	if err := l.svcCtx.SMSService.SendCode(l.ctx, req.Mobile, "reset"); err != nil {
		return nil, err
	}

	return &types.ResetPasswordRequestResp{
		Message: "验证码已发送",
	}, nil
}
