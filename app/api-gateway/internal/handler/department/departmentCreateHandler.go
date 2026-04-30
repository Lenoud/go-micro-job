// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package department

import (
	"net/http"

	"api-gateway/internal/logic/department"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 管理员-创建部门
func DepartmentCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateDepartmentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := department.NewDepartmentCreateLogic(r.Context(), svcCtx)
		resp, err := l.DepartmentCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
