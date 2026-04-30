package config

import (
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestLoadDepartmentServiceConfigs(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantMode string
	}{
		{name: "local", path: "../../etc/department-local.yaml", wantMode: "dev"},
		{name: "prod", path: "../../etc/department.yaml", wantMode: "pro"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			if err := conf.Load(tt.path, &cfg); err != nil {
				t.Fatalf("load config: %v", err)
			}

			if cfg.Mode != tt.wantMode {
				t.Fatalf("expected mode %q, got %q", tt.wantMode, cfg.Mode)
			}
			if cfg.Timeout != 2500 || cfg.CpuThreshold != 900 {
				t.Fatalf("unexpected rpc timeout/cpu: %d/%d", cfg.Timeout, cfg.CpuThreshold)
			}
			if !cfg.Health {
				t.Fatalf("expected rpc health check enabled")
			}
			if cfg.DevServer.Port != 6062 {
				t.Fatalf("expected dev server port 6062, got %d", cfg.DevServer.Port)
			}
			if !cfg.Middlewares.Recover || !cfg.Middlewares.Stat || !cfg.Middlewares.Prometheus || !cfg.Middlewares.Breaker {
				t.Fatalf("expected rpc server middlewares enabled")
			}
		})
	}
}
