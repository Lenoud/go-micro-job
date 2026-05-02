package svc

import (
	"context"
	"sync"
	"time"

	"oplog-service/internal/config"
	"oplog-service/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	OpLogModel  model.OpLogModel
	stopCleanup chan struct{}
	stopOnce    sync.Once
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	ctx := &ServiceContext{
		Config:      c,
		OpLogModel:  model.NewOpLogModel(conn),
		stopCleanup: make(chan struct{}),
	}
	ctx.startCleanupLoop()
	return ctx
}

func (s *ServiceContext) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCleanup)
	})
}

func (s *ServiceContext) startCleanupLoop() {
	retentionDays := s.Config.RetentionDays
	if retentionDays <= 0 {
		retentionDays = 90
	}

	go func() {
		s.cleanup(retentionDays)
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.cleanup(retentionDays)
			case <-s.stopCleanup:
				return
			}
		}
	}()
}

func (s *ServiceContext) cleanup(retentionDays int64) {
	before := time.Now().AddDate(0, 0, -int(retentionDays)).UnixMilli()
	if err := s.OpLogModel.DeleteBefore(context.Background(), before); err != nil {
		logx.Errorf("[oplog-service] cleanup old logs failed: %v", err)
	}
}
