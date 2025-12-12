## 背景 / 目标 / 非目标
- 背景：从 0 开始实现一个用户微服务组件，提供 gRPC 用户增删改查。
- 目标：交付 DDD 分层目录、gRPC 接口定义、数据库/缓存设计、配置样例、日志方案、迁移脚本与实现思路。
- 非目标：鉴权、多租户、审计、灰度发布暂不在本迭代。

## 依赖与版本
- Go: 1.23.9 darwin/arm64（go.mod 需声明）。
- 框架：@mfycommon/go-zero（使用 rpc、sqlx、redis 配置；版本按 @mfycommon 最新稳定版，需在方案中明确）。
- 日志：@mfycommon/zap（版本按 @mfycommon 最新稳定版，需在方案中明确）。
- 密码哈希：优先 bcrypt（若 @mfycommon 有推荐版本，按推荐；否则使用成熟版本并注明）。
- 迁移工具：goose 或 migrate（二选一并说明版本）。

## 需求（功能）
- 使用 DDD 设计模式，基于 goctl rpc 生成骨架并调整为 DDD 分层。
- 提供用户 gRPC 增删改查接口：AddUser、UpdateUser、RemoveUser、GetUser。
- 用户 status 状态信息需要缓存到 Redis，优先使用 @mfycommon/go-zero 的 redis 配置。
- 日志必须使用 @mfycommon/zap，输出英文、JSON 格式。
- 框架使用 go-zero，gRPC 使用 @mfycommon/go-zero 的 rpc。

## 接口与错误码
- gRPC Service：UserService
- 方法：AddUser、UpdateUser、RemoveUser、GetUser
- 请求/响应字段：需在 proto 中定义，时间字段使用 RFC3339 字符串。
- gRPC status 映射：InvalidArgument（校验失败）、AlreadyExists（用户名重复）、NotFound（用户不存在）、Internal（未知错误）。

## 数据库设计
- 表：users
- 字段：id, username, password(hash), email, status, created_at, updated_at
- 需明确类型、约束、默认值、索引；说明软删/硬删策略（默认硬删），写操作的事务策略。

## 环境
- go version go1.23.9 darwin/arm64
- 公共组件版本以 @mfycommon 文档为准；若未指定，使用最新稳定版并在方案中明确版本号。

## 参考
- 使用 @mfycommon/go-zero
- 使用 @mfycommon/zap

## 目录/分层约定（DDD + go-zero）
- 建议目录：
  - internal/app/user/api：proto、生成 pb、server 启动（goctl rpc new user 生成后调整）。
  - internal/app/user/handler：gRPC handler，请求校验与 DTO 转换。
  - internal/app/user/usecase：应用服务，编排事务/缓存/日志。
  - internal/domain/user：实体、领域常量（status 枚举、redis key 前缀、错误）、领域校验。
  - internal/domain/user/repo.go：UserRepo（DB）、UserCacheRepo（Redis）接口。
  - internal/infra/persistence：UserRepo 实现（go-zero sqlx），mapper。
  - internal/infra/cache：UserCacheRepo 实现（Redis）。
  - internal/pkg/logger：@mfycommon/zap 初始化封装，与 go-zero logx 适配。
  - configs：go-zero rpc 配置（service/db/redis/log）。
  - migrations：数据库迁移文件。
  - scripts：goctl 生成、迁移、启动脚本。
- 职责边界：handler 只做校验/转换；usecase 编排；domain 定义规则；repo/cache 实现接口。

## 缓存策略
- 缓存对象：用户 status。
- key 规范：`user:status:{id}`，值包含 status、updated_at；TTL 可配置（如 24h）。
- 策略：写操作完成后写/删缓存；写失败不影响主流程但需 error 日志；Get 先查缓存，miss 回源 DB 并回写缓存。

## 日志
- 使用 @mfycommon/zap，JSON 输出，英文消息。
- 日志点：入口/关键分支/外部 IO 前后使用 info；所有错误必须 error。
- 日志字段：request_id、user_id、username、err。
- 注释可中英，日志必须英文。

## 配置化与常量化
- 禁止硬编码：status 枚举、redis key 前缀、错误码、日志字段、连接信息。
- 配置文件：YAML/JSON，路径如 configs/user.yaml，包含 rpc、db、redis、log。

## 安全要求
- 密码只存 hash，不输出到日志；参数化 SQL；输入校验（用户名长度/字符集、密码最小长度>=8、email 格式）。
- TLS/鉴权预留（方案中注明扩展位）。

## 非功能性需求
- 不要用硬编码，使用常量。
- 日志尽可能简短但可定位问题，必须英文；关键路径 info，不放过错误日志。
- 考虑鲁棒性、简洁性、可维护性、可扩展性、可读性、性能、可测试性、安全性、兼容性、可复用性。

## 输出要求
- 整体方案调用链（若有多个，分多个链给出）。
- 技术方案按接口/需求逐项展开：包含 gRPC proto 定义、每层实现逻辑、字段取值说明。
- DB 方案需先给 schema 设计；接口方案需给字段定义；关联模块需列出；注意事项需列出。
- 明确依赖版本、目录结构、执行顺序（版本确认 -> goctl 生成骨架 -> 目录调整 -> 迁移/配置 -> repo/cache/usecase/handler -> 测试）。

请输出方案到 AI/v1.0.0/tech_solution/baixs/user_solution.md
