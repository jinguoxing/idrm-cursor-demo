// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/svc"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/types"
	"github.com/jinguoxing/idrm-cursor-demo/model/auth/login_history"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/errorx"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/uuid"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 查询用户
	user, err := l.svcCtx.UserModel.FindByMobile(l.ctx, req.Mobile)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeAuth, "手机号或密码错误")
	}

	// 检查账户状态
	if user.Status != 1 {
		return nil, errorx.NewWithMsg(errorx.ErrCodeForbidden, "当前用户存在异常，请联系管理员")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeAuth, "手机号或密码错误")
	}

	// 生成JWT Token
	token, err := l.svcCtx.JWTService.GenerateToken(l.ctx, user.Id)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeSystem, "生成Token失败")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := l.svcCtx.UserModel.Update(l.ctx, user); err != nil {
		l.Errorf("更新最后登录时间失败: %v", err)
	}

	// 解析设备信息
	userAgent := l.ctx.Value("User-Agent")
	if userAgent == nil {
		if r, ok := l.ctx.Value("httpRequest").(*http.Request); ok {
			userAgent = r.Header.Get("User-Agent")
		}
	}
	deviceType, deviceID := l.svcCtx.DeviceParser(getString(userAgent))

	// 获取IP地址
	ip := getClientIP(l.ctx)

	// 记录登录历史
	historyID, err := uuid.GenerateUUID()
	if err != nil {
		l.Errorf("生成登录历史ID失败: %v", err)
	} else {
		history := &login_history.LoginHistory{
			Id:         historyID,
			UserID:     user.Id,
			IP:         ip,
			DeviceType: deviceType,
			DeviceID:   deviceID,
			UserAgent:  getString(userAgent),
			LoginAt:    now,
		}
		if _, err := l.svcCtx.LoginHistoryModel.Insert(l.ctx, history); err != nil {
			l.Errorf("记录登录历史失败: %v", err)
		}
	}

	return &types.LoginResp{
		Token: token,
	}, nil
}

// getString 安全获取字符串
func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// getClientIP 获取客户端IP
func getClientIP(ctx context.Context) string {
	// TODO: 从请求头中获取真实IP
	return "127.0.0.1"
}
