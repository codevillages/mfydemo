## 系统背景
- 这是一个用户系统的微服务组件，主要提供grpc 接口，提供用户的增删改查，操作
- 用户信息的查询优先为了提速，将用户的状态放到redis 缓存中

## 系统设计模式
- 采用 DDD 设计模式
- 使用的是 mysql 数据库
- 微服务框架go-zero, gorm连接数据库


## 目录结构
```
./
├── AI
│   ├── infrastructure.md
│   ├── 需求
│   │   ├── v1.0.0
│   │   │   └── user_prd.md
│   │   └── v1.0.1
│   ├── 执行计划
│   │   ├── v1.0.0
│   │   │   └── user_plan.md
│   │   └── v1.0.1
│   └── 方案设计
│       ├── v1.0.0
│       │   └── user_solution.md
│       └── v1.0.1
├── README.md
├── main.go
├── mfy_components
│   ├── go-zero
│   │   ├── COORDINATION_PLAN.md
│   │   ├── INTEGRATION.md
│   │   ├── INTEGRATION_GUIDE.md
│   │   ├── README.md
│   │   ├── README_CN.md
│   │   ├── articles
│   │   │   └── ai-ecosystem-guide.md
│   │   ├── best-practices
│   │   │   └── overview.md
│   │   ├── getting-started
│   │   │   └── quick-start.md
│   │   ├── patterns
│   │   │   ├── database-patterns.md
│   │   │   ├── resilience-patterns.md
│   │   │   ├── rest-api-patterns.md
│   │   │   └── rpc-patterns.md
│   │   └── troubleshooting
│   │       └── common-issues.md
│   └── log
└── src
    ├── auth
    │   └── readme.md
    └── user
        └── readme.md
```

