// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/svc"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendRegisterCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendRegisterCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRegisterCodeLogic {
	return &SendRegisterCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendRegisterCodeLogic) SendRegisterCode(req *types.SendRegisterCodeReq) (resp *types.SendRegisterCodeResp, err error) {
	// 发送注册验证码
	if err := l.svcCtx.SMSService.SendCode(l.ctx, req.Mobile, "register"); err != nil {
		return nil, err
	}

	return &types.SendRegisterCodeResp{
		Message: "验证码已发送",
	}, nil
}
