# 用户管理组件实现任务拆解

> 依据 `AI/v1.0.0/tech_solution/baixs/user_solution.md`，按接口维度拆分实现步骤，不生成代码。

## 0. 通用与基础准备（所有接口共享）
- 定义统一响应结构与错误码：HTTP 200 为成功，非 200 为失败。
- 定义用户实体与数据访问层：字段、索引约束、软删除约定。
- 定义密码处理工具：hash+salt 生成与校验策略（仅保存 hash+salt）。
- 定义路由组与版本前缀：`/api/v1`，预留鉴权中间件接入点。
- 定义输入校验与错误返回规范：参数必填、唯一性冲突、未找到等。

## 1. 创建用户 `POST /api/v1/users`
- 设计请求/响应 DTO：包含 username、password、nickname、email、phone、status。
- 校验逻辑：username 必填且唯一；password 必填；email/phone 若填写需唯一。
- 密码处理：将明文 password 转为 hash+salt 存储。
- 持久化：写入 MySQL，处理唯一索引冲突并返回冲突错误。
- 返回结果：HTTP 200 + 用户详情（不返回 password）。

## 2. 用户列表 `GET /api/v1/users`
- 设计查询参数：page、page_size、keyword、status、include_deleted。
- 构建查询条件：关键词匹配 username/nickname/email/phone；状态过滤；软删除过滤。
- 分页查询：返回 total + list。
- 返回结果：HTTP 200 + 分页列表（不返回 password）。

## 3. 用户详情 `GET /api/v1/users/{id}`
- 路径参数校验：id 为正整数。
- 查询逻辑：按 id 获取用户，默认排除软删除数据。
- 未找到处理：返回非 200 状态码。
- 返回结果：HTTP 200 + 用户详情（不返回 password）。

## 4. 更新用户 `PUT /api/v1/users/{id}`
- 设计请求 DTO：nickname、email、phone、status（username 不可改）。
- 校验逻辑：email/phone 若更新需唯一性检查。
- 更新逻辑：仅更新允许字段，记录更新时间。
- 未找到处理：返回非 200 状态码。
- 返回结果：HTTP 200 + 更新后的用户详情（不返回 password）。

## 5. 删除用户 `DELETE /api/v1/users/{id}`
- 路径参数校验：id 为正整数。
- 删除策略：默认软删除；`hard=true` 时执行硬删除。
- 删除结果：未找到返回非 200；成功返回 HTTP 200。
- 兼容列表：软删除后列表默认不返回该用户。
