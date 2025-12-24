# 用户管理组件需求规格说明

## 1. 背景与目标
- 基于 GoFrame 框架实现一个 HTTP Server 组件，提供用户管理的增删改查能力。
- 组件可复用、可扩展，支持后续接入权限、审计等能力。

## 2. 范围
- **包含**：用户创建、查询（列表与详情）、更新、删除。
- **不包含**：登录鉴权、权限控制、第三方账号体系、消息通知等（可在后续版本扩展）。

## 3. 角色与使用场景
- **管理员**：通过 API 管理用户数据（增删改查）。
- **系统集成方**：通过 HTTP 接口对接业务系统。

## 4. 数据模型（建议）
- id: int64，自增主键
- username: string，用户名（唯一）
- nickname: string，昵称
- email: string，邮箱（唯一，可选）
- phone: string，手机号（唯一，可选）
- password: string，密码（hash+salt 存储）
- status: int，状态（0=禁用，1=启用）
- created_at: datetime，创建时间
- updated_at: datetime，更新时间
- deleted_at: datetime，软删除时间（可选）

## 5. 功能性需求
### 5.1 创建用户
- 允许创建新用户，username 必填且唯一。
- 密码在创建时必填，持久化存储为 hash+salt，不保存明文。
- email、phone 可选，但若填写则需唯一。
- 默认 status=1（启用）。

### 5.2 查询用户
- **列表查询**：支持分页、关键词检索（username/nickname/email/phone）。
- **详情查询**：按 id 查询单个用户。
- 支持按 status 过滤。

### 5.3 更新用户
- 允许更新 nickname、email、phone、status。
- username 不允许修改（如需修改需走专门流程，当前不支持）。
- 更新时需校验 email、phone 的唯一性。

### 5.4 删除用户
- 支持软删除（推荐）与硬删除（可选配置）。
- 删除后列表不返回已删除用户，除非指定包含删除记录。

## 6. 接口设计（建议）
> 以 RESTful 风格为参考，可根据现有路由规范调整。

### 6.1 创建用户
- `POST /api/v1/users`
- 请求体：
  - username (string, required)
  - password (string, required)
  - nickname (string, optional)
  - email (string, optional)
  - phone (string, optional)
  - status (int, optional)
- 响应：创建成功返回用户详情

### 6.2 用户列表
- `GET /api/v1/users`
- 查询参数：
  - page (int, default=1)
  - page_size (int, default=20)
  - keyword (string, optional)
  - status (int, optional)
  - include_deleted (bool, optional)
- 响应：分页列表（total、list）

### 6.3 用户详情
- `GET /api/v1/users/{id}`
- 响应：用户详情

### 6.4 更新用户
- `PUT /api/v1/users/{id}`
- 请求体：
  - nickname (string, optional)
  - email (string, optional)
  - phone (string, optional)
  - status (int, optional)
- 响应：更新后的用户详情

### 6.5 删除用户
- `DELETE /api/v1/users/{id}`
- 查询参数（可选）：
  - hard (bool, default=false)
- 响应：删除结果

## 7. 校验与错误码（建议）
- 400: 参数校验失败
- 404: 用户不存在
- 409: username/email/phone 冲突
- 500: 服务内部错误

## 8. 非功能性需求
- **性能**：列表查询支持索引字段检索；分页默认 20。
- **安全**：密码以 hash+salt 方式存储，预留鉴权中间件接入点（如 JWT/Session）。
- **可观测性**：日志记录请求与异常信息；错误统一返回格式。

## 9. 约束与假设
- 组件基于 GoFrame（建议 ghttp 或 gf/v2 框架）。
- 采用统一响应结构（如 code/message/data），细节遵循现有项目规范。
- 数据存储为 MySQL，可通过 DAO/Model 抽象。

## 10. 输出物
- 用户管理组件代码（HTTP handlers、service、dao/model）。
- 接口文档与示例。
- 基础单元测试（可选）。
