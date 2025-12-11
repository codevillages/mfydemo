## 目标
- 基于 go-zero（优先 @mfycommon/go-zero）搭建用户 gRPC 微服务，接口包含 AddUser/UpdateUser/RemoveUser（如需查询可增补 GetUser），遵循 DDD，日志用 @mfycommon/zap，用户 status 缓存走 redis（优先用 @mfycommon/go-zero 的 redis 配置）。
- 输出可直接指导后续开发：目录、proto、方法签名、调用链、日志/错误/配置规范，避免硬编码（状态等常量由 @mos/tidb/model 包装后引用）。

## 快速落地步骤
1) 设计 proto：`proto/user.proto` 定义 Add/Update/Remove（可预留 Get 便于 CRUD 完整），字段 status 使用 pkg/constant 中对 @mos/tidb/model 的包装常量。
2) 用 goctl 生成基础骨架：`goctl rpc protoc proto/user.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./cmd/user.rpc --style=goZero`，确保依赖指向 @mfycommon/go-zero。
3) 搭建 DDD 目录并放空实现：domain/entity|repository|service，infra/dao|repository|cache，internal/logic|server|svc|config，pkg/constant。
4) 初始化配置与依赖：config.go 引入 zap、redis、db、ttl；service_context 注入 zap logger、redis client、dao、repo。
5) 实现 repo + cache + service + logic，补日志/错误处理；最后补测试（repo/cache/logic、端到端）。

## 目录结构（DDD）
- `proto/user.proto`
- `cmd/user.rpc`：main/注册/启动。
- `internal/config/config.go`：`RpcServerConf`, `LogConf`(@mfycommon/zap), `RedisConf`, `DBConf`, `CacheTTL`, `Etcd`.
- `internal/svc/service_context.go`：初始化 logger、redis、dao、repo。
- `internal/server/user_server.go`：gRPC 入口，调 logic。
- `internal/logic/user/*.go`：Add/Update/Remove（可含 Get）。
- `internal/domain/entity/user.go`：ID, Name, Email, Status, Version, CreatedAt, UpdatedAt。
- `internal/domain/repository/user_repository.go`：接口定义。
- `internal/domain/service/user_service.go`：业务校验与编排。
- `internal/infra/dao/user_model.go`：go-zero model（Insert/FindOne/Update/Delete）。
- `internal/infra/repository/user_repo.go`：实现仓储，聚合 dao+cache。
- `internal/infra/cache/user_status_cache.go`：redis 封装。
- `pkg/constant/status.go`：包装 @mos/tidb/model 常量，统一 key/错误码。

## Proto 建议
- Messages:
  - `AddUserReq { string name; string email; int32 status; }` -> `AddUserResp { int64 id; }`
  - `UpdateUserReq { int64 id; string name; string email; int32 status; int64 version; }` -> `google.protobuf.Empty`
  - `RemoveUserReq { int64 id; }` -> `google.protobuf.Empty`
  - 可选 `GetUserReq { int64 id; }` -> `GetUserResp { User user; }`（满足“查”需求）
  - `User { int64 id; string name; string email; int32 status; int64 version; int64 created_at; int64 updated_at; }`
- Service: `rpc AddUser`, `UpdateUser`, `RemoveUser`（如需要查询则加 `GetUser`）。

## 配置示例（config.yaml 提示）
```yaml
RpcServerConf:
  ListenOn: 0.0.0.0:9001
Etcd:
  Hosts: [127.0.0.1:2379]
  Key: user.rpc
LogConf: # @mfycommon/zap
  Level: info
RedisConf: # @mfycommon/go-zero redis
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""
CacheTTL:
  UserStatusSeconds: 600
DBConf:
  DataSource: user:pwd@tcp(127.0.0.1:3306)/userdb?charset=utf8mb4&parseTime=true
```

