## 可执行小任务清单（按序）
1. 初始化 Go 模块与基础目录：创建 `go.mod`（module 建议 `mfyai/mfydemo`），新增 `proto/`、`cmd/user.rpc/`、`internal/{config,svc,server,logic,domain,infra}/`、`pkg/constant/` 等目录。
2. 编写 `proto/user.proto`（含 AddUser/UpdateUser/RemoveUser/可选 GetUser），运行 `protoc --go_out=. --go-grpc_out=.` 生成 pb 代码。
3. 定义配置结构：`internal/config/config.go` 包含 rpc/etcd/redis/db/cacheTTL/log；编写示例 `cmd/user.rpc/etc/user.yaml`。
4. 搭建领域层：`internal/domain/entity/user.go`、`internal/domain/repository/user_repository.go`、`internal/domain/service/user_service.go`（含校验、乐观锁预期）。
5. 搭建基础设施层：`internal/infra/dao/user_model.go`（go-zero sqlx model & CRUD）、`internal/infra/cache/user_status_cache.go`、`internal/infra/repository/user_repo.go`（聚合 dao+cache+事务）。
6. 搭建接口层：`internal/svc/service_context.go` 初始化 zap、redis、dao、repo；`internal/logic/user/*.go` 实现用例；`internal/server/user_server.go` 注册 gRPC。
7. 编写入口 `cmd/user.rpc/user.go` 启动 zrpc 服务（etcd 注册）并引用 server。
8. 编写 `pkg/constant/status.go`、错误码/redis key 常量（包装 @mos/tidb/model）。
9. 补充单元测试：仓储缓存逻辑/领域服务；（若时间不足，至少保证核心逻辑可编译）。
10. 运行构建/启动：`go test ./...` 或 `go run cmd/user.rpc/user.go -f cmd/user.rpc/etc/user.yaml`，修复报错至运行通过。
