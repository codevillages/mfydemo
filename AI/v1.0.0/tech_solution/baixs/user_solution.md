# 用户管理组件技术方案

## 1. 实现的需求描述
本方案基于 `AI/v1.0.0/spec/baixs/user_spec.md`，使用 GoFrame 实现 HTTP Server 组件，提供用户管理增删改查能力，包含：
- 用户创建：用户名唯一、密码必填并以 hash+salt 存储，email/phone 可选且唯一。
- 用户查询：支持分页、关键词搜索、按状态过滤，支持详情查询。
- 用户更新：允许更新昵称、邮箱、手机号、状态；用户名不可修改。
- 用户删除：默认软删除，可选硬删除。
- 统一响应结构与错误码规范，HTTP 状态码 200 表示成功，非 200 表示失败，预留鉴权接入点。

## 2. 调用链流程图与函数调用流程

### 2.1 调用链流程图（文本版）
```
HTTP Request
   |
   v
Router (ghttp)
   |
   v
Controller (api/v1/user)
   |
   v
Service (user_service)
   |
   v
DAO/Model (user_dao)
   |
   v
MySQL
```

### 2.2 函数调用流程
以“创建用户”为例：
```
CreateUserHandler(ctx)
  -> userService.Create(ctx, req)
     -> userService.validateCreate(req)
     -> userService.hashPassword(req.password)
     -> userDao.Insert(ctx, userEntity)
  <- 返回 userEntity
<- 返回统一响应
```

以“用户列表”为例：
```
ListUsersHandler(ctx)
  -> userService.List(ctx, query)
     -> userService.buildQuery(query)
     -> userDao.SelectPage(ctx, filters)
  <- 返回 total + list
<- 返回统一响应
```

## 3. 数据库表设计（MySQL）

表名：`users`

```sql
CREATE TABLE `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL,
  `password` VARCHAR(255) NOT NULL COMMENT 'hash+salt',
  `nickname` VARCHAR(64) DEFAULT NULL,
  `email` VARCHAR(128) DEFAULT NULL,
  `phone` VARCHAR(32) DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  UNIQUE KEY `uk_users_email` (`email`),
  UNIQUE KEY `uk_users_phone` (`phone`),
  KEY `idx_users_status` (`status`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

说明：
- `password` 存储 hash+salt，不保存明文。
- email/phone 为可选字段，若为空需要注意唯一索引冲突策略（可改为组合索引或应用层控制）。
- 软删除使用 `deleted_at`，查询默认过滤非空数据。

## 4. API 接口设计

### 4.1 统一响应结构
JSON：
```
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```
XML：
```
<response>
  <code>0</code>
  <message>ok</message>
  <data></data>
</response>
```
Proto（示例）：
```
message ApiResponse {
  int32 code = 1;
  string message = 2;
  bytes data = 3; // 具体业务返回体序列化
}
```

### 4.2 接口列表

#### 3.2.1 创建用户
- `POST /api/v1/users`
- 请求 JSON：
```
{
  "username": "string",
  "password": "string",
  "nickname": "string",
  "email": "string",
  "phone": "string",
  "status": 1
}
```
- 响应 JSON：
```
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "username": "string",
    "nickname": "string",
    "email": "string",
    "phone": "string",
    "status": 1,
    "created_at": "2024-01-01 12:00:00",
    "updated_at": "2024-01-01 12:00:00"
  }
}
```

#### 3.2.2 用户列表
- `GET /api/v1/users`
- 查询参数：`page` `page_size` `keyword` `status` `include_deleted`
- 响应 JSON：
```
{
  "code": 0,
  "message": "ok",
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "username": "string",
        "nickname": "string",
        "email": "string",
        "phone": "string",
        "status": 1,
        "created_at": "2024-01-01 12:00:00",
        "updated_at": "2024-01-01 12:00:00"
      }
    ]
  }
}
```

#### 3.2.3 用户详情
- `GET /api/v1/users/{id}`

#### 3.2.4 更新用户
- `PUT /api/v1/users/{id}`
- 请求 JSON：
```
{
  "nickname": "string",
  "email": "string",
  "phone": "string",
  "status": 1
}
```

#### 3.2.5 删除用户
- `DELETE /api/v1/users/{id}`
- 查询参数：`hard`（默认 false）

### 4.3 字段与枚举说明
- status：`0=禁用`，`1=启用`
- include_deleted：`true=包含软删除`，`false=不包含`

### 4.4 XML 与 Proto 字段说明
XML 与 JSON 字段一致。
Proto 推荐定义：
```
message User {
  int64 id = 1;
  string username = 2;
  string nickname = 3;
  string email = 4;
  string phone = 5;
  int32 status = 6;
  string created_at = 7;
  string updated_at = 8;
}

message UserListResponse {
  int64 total = 1;
  repeated User list = 2;
}

enum UserStatus {
  USER_STATUS_DISABLED = 0;
  USER_STATUS_ENABLED = 1;
}
```
