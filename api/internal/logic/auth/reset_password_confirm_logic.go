// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/svc"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/types"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/errorx"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/validator"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordConfirmLogic {
	return &ResetPasswordConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordConfirmLogic) ResetPasswordConfirm(req *types.ResetPasswordConfirmReq) (resp *types.ResetPasswordConfirmResp, err error) {
	// 验证验证码
	if err := l.svcCtx.SMSService.VerifyCode(l.ctx, req.Mobile, "reset", req.Code); err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "验证码错误或已过期")
	}

	// 验证密码强度
	if err := validator.ValidatePasswordStrength(req.Password); err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeParamInvalid, err.Error())
	}

	// 查询用户
	user, err := l.svcCtx.UserModel.FindByMobile(l.ctx, req.Mobile)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeNotFound, "用户不存在")
	}

	// 检查新密码是否与旧密码相同
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err == nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "新密码不能与旧密码相同")
	}

	// 加密新密码
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeSystem, "密码加密失败")
	}

	// 更新密码
	if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, user.Id, string(passwordHash)); err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeSystem, "更新密码失败")
	}

	// 使旧Token失效
	if err := l.svcCtx.JWTService.InvalidateUserTokens(l.ctx, user.Id); err != nil {
		l.Errorf("使Token失效失败: %v", err)
		// 不返回错误，密码重置已成功
	}

	return &types.ResetPasswordConfirmResp{
		Message: "密码重置成功，请重新登录",
	}, nil
}
