## 系统背景
- 这是一个用户管理 HTTP 服务组件，基于 GoFrame 提供用户增删改查接口。
- 统一响应结构，HTTP 200 表示成功，非 200 表示失败。
- 密码仅保存 hash+salt（bcrypt），不存明文。

## 系统设计模式
- 分层架构：Controller -> Service -> DAO -> MySQL。
- 使用 MySQL 数据库，GoFrame gdb 作为数据库访问层。
- 框架版本：GoFrame v2.9.0（详见 `mfycommon/goframe/skills.md`）。

## 配置与启动
- 配置文件：`config/config.yaml`。
  - `database.default.link`：应用连接 MySQL 的 DSN。
  - `database.bootstrap.*`：自动建库建表参数。
- 启动入口：`main.go`，启动时会执行 `bootstrap.EnsureDatabase` 自动创建库表。

## 代码规约
- 统一响应：使用 `internal/response` 输出 `{code,message,data}`。
- 错误处理：参数错误返回 400，未找到 404，冲突 409，系统错误 500。
- 密码处理：使用 `internal/security` 进行 hash+salt，不输出敏感字段。
- 接口命名：RESTful 风格，版本前缀 `/api/v1`。

## 中间件与横切关注点
- `internal/middleware/auth.go`：鉴权占位中间件，便于后续接入 JWT/Session。
- 日志：由 GoFrame 默认日志输出，后续可按需扩展访问日志与 tracing。

## 测试
- 集成测试：`internal/integration/user_create_test.go` 覆盖创建用户接口。

## 目录结构
```
./
├── AI                              # 需求/方案/任务文档入口
│   ├── v1.0.0                      # 当前版本文档
│   │   ├── spec/baixs              # 需求规格说明
│   │   ├── tech_solution/baixs     # 技术方案
│   │   └── tasks/baixs             # 任务拆解
│   └── v1.0.1                      # 历史/后续版本文档
├── config                          # 服务配置目录
│   └── config.yaml                 # MySQL 连接与初始化配置
├── internal                        # 业务实现主目录（不可对外导出）
│   ├── bootstrap                   # 启动引导与初始化（建库建表）
│   ├── controller                  # HTTP 控制器层（路由处理）
│   ├── dao                         # 数据访问层接口
│   │   └── mysql                   # MySQL DAO 实现
│   ├── integration                 # 集成测试
│   ├── middleware                  # 中间件（鉴权占位）
│   ├── model                       # 领域模型/实体
│   │   └── entity                  # 实体定义
│   ├── response                    # 统一响应结构与错误码
│   ├── security                    # 密码 hash+salt 等安全工具
│   └── service                     # 业务服务层
│   └── skills.md                   # 目录职责与使用说明
├── mfycommon                       # 规范/技能/示例库
│   └── goframe                     # GoFrame 规范与示例
├── go.mod                          # Go 模块定义
├── go.sum                          # Go 依赖校验
└── main.go                         # 应用启动入口
```

## AI 提示
- 每个目录下都会提供 `skills.md`，用于帮助 AI 快速了解该目录的职责与使用方式。
