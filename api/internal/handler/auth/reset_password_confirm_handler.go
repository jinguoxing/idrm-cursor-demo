// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"net/http"

	"github.com/jinguoxing/idrm-cursor-demo/api/internal/logic/auth"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/svc"
	"github.com/jinguoxing/idrm-cursor-demo/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ResetPasswordConfirmHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetPasswordConfirmReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewResetPasswordConfirmLogic(r.Context(), svcCtx)
		resp, err := l.ResetPasswordConfirm(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
