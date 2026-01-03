// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"time"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/svc"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/types"
	"github.com/jinguoxing/idrm-cursor-demo/model/auth/users"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/errorx"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/uuid"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/validator"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 验证验证码
	if err := l.svcCtx.SMSService.VerifyCode(l.ctx, req.Mobile, "register", req.Code); err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "验证码错误或已过期")
	}

	// 验证密码强度
	if err := validator.ValidatePasswordStrength(req.Password); err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeParamInvalid, err.Error())
	}

	// 检查手机号是否已注册
	_, err = l.svcCtx.UserModel.FindByMobile(l.ctx, req.Mobile)
	if err == nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeAlreadyExists, "该手机号已注册")
	}

	// 生成UUID v7作为用户ID
	userID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeSystem, "生成用户ID失败")
	}

	// 加密密码
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeSystem, "密码加密失败")
	}

	// 创建用户
	now := time.Now()
	user := &users.User{
		Id:           userID,
		Mobile:       req.Mobile,
		PasswordHash: string(passwordHash),
		Status:       1, // 启用
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		UserId: userID,
	}, nil
}