## 域模型与仓储接口
```go
type UserRepository interface {
    Save(ctx context.Context, u entity.User) (int64, error)
    Update(ctx context.Context, u entity.User) error
    Delete(ctx context.Context, id int64) error
    Get(ctx context.Context, id int64) (*entity.User, error) // 便于 Remove 前校验/查
    SetStatusCache(ctx context.Context, id int64, status int32, version int64, ttl time.Duration) error
    GetStatusCache(ctx context.Context, id int64) (*StatusCache, error)
    DelStatusCache(ctx context.Context, id int64) error
}
type StatusCache struct {
    Status  int32
    Version int64
}
```
- `User` 字段与 proto 对齐；状态常量从 pkg/constant/status.go。

## 领域服务（internal/domain/service/user_service.go）
- `Add(ctx, entity.User) (int64, error)`: 校验必填/状态合法/邮箱格式 -> repo.Save -> SetStatusCache。
- `Update(ctx, entity.User) error`: 先校验存在与版本，repo.Update（乐观锁）-> SetStatusCache。
- `Remove(ctx, id int64) error`: 校验存在 -> repo.Delete -> DelStatusCache。
- 可选 `Get(ctx, id int64) (*entity.User, error)`: 读缓存 status/version，查 DB 补全，miss 回填。

## 缓存策略
- key: `user:status:{id}`；value: JSON/struct `{status,version}`；TTL 从 config。
- 写成功后刷新缓存；删除后立刻删缓存；读优先缓存，miss 回源再回填。
- 并发：Update 时带 version，缓存写入前比对版本，防止旧数据覆盖。

## 事务与并发
- Add/Update/Remove 使用 go-zero sqlx Session 保证 DB 一致性；缓存更新在事务成功后执行。
- Update WHERE id=? AND version=?，影响行数为 0 时返回 ErrVersionConflict。

## 日志与错误
- 全英文日志，关键链路 info：请求入参（id/status/version）与结果；错误处 error 记录 err、trace id。
- server 层收敛 gRPC error；logic 记录业务意图并映射错误码；repo/dao 返回具体错误。
- 常量化错误码/状态/redis key 前缀，放 pkg/constant（包装 @mos/tidb/model，避免硬编码）。

## 调用链（含层次）
- AddUser: `server.AddUser -> logic.AddUser -> domain.UserService.Add -> repo.Save -> dao.Insert -> repo.SetStatusCache -> cache.SetStatus`.
- UpdateUser: `server.UpdateUser -> logic.UpdateUser -> domain.UserService.Update -> repo.Update -> dao.Update -> repo.SetStatusCache -> cache.SetStatus`.
- RemoveUser: `server.RemoveUser -> logic.RemoveUser -> domain.UserService.Remove -> repo.Delete -> dao.Delete -> repo.DelStatusCache -> cache.DelStatus`.
- 可选 GetUser: `server.GetUser -> logic.GetUser -> domain.UserService.Get -> repo.GetStatusCache? -> repo.Get -> dao.FindOne -> repo.SetStatusCache(miss)`。

## 每层封装的函数
- server: gRPC Handler（AddUser/UpdateUser/RemoveUser/可选 GetUser）。
- logic: `AddUser`, `UpdateUser`, `RemoveUser`, `GetUser(可选)`，负责 DTO<->实体、记录日志、错误映射。
- domain service: `Add`, `Update`, `Remove`, `Get(可选)`。
- repository impl: `Save`, `Update`, `Delete`, `Get`, `SetStatusCache`, `GetStatusCache`, `DelStatusCache`。
- dao: `Insert`, `FindOne`, `Update`, `Delete`（乐观锁），必要时 `WithSession`。
- cache: `SetStatus`, `GetStatus`, `DelStatus`。
- pkg/constant: 状态/错误码/redis key 前缀。

## 测试指引
- repo/cache 单测：mock dao/redis，验证缓存命中/穿透/回填、版本写保护。
- logic 单测：mock svcCtx/repo，断言错误码、日志字段是否包含 id/status。
- 集成测试：启动 rpc，调用 Add/Update/Remove（及可选 Get），验证 redis 缓存写入与删除、日志输出。
