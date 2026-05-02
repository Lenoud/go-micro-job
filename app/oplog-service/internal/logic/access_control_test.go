package logic

import (
	"context"
	"testing"

	"oplog-service/internal/common"
	"oplog-service/internal/model"
	"oplog-service/internal/svc"
	"oplog-service/oplog"
)

type fakeOpLogModel struct {
	batchCalled     bool
	opLogListCalled bool
	loginListCalled bool
	batchLogs       []*model.OpLog
}

func (m *fakeOpLogModel) BatchInsert(ctx context.Context, logs []*model.OpLog) error {
	m.batchCalled = true
	m.batchLogs = logs
	return nil
}

func (m *fakeOpLogModel) CountOpLogList(ctx context.Context) (int64, error) {
	return 1, nil
}

func (m *fakeOpLogModel) CountLoginLogList(ctx context.Context) (int64, error) {
	return 1, nil
}

func (m *fakeOpLogModel) FindOpLogList(ctx context.Context, page, pageSize int64) ([]*model.OpLog, error) {
	m.opLogListCalled = true
	return []*model.OpLog{{
		Id:          "10",
		RequestId:   "req-1",
		UserId:      "7",
		ReIp:        "127.0.0.1",
		RequestTime: 1710000000000,
		ReUa:        "ua",
		ReUrl:       "/api/user/create",
		ReMethod:    "POST",
		ReContent:   "username=a",
		Success:     "1",
		BizCode:     200,
		BizMsg:      "操作成功",
		ResponseMs:  12,
	}}, nil
}

func (m *fakeOpLogModel) FindLoginLogList(ctx context.Context, page, pageSize int64) ([]*model.OpLog, error) {
	m.loginListCalled = true
	return []*model.OpLog{{
		Id:         "11",
		RequestId:  "req-login",
		UserId:     "1",
		ReUrl:      "/api/user/login",
		ReMethod:   "POST",
		ReUa:       "login-ua",
		ResponseMs: 8,
	}}, nil
}

func (m *fakeOpLogModel) DeleteBefore(ctx context.Context, beforeTime int64) error {
	return nil
}

