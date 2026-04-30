package config

import (
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestLoadGatewayConfigs(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		wantMode      string
		wantDevServer int
	}{
		{name: "local", path: "../../etc/apigateway-local.yaml", wantMode: "dev", wantDevServer: 6060},
		{name: "prod", path: "../../etc/apigateway.yaml", wantMode: "pro", wantDevServer: 6060},
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
			if cfg.Timeout != 3000 || cfg.CpuThreshold != 900 {
				t.Fatalf("unexpected rest timeout/cpu: %d/%d", cfg.Timeout, cfg.CpuThreshold)
			}
			if !cfg.Middlewares.Breaker || !cfg.Middlewares.Timeout || !cfg.Middlewares.Recover {
				t.Fatalf("expected gateway breaker/timeout/recover middleware enabled")
			}
			if cfg.DevServer.Port != tt.wantDevServer {
				t.Fatalf("expected dev server port %d, got %d", tt.wantDevServer, cfg.DevServer.Port)
			}
			assertRpcClientConfig(t, "user", cfg.UserRpc.Timeout, cfg.UserRpc.NonBlock, cfg.UserRpc.KeepaliveTime, cfg.UserRpc.Middlewares.Breaker, cfg.UserRpc.Middlewares.Timeout)
			assertRpcClientConfig(t, "department", cfg.DepartmentRpc.Timeout, cfg.DepartmentRpc.NonBlock, cfg.DepartmentRpc.KeepaliveTime, cfg.DepartmentRpc.Middlewares.Breaker, cfg.DepartmentRpc.Middlewares.Timeout)
		})
	}
}

func assertRpcClientConfig(t *testing.T, name string, timeout int64, nonBlock bool, keepalive time.Duration, breaker bool, timeoutMiddleware bool) {
	t.Helper()

	if timeout != 2000 {
		t.Fatalf("expected %s rpc timeout 2000, got %d", name, timeout)
	}
	if !nonBlock {
		t.Fatalf("expected %s rpc nonblock enabled", name)
	}
	if keepalive != 20*time.Second {
		t.Fatalf("expected %s rpc keepalive 20s, got %s", name, keepalive)
	}
	if !breaker || !timeoutMiddleware {
		t.Fatalf("expected %s rpc breaker/timeout middleware enabled", name)
	}
}
