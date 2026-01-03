// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"time"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/svc"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/types"
	"github.com/jinguoxing/idrm-cursor-demo/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginHistoryLogic {
	return &LoginHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginHistoryLogic) LoginHistory(req *types.LoginHistoryReq) (resp *types.LoginHistoryResp, err error) {
	// TODO: 从Token获取用户ID（需要实现JWT中间件）
	// 暂时从请求中获取，实际应该从JWT Token中解析
	userID := l.ctx.Value("user_id")
	if userID == nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeUnauthorized, "未登录")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil, errorx.NewWithMsg(errorx.ErrCodeUnauthorized, "用户ID格式错误")
	}

	// 解析分页参数
	page := req.Offset
	if page < 1 {
		page = 1
	}
	pageSize := req.Limit
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询登录历史
	histories, total, err := l.svcCtx.LoginHistoryModel.FindByUserID(l.ctx, userIDStr, page, pageSize)
	if err != nil {
		return nil, errorx.NewWithMsg(errorx.ErrCodeSystem, "查询登录历史失败")
	}

	// 格式化返回数据
	items := make([]types.LoginHistoryItem, 0, len(histories))
	for _, h := range histories {
		// 应用筛选条件
		if req.StartTime != "" {
			startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
			if err == nil && h.LoginAt.Before(startTime) {
				continue
			}
		}
		if req.EndTime != "" {
			endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
			if err == nil && h.LoginAt.After(endTime) {
				continue
			}
		}
		if req.IP != "" && h.IP != req.IP {
			continue
		}

		items = append(items, types.LoginHistoryItem{
			Id:         h.Id,
			IP:         h.IP,
			DeviceType: h.DeviceType,
			DeviceID:   h.DeviceID,
			UserAgent:  h.UserAgent,
			LoginAt:    h.LoginAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.LoginHistoryResp{
		Entries:    items,
		TotalCount: total,
	}, nil
}
