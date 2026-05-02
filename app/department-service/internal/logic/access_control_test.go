package logic

import (
	"context"
	"testing"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"
)

func TestListRejectsMissingAuthContext(t *testing.T) {
	t.Parallel()

	logic := NewListLogic(context.Background(), &svc.ServiceContext{})

	resp, err := logic.List(&department.DepartmentListReq{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeForbidden {
		t.Fatalf("List() code = %v, want %v", resp.GetCode(), common.CodeForbidden)
	}
}
