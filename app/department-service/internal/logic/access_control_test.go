package logic

import (
	"context"
	"database/sql"
	"testing"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/model"
	"department-service/internal/svc"
)

type fakeDepartmentModel struct {
	listCalled bool
}

func (m *fakeDepartmentModel) Insert(ctx context.Context, data *model.Department) (sql.Result, error) {
	return nil, nil
}

func (m *fakeDepartmentModel) FindOne(ctx context.Context, id string) (*model.Department, error) {
	return nil, nil
}

func (m *fakeDepartmentModel) FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*model.Department, int64, error) {
	m.listCalled = true
	return []*model.Department{}, 0, nil
}

func (m *fakeDepartmentModel) Update(ctx context.Context, data *model.Department) error {
	return nil
}

func (m *fakeDepartmentModel) Delete(ctx context.Context, ids string) error {
	return nil
}

func TestListRejectsMissingAuthContext(t *testing.T) {
	t.Parallel()

	departmentModel := &fakeDepartmentModel{}
	logic := NewListLogic(context.Background(), &svc.ServiceContext{DepartmentModel: departmentModel})

	resp, err := logic.List(&department.DepartmentListReq{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeForbidden {
		t.Fatalf("List() code = %v, want %v", resp.GetCode(), common.CodeForbidden)
	}
	if departmentModel.listCalled {
		t.Fatal("List() queried departments without auth context")
	}
}