func TestListRejectsMissingAuthContext(t *testing.T) {
	t.Parallel()

	opLogModel := &fakeOpLogModel{}
	logic := NewListLogic(context.Background(), &svc.ServiceContext{OpLogModel: opLogModel})

	resp, err := logic.List(&oplog.OpLogListReq{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeForbidden {
		t.Fatalf("List() code = %v, want %v", resp.GetCode(), common.CodeForbidden)
	}
	if opLogModel.opLogListCalled {
		t.Fatal("List() queried operation logs without auth context")
	}
}

func TestListRejectsNonAdmin(t *testing.T) {
	t.Parallel()

	opLogModel := &fakeOpLogModel{}
	logic := NewListLogic(context.Background(), &svc.ServiceContext{OpLogModel: opLogModel})

	resp, err := logic.List(&oplog.OpLogListReq{
		Page:     1,
		PageSize: 10,
		Auth:     &oplog.OpLogContext{UserId: "2", Username: "hr", Role: common.RoleRecruiter},
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeForbidden {
		t.Fatalf("List() code = %v, want %v", resp.GetCode(), common.CodeForbidden)
	}
	if opLogModel.opLogListCalled {
		t.Fatal("List() queried operation logs for non-admin")
	}
}

func TestListAllowsAdminAndNormalizesPagination(t *testing.T) {
	t.Parallel()

	opLogModel := &fakeOpLogModel{}
	logic := NewListLogic(context.Background(), &svc.ServiceContext{OpLogModel: opLogModel})

	resp, err := logic.List(&oplog.OpLogListReq{
		Auth: &oplog.OpLogContext{UserId: "1", Username: "admin", Role: common.RoleAdmin},
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeSuccess {
		t.Fatalf("List() code = %v, want %v", resp.GetCode(), common.CodeSuccess)
	}
	if !opLogModel.opLogListCalled {
		t.Fatal("List() did not query operation logs for admin")
	}
	if resp.GetData().GetPage() != 1 || resp.GetData().GetPageSize() != 10 {
		t.Fatalf("List() page/pageSize = %d/%d, want 1/10", resp.GetData().GetPage(), resp.GetData().GetPageSize())
	}
	got := resp.GetData().GetList()[0]
	if got.GetReResponseTime() != "12" || got.GetReUserAgent() != "ua" {
		t.Fatalf("List() aliases = reResponseTime:%q reUserAgent:%q", got.GetReResponseTime(), got.GetReUserAgent())
	}
}

func TestLoginListUsesLoginFilter(t *testing.T) {
	t.Parallel()

	opLogModel := &fakeOpLogModel{}
	logic := NewLoginListLogic(context.Background(), &svc.ServiceContext{OpLogModel: opLogModel})

	resp, err := logic.LoginList(&oplog.OpLogListReq{
		Page:     1,
		PageSize: 10,
		Auth:     &oplog.OpLogContext{UserId: "1", Username: "admin", Role: common.RoleAdmin},
	})
	if err != nil {
		t.Fatalf("LoginList() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeSuccess {
		t.Fatalf("LoginList() code = %v, want %v", resp.GetCode(), common.CodeSuccess)
	}
	if !opLogModel.loginListCalled || opLogModel.opLogListCalled {
		t.Fatalf("LoginList() loginCalled=%v opCalled=%v", opLogModel.loginListCalled, opLogModel.opLogListCalled)
	}
	if got := resp.GetData().GetList()[0].GetReUserAgent(); got != "login-ua" {
		t.Fatalf("LoginList() reUserAgent = %q, want login-ua", got)
	}
}

func TestBatchCreateSkipsEmptyBatch(t *testing.T) {
	t.Parallel()

	opLogModel := &fakeOpLogModel{}
	logic := NewBatchCreateLogic(context.Background(), &svc.ServiceContext{OpLogModel: opLogModel})

	resp, err := logic.BatchCreate(&oplog.BatchCreateReq{})
	if err != nil {
		t.Fatalf("BatchCreate() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeSuccess {
		t.Fatalf("BatchCreate() code = %v, want %v", resp.GetCode(), common.CodeSuccess)
	}
	if opLogModel.batchCalled {
		t.Fatal("BatchCreate() inserted empty batch")
	}
}

func TestBatchCreatePersistsRecords(t *testing.T) {
	t.Parallel()

	opLogModel := &fakeOpLogModel{}
	logic := NewBatchCreateLogic(context.Background(), &svc.ServiceContext{OpLogModel: opLogModel})

	resp, err := logic.BatchCreate(&oplog.BatchCreateReq{
		Logs: []*oplog.OpLogRecord{{
			RequestId:  "req-1",
			UserId:     "1",
			ReIp:       "127.0.0.1",
			ReTime:     1710000000000,
			ReUa:       "ua",
			ReUrl:      "/api/user/create",
			ReMethod:   "POST",
			ReContent:  "username=a",
			Success:    "1",
			BizCode:    200,
			BizMsg:     "操作成功",
			AccessTime: 12,
		}},
	})
	if err != nil {
		t.Fatalf("BatchCreate() error = %v", err)
	}
	if resp == nil || resp.Code != common.CodeSuccess {
		t.Fatalf("BatchCreate() code = %v, want %v", resp.GetCode(), common.CodeSuccess)
	}
	if !opLogModel.batchCalled || len(opLogModel.batchLogs) != 1 {
		t.Fatalf("BatchCreate() batchCalled=%v logs=%d", opLogModel.batchCalled, len(opLogModel.batchLogs))
	}
	if got := opLogModel.batchLogs[0]; got.RequestId != "req-1" || got.ReUrl != "/api/user/create" {
		t.Fatalf("BatchCreate() stored log = %+v", got)
	}
}
