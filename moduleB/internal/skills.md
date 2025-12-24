# mfydemo-user-service v1.0.0 — 用户管理组件技能说明

## 标题/目的/适用场景
- 库名+版本：`mfydemo-user-service` v1.0.0（基于 GoFrame v2.9.0）
- 目的：提供用户管理 HTTP API（创建/列表/详情/更新/删除），包含密码 hash+salt、安全校验与统一响应。
- 推荐用在：中小型用户管理服务、后台管理系统、需要标准 CRUD 的内部 API。
- 替代方案：仅需轻量 HTTP 用 `net/http` 或 `gin`；强约束/生成式服务用 go-zero。
- 不适用：仅 gRPC 的纯 RPC 服务；高并发超大规模用户中心且需要分库分表的场景。

## 所需输入
- 配置文件：`config/config.yaml`
  - `database.default.link`：MySQL DSN（必填）
  - `database.bootstrap.*`：自动建库建表参数（必填）
- 环境：MySQL 服务可用（默认 `127.0.0.1:3306`）
- 推荐值：
  - `server.address`: `:8080`
  - `database.default.maxIdle`: 10
  - `database.default.maxOpen`: 100

## 流程/工作流程
1. 启动前调用 `bootstrap.EnsureDatabase(ctx)` 自动创建库表。
2. 初始化 DAO 与 Service：`mysql.NewUserDAO()` -> `service.NewUserService(dao)`。
3. 注册路由：在 `main.go` 的 `/api/v1` 分组下注册接口。
4. 请求处理流程：Controller -> Service -> DAO -> MySQL。
5. 响应输出：使用 `internal/response` 统一结构返回；HTTP 200 表示成功，非 200 表示失败。

## 何时使用该技能
- 需要快速搭建用户 CRUD HTTP 服务时。
- 需要统一响应结构与密码安全存储时。

## 输出格式
- HTTP 状态码：`200` 成功，非 `200` 失败。
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```
- 密码字段不返回（JSON `-`）。

## 示例
- 初始化与启动：
```go
if err := bootstrap.EnsureDatabase(context.Background()); err != nil {
	log.Fatalf("bootstrap failed: %v", err)
}
userDAO := mysql.NewUserDAO()
userService := service.NewUserService(userDAO)
userController := controller.NewUserController(userService)
```
- 常规操作：
  - 创建用户：`POST /api/v1/users`
  - 列表查询：`GET /api/v1/users`
  - 详情查询：`GET /api/v1/users/{id}`
- 错误/重试：仅对幂等读操作考虑重试，写操作不重试。
- 并发/连接管理：使用 `database.default.maxIdle/maxOpen` 控制连接池。
- 收尾清理：服务退出时调用 `s.Shutdown()`。

## 限制条件与安全规则
- 禁止记录密码明文；日志中不输出敏感字段。
- 写接口不做自动重试；读接口最多 3 次重试且总耗时 < 10s。
- 默认分页 `page_size=20`，防止大查询。
- 必须使用统一响应结构，不手写裸 JSON。

## 常见坑/FAQ
- 高：MySQL 未启动导致连接失败，确认 `127.0.0.1:3306` 可用。
- 中：未导入 GoFrame MySQL 驱动导致报错，需 `github.com/gogf/gf/contrib/drivers/mysql/v2`。
- 低：分页参数未传导致默认值生效，检查 `page/page_size`。

## 可观测性/诊断
- 通过 GoFrame 日志输出请求链路与错误。
- 关键字段：`status`、`latency`、`path`、`request_id`。

## 版本与依赖
- GoFrame: v2.9.0
- MySQL 驱动：`github.com/gogf/gf/contrib/drivers/mysql/v2` v2.9.0
- 依赖服务：MySQL
- 内部路径：`internal/controller`、`internal/service`、`internal/dao/mysql`

## 更新记录/Owner
- 最后更新时间：2025-12-18
- 维护人：@backend
- 评审人：@tech-lead
