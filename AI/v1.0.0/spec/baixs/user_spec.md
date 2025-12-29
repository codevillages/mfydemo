# 用户增删改查 HTTP 服务需求说明

## 背景与目标
- 基于 GoFrame v2.9.0（参考 `mfycommon/goframe/skills.md`）提供用户增删改查 HTTP 组件，可直接作为独立服务启动。
- 输出统一响应 `{code,message,data}`，对接后台管理或内部 API；后续可挂接鉴权/审计。
- 默认使用 MySQL（gdb）和内置日志，配置化端口/超时，支持容器化部署。
- 入口 `main.go` 启动 HTTP Server，并按模块加载（例如 `module_user`），该模块暴露用户路由与依赖初始化。

## 作用域
- **功能内**：用户创建/查询（单个与列表）/更新/删除；密码安全存储；基本查询过滤；软删除占位。
- **功能外**：注册/登录、权限体系、短信/邮件、三方 OAuth、批量导入导出。

## 领域模型与约束
- 表：`users`（推荐 InnoDB，utf8mb4）。
- 字段：`id`(bigint PK, 自增)、`username`(唯一, 4-32)、`email`(唯一, 可空)、`phone`(唯一, 可空, 国内 11 位)、`nickname`、`avatar`、`password_hash`(bcrypt)、`status`(0 正常/1 禁用)、`created_at/updated_at/deleted_at`。
- 约束：username/email/phone 唯一；删除默认软删（`deleted_at` 置值）；禁用用户不允许登录（预留）。
- 密码：仅存 bcrypt hash，不回传给客户端；更新密码需二次确认字段。

## API 规格（前缀 `/api/v1/users`）
- **创建** `POST /`：请求体 `{username, password, email?, phone?, nickname?, avatar?}`；返回 `id`。校验 username 必填且唯一、password 长度 ≥8，email/phone 格式。
- **详情** `GET /{id}`：路径参数 `id`；返回用户基础字段（不含密码 hash）。
- **列表** `GET /`：查询参数 `page`(默认1)、`pageSize`(默认20, ≤100)、`keyword`(模糊匹配 username/nickname/email/phone)、`status`、`createdFrom/createdTo`；返回分页 `{list,total,page,pageSize}`。
- **更新** `PUT /{id}`：允许更新 `email/phone/nickname/avatar/status`，可选 `password`（仅在提供时重置）；保持唯一约束。
- **删除** `DELETE /{id}`：软删除；重复删除幂等返回成功。
- 统一响应：成功 `{"code":0,"message":"OK","data":{...}}`；错误码见下。
- Headers：`X-Request-Id` 透传；`Content-Type: application/json`。

## 错误码与 HTTP 状态
- 成功：HTTP 200 + 业务码 `0`。
- 失败：HTTP 非 200（建议：400 参数错误/校验失败 -> 业务码 `40001`；404 未找到 -> 业务码 `40401`；409 唯一冲突 -> 业务码 `40901`；500 内部错误 -> 业务码 `50000`）。
- message 采用 gerror 文案；响应体保持 `{code,message,data}`。

## 配置与启动
- 配置文件：`config/config.yaml`（或通过 `gf.gcfg.file/path` 覆盖），关键项：
  - `server.address`、`readTimeout`、`writeTimeout`、`idleTimeout`、`maxHeaderBytes`、`clientMaxBodySize`。
  - `server.logLevel/logStdout/accessLogEnabled/errorLogEnabled`。
  - `database.default.*`：MySQL 连接；`queryTimeout/execTimeout` 按 1-3s 设定。
- 启动：`main.go` 创建 `s := g.Server()`，挂载中间件（RequestID -> AccessLog -> `ghttp.MiddlewareHandlerResponse` -> Auth 占位），`group := s.Group("/api/v1")` 注册用户路由；`s.SetGraceful(true)` 优雅退出。
- 依赖：Go 1.22，GoFrame v2.9.0，bcrypt（`golang.org/x/crypto/bcrypt`），可选 Redis 缓存用户详情。

## 校验与安全
- 参数校验：使用 `g.Meta` + `v` 标签，失败返回 400；email 使用正则，phone 使用 11 位数字校验。
- 密码：创建/重置必须两次确认（`password`+`passwordConfirm`）；bcrypt `cost` 12-14。
- 日志脱敏：日志中屏蔽 password/passwordConfirm，隐藏手机号/邮箱中间位。
- CORS：默认关闭，如需对接前端可开启白名单。
- 速率限制：预留中间件插槽，可按 IP/用户限流。

## 数据与并发一致性
- 数据库层使用事务保障创建与唯一校验一致性；遇到唯一约束错误映射为 409。
- 软删除查询默认过滤 `deleted_at IS NULL`；需要包含已删记录的接口暂不支持。
- 列表查询排序：默认 `created_at DESC`，允许通过 `sortBy=created_at|id` + `order=asc|desc` 白名单。

## 可观测性与运维
- AccessLog 含 `request_id/method/path/status/latency_ms/client_ip`；慢请求 >300ms 记 `warn`。
- 健康检查：`GET /healthz`，返回 `{code:0,message:"OK"}`。
- pprof/metrics：可通过配置开关（如 `server.pprofEnabled`）开启，默认关闭生产。

## 交付物
- 代码目录（参考 `AI/v1.0.0/infra.md`）：`internal/controller/user.go`、`internal/service/user.go`、`internal/dao/mysql/user.go`、`internal/model/entity/user.go`、`internal/response`、`internal/middleware`。
- 初始化脚本：`internal/bootstrap/ensure_database.go` 创建 `users` 表及唯一索引。
- 文档：本文件 + `skills.md`（目录职责说明）+ OpenAPI/Swagger 占位。
