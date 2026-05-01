// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package oplog

import (
	"net/http"

	"api-gateway/internal/logic/oplog"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 操作日志列表
func OpLogListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpLogListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := oplog.NewOpLogListLogic(r.Context(), svcCtx)
		resp, err := l.OpLogList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
