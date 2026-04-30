package logic

import (
	"context"
	"database/sql"
	"testing"

	"user-service/internal/common"
	"user-service/internal/model"
	"user-service/internal/svc"
	"user-service/user"
)

type fakeUserModel struct {
	listCalled bool
}

func (m *fakeUserModel) Insert(ctx context.Context, data *model.User) (sql.Result, error) {
	return nil, nil
}

func (m *fakeUserModel) FindOne(ctx context.Context, id string) (*model.User, error) {
	return nil, nil
}

func (m *fakeUserModel) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return nil, nil
}

func (m *fakeUserModel) FindAdminUser(ctx context.Context, username, password string) (*model.User, error) {
	return nil, nil
}

func (m *fakeUserModel) FindNormalUser(ctx context.Context, username, password string) (*model.User, error) {
	return nil, nil
}

func (m *fakeUserModel) FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*model.User, int64, error) {
	m.listCalled = true
	return []*model.User{}, 0, nil
}

func (m *fakeUserModel) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *fakeUserModel) DeleteBatch(ctx context.Context, ids []string) error {
	return nil
}

func (m *fakeUserModel) Update(ctx context.Context, data *model.User) error {
	return nil
}

func (m *fakeUserModel) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	return nil
}

func TestListRejectsMissingAuthContext(t *testing.T) {
	t.Parallel()

	userModel := &fakeUserModel{}
	logic := NewListLogic(context.Background(), &svc.ServiceContext{UserModel: userModel})

	resp, err := logic.List(&user.UserListReq{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeForbidden {
		t.Fatalf("List() code = %v, want %v", resp.GetCode(), common.CodeForbidden)
	}
	if userModel.listCalled {
		t.Fatal("List() queried users without auth context")
	}
}
