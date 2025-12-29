# 用户增删改查任务拆解（module_user）

## 通用准备
- `main.go` 启动 Server，注册中间件（RequestID/AccessLog/HandlerResponse/Auth 占位），调用 `module_user.Init(s)` 注册 `/api/v1/users` 和 `/healthz`。
- `internal/bootstrap/ensure_database.go` 自动建库建表（`users` 字段与索引按技术方案）。
- 配置补全：`config/config.yaml` server/DB/超时/日志项；支持 `gf.gcfg.file/path`。

## 创建用户（POST /api/v1/users）
- DTO/校验：username 4-32 必填；password ≥8 且与 confirm 一致；email/phone 格式校验。
- Service：bcrypt hash(cost 12-14)，唯一校验映射 40901，事务内写入并返回 id。
- Controller：绑定路由/方法，调用 service，使用统一响应（200 code 0；失败按业务码）。
- 测试：必填校验、唯一冲突、成功创建返回 id。

## 查询用户详情（GET /api/v1/users/{id}）
- DAO：根据 id 且 `deleted_at IS NULL` 查询。
- Service：未找到映射 40401；输出脱敏（不含 password_hash，可对 phone/email 部分脱敏）。
- Controller：路由绑定与参数校验；统一响应。
- 测试：存在用户返回详情；不存在返回 40401。

## 用户列表查询（GET /api/v1/users）
- Query 处理：page/pageSize（默认 1/20，max 100）、keyword 模糊、status、createdFrom/To、sortBy(id|created_at)/order(asc|desc) 白名单。
- DAO：带条件的分页查询和 total 统计，过滤软删。
- Controller：参数校验，调用 service/dao，返回 `{list,total,page,pageSize}`。
- 测试：分页默认值，排序白名单，keyword 模糊，pageSize>100 被拒。

## 更新用户（PUT /api/v1/users/{id}）
- DTO/校验：允许 email/phone/nickname/avatar/status，可选 password+confirm；唯一约束校验。
- Service：事务内更新；密码字段时重新 bcrypt hash；唯一冲突映射 40901；未找到映射 40401。
- Controller：路由绑定、参数校验、统一响应。
- 测试：更新基础信息成功；密码更新；唯一冲突；不存在用户 40401。

## 删除用户（DELETE /api/v1/users/{id}）
- DAO/Service：软删（设置 `deleted_at`），重复删除幂等视为成功。
- Controller：路由绑定、统一响应。
- 测试：删除成功；重复删除仍返回成功。

## 健康检查与可观测性
- `GET /healthz` 返回 `{code:0,message:"OK",data:{}}`。
- AccessLog 含 request_id/method/path/status/latency_ms/client_ip；>300ms 记 warn。

## 验收清单
- API 行为符合需求/技术方案；成功 HTTP 200 code 0，失败 HTTP 非 200 + 业务码 40001/40401/40901/50000。
- 唯一冲突返回 40901；软删除幂等；密码不落日志。
- 自动建表生效；健康检查可用；列表分页/排序/过滤按白名单。
