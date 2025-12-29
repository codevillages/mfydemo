# 用户增删改查技术方案（module_user）

## 1. 需求与目标
- 基于 GoFrame v2.9.0（参考 `mfycommon/goframe/skills.md`）在 `main.go` 启动 HTTP Server，并加载 `module_user` 模块提供用户增删改查接口。
- 统一响应：成功 HTTP 200 + `{"code":0,"message":"OK","data":...}`；失败 HTTP 非 200 + 业务码（40001/40401/40901/50000）。
- 使用 MySQL 作为主存；软删除、唯一约束、密码仅存 bcrypt hash。

## 2. 模块与目录
- `main.go`：创建 `g.Server()`，加载配置/中间件，调用 `module_user.Init(s)` 注册路由与依赖。
- `module_user`：
  - `internal/controller/user.go`：HTTP Handler（绑定 `g.Meta`，参数校验，调用服务）。
  - `internal/service/user.go`：业务编排、事务控制、错误映射。
  - `internal/dao/mysql/user.go`：gdb 访问、唯一校验、软删。
  - `internal/model/entity/user.go`：实体/DTO 定义。
  - `internal/response`：统一响应封装（复用项目已有）。
  - `internal/middleware`：RequestID/AccessLog/HandlerResponse/Auth 占位。
  - `internal/bootstrap/ensure_database.go`：建库建表。

## 3. 启动与调用链
```
main.go
  -> bootstrap.LoadConfig() / InitLogger()
  -> s := g.Server()
  -> register middlewares (RequestID -> AccessLog -> MiddlewareHandlerResponse -> Auth placeholder)
  -> module_user.Init(s)
      -> controller.RegisterRoutes(group /api/v1/users)
      -> inject service/dao
  -> s.Run()
```
请求调用链（按成功路径）：
`HTTP Request -> Middleware(RequestID/AccessLog/HandlerResponse/Auth) -> controller.UserXxx -> service.UserXxx -> dao.UserXxx -> MySQL -> service -> controller -> HandlerResponse(JSON)`。

## 4. API 设计（前缀 `/api/v1/users`）
- 统一响应字段：`code`(int)、`message`(string)、`data`(object)。
- 公共 Headers：`X-Request-Id` 透传；`Content-Type: application/json`。

1) 创建用户 `POST /`
```json
Request: { "username": "u1", "password": "secret123", "passwordConfirm": "secret123", "email": "a@b.com", "phone": "13800000000", "nickname": "张三", "avatar": "https://..." }
Response 200: { "code": 0, "message": "OK", "data": { "id": 1 } }
```
校验：username 必填 4-32；password 长度≥8 且与 confirm 相同；email/phone 格式；唯一约束 username/email/phone。

2) 用户详情 `GET /{id}`
```json
Response 200: { "code": 0, "message": "OK", "data": { "id":1,"username":"u1","email":"a@b.com","phone":"138****0000","nickname":"张三","avatar":"https://...","status":0,"createdAt":"2025-01-01T10:00:00Z","updatedAt":"2025-01-02T10:00:00Z" } }
```

3) 用户列表 `GET /`
- Query: `page`(default 1)、`pageSize`(default 20, max 100)、`keyword`(模糊 username/nickname/email/phone)、`status`、`createdFrom`、`createdTo`、`sortBy`(id|created_at)、`order`(asc|desc)。
```json
Response 200: { "code":0,"message":"OK","data":{ "list":[{...}], "total":123, "page":1, "pageSize":20 } }
```

4) 更新用户 `PUT /{id}`
```json
Request: { "email":"new@a.com","phone":"13900000000","nickname":"张三2","avatar":"https://...","status":1,"password":"newsecret","passwordConfirm":"newsecret" }
Response 200: { "code":0,"message":"OK","data":{} }
```
仅在提供 password 时重置；唯一冲突返回 409/业务码 40901。

5) 删除用户 `DELETE /{id}`
```json
Response 200: { "code":0,"message":"OK","data":{} }
```
软删除，幂等。

6) 健康检查 `GET /healthz`
```json
Response 200: { "code":0,"message":"OK","data":{} }
```

错误响应示例：
```json
{ "code":40001, "message":"username is required", "data":{} }
```

状态/枚举：
- `status`: 0 正常（default），1 禁用。

## 5. 数据库设计（MySQL）
- 库：从 `config/config.yaml` 的 `database.default` 读取；启动 `bootstrap.EnsureDatabase` 自动建库建表。
- 表：`users`（InnoDB，utf8mb4_unicode_ci）
  - `id` BIGINT UNSIGNED PK AUTO_INCREMENT
  - `username` VARCHAR(32) UNIQUE NOT NULL
  - `email` VARCHAR(128) UNIQUE NULL
  - `phone` VARCHAR(16) UNIQUE NULL
  - `nickname` VARCHAR(64) NULL
  - `avatar` VARCHAR(255) NULL
  - `password_hash` VARCHAR(255) NOT NULL
  - `status` TINYINT NOT NULL DEFAULT 0
  - `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  - `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
  - `deleted_at` DATETIME NULL
- 索引：PK(id)，UNIQUE(username)，UNIQUE(email)，UNIQUE(phone)，索引(status, created_at)。
- 约束策略：唯一约束由 DB 保证；gdb 捕获唯一冲突并映射 409；所有查询默认 `deleted_at IS NULL`。

## 6. 中间件与安全
- 中间件顺序：RequestID -> AccessLog -> `ghttp.MiddlewareHandlerResponse` -> Auth 占位 -> 业务。
- 校验：使用 `g.Meta` + `v` 标签做必填/格式校验，失败返回 HTTP 400 + 业务码 40001。
- 密码：bcrypt cost 12-14；日志脱敏 password/passwordConfirm；返回数据隐藏 password_hash。
- CORS：默认关闭，按需开启白名单。
- 速率限制：预留中间件插槽，可基于 IP/用户。

## 7. 事务与并发
- 创建/更新在 service 层开启事务，包含唯一校验与写入；捕获唯一异常映射 409。
- 软删除通过更新 `deleted_at`；重复删除视为幂等成功。
- 列表分页：`LIMIT/OFFSET`；排序白名单控制。

## 8. 可观测性与运维
- AccessLog 字段：`request_id, method, path, status, latency_ms, client_ip, err`；>300ms 记 warn。
- 配置项：`server.address/readTimeout/writeTimeout/idleTimeout/maxHeaderBytes/clientMaxBodySize`；`server.logLevel/logStdout/accessLogEnabled/errorLogEnabled`。
- pprof/metrics：通过配置开关启用，默认关闭生产。

## 9. 交付物清单
- 代码：`module_user/internal/controller|service|dao|model|response|middleware|bootstrap`。
- 表初始化：`internal/bootstrap/ensure_database.go` 中的建表 SQL。
- 文档：`AI/v1.0.0/spec/baixs/user_spec.md`（需求）与本技术方案；后续可补充 OpenAPI/Swagger。
