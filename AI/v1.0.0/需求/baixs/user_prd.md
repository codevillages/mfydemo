
## 背景
- 我想实现一个微服务组件，提供 grpc 接口，实现用户的增删改查，目前项目从 0 开始
- 优先了解 @mfycommon 这些公共组件，能使用的优先使用

## 环境
- go version go1.23.9 darwin/arm64
- 利用到的公共组件版本以 @mfycommon下的文档为准
## 需求
- 用 DDD设计模式
- 帮我实现组件的目录搭建
- 提供用户的 grpc 增删改查接口
- 用户的status 状态信息需要缓存在 redis中，优先使用 @mfycommon/go-zero的redis 配置
- 日志请使用 @mfycommon/zap
- 框架使用使用 go-zero，grpc 优先使用 @mfycommon/go-zero的 rpc

## 接口
- AddUser 添加用户
- UpdateUser 修改用户
- RemoveUser 删除用户


## 参考
- @mfycommon/go-zero
- @mfycommon/zap


## 非功能性需求
- 不要用硬编码，尽量使用常量，一般情况我的常量都是配置在 @mos/tidb/model中的
- 日志经可能的简短，但是必须保证能查出问题，日志必须是英文，注释可以是中文或者英文
- 关键地方需要有 info 日志，不能放过任何错误日志

## 输出要求
- 要求输出每个函数的调用链
- 要求输出每一层封装哪些函数


请给我先输出方案 输出到 @AI/v1.0.0/方案设计/baixs/user_solution.md中