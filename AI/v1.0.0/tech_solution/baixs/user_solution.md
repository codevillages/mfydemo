## 调用链（整体）
- Client -> gRPC（@mfycommon/go-zero rpc） -> Handler（仅校验/DTO 转换） -> Usecase（应用服务，编排逻辑/事务/日志） -> Domain（实体/规则） -> Repository 接口 -> Infra DB Repo（MySQL via go-zero sqlx）
- Usecase -> Cache Repo（Redis via @mfycommon/go-zero 配置）存取用户 status
- Usecase -> Logger（@mfycommon/zap）在入口、关键分支、错误、外部 IO 前后记录英文日志

## 目录/工程规划（DDD + go-zero）
- `internal/app/user/api`：proto、生成 pb、server 启动（goctl rpc new user 后调整）。
- `internal/app/user/handler`：gRPC handler，请求校验与 DTO 转换。
- `internal/app/user/usecase`：应用服务（Add/Update/Remove/Get），组织事务、缓存、日志。
- `internal/domain/user`：User 聚合、领域常量（status 枚举、redis key 前缀、错误）、领域校验。
- `internal/domain/user/repo.go`：`UserRepo`（DB）、`UserCacheRepo`（Redis）接口。
- `internal/infra/persistence`：UserRepo 实现（go-zero sqlx），mapper。
- `internal/infra/cache`：UserCacheRepo 实现（Redis）。
- `internal/pkg/logger`：@mfycommon/zap 初始化封装，与 go-zero logx 适配。
- `configs`：go-zero rpc 配置（service、db、redis、log）。
- `migrations`：用户表 schema 及变更。
- `scripts`：goctl 生成、迁移、启动脚本。

## 数据库 Schema 设计
- 表：`users`
  - `id` BIGINT UNSIGNED PK AUTO_INCREMENT
  - `username` VARCHAR(64) UNIQUE NOT NULL
  - `password` VARCHAR(128) NOT NULL（存 hash）
  - `email` VARCHAR(128) NOT NULL
  - `status` TINYINT NOT NULL（0=inactive,1=active,2=blocked，常量定义）
  - `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  - `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
- 索引：`uk_username` 唯一，`idx_status`，可选 `idx_email`
- 事务：写操作走事务；更新/删除后更新或删除缓存；失败回滚

## gRPC 接口设计
- Service：`user.UserService`
- 请求/响应：
  - `AddUserRequest { string username; string password; string email; int32 status; }`
  - `AddUserResponse { int64 id; }`
  - `UpdateUserRequest { int64 id; string password; string email; int32 status; }`
  - `UpdateUserResponse { bool success; }`
  - `RemoveUserRequest { int64 id; }`
  - `RemoveUserResponse { bool success; }`
  - `GetUserRequest { int64 id; }`
  - `GetUserResponse { int64 id; string username; string email; int32 status; string created_at; string updated_at; }`
- 错误映射：InvalidArgument、AlreadyExists、NotFound、Internal
- 时间格式：响应时间字段使用 RFC3339 字符串

## 需求拆解与方案
- 目录搭建：goctl 生成 rpc 骨架后按 DDD 目录调整；常量集中定义，禁止硬编码。
- DDD 设计：实体 `User`（ID, Username, PasswordHash, Email, Status, CreatedAt, UpdatedAt）；领域常量（status 值、redis key 前缀、日志字段、错误码）；领域校验（用户名/邮箱格式、密码长度>=8、status 合法）。
- gRPC CRUD：
  - AddUser：校验 -> 检查 username 唯一 -> hash 密码（bcrypt）-> DB 插入 -> 写 Redis status -> 返回 id；入口/成功 info，错误 error。
  - UpdateUser：校验 -> 查存在 -> 更新邮箱/密码/status -> DB 提交 -> 写 Redis status；缓存写失败 error 但不阻断。
  - RemoveUser：硬删；事务删除 -> 删除 Redis status；不存在返回 NotFound；日志覆盖入口/成功/错误。
  - GetUser：先读 Redis status，命中则 DB 查基础信息并返回；miss 时查 DB，回写缓存 status；错误记日志。
- 缓存策略：key `user:status:{id}`；值 `{status:int, updated_at:timestamp}`；TTL 配置（如 24h）；缓存写失败不影响主流程但需 error log。
- 日志：@mfycommon/zap JSON；入口 info，关键状态变更 info，错误 error，字段含 err/user_id/username/request_id；日志英文。
- 配置化：RPC、DB、Redis、Log 写在配置文件；常量集中 `internal/constants`；避免硬编码。

## 关联模块
- @mfycommon/go-zero：rpc 框架、sqlx、redis、配置结构、logx 适配
- @mfycommon/zap：日志初始化与格式，go-zero logx bridge
- bcrypt：密码哈希（若 @mfycommon 有推荐版本则使用）
- goose/migrate：迁移执行

## 注意事项
- 安全：只存密码 hash，不记录密码；参数化 SQL 防注入；gRPC TLS/鉴权后续扩展
- 常量化：status、redis key 前缀、错误码、日志字段全部常量
- 性能/扩展：缓存降低 DB 读；连接池配置化；预留批量接口扩展
- 可测试性：Usecase/Repo/Cache 接口可 mock；单测覆盖校验、缓存命中/回源、错误分支；gRPC 集成测试覆盖 CRUD

## 实施步骤（建议）
- 确认 @mfycommon/go-zero 与 @mfycommon/zap 版本，生成 goctl rpc 骨架
- 编写迁移脚本与配置模板（db/redis/log/service）
- 实现 repo/cache/usecase/handler，补充单元/集成测试
