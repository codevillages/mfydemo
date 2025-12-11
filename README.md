## 全局说明
- 最新版本指 v1.0.1 > v1.0.0，按字段序排就行了

## 执行需求
- 请先了解我的整体架构  @AI/infrastructure.md
- 请获取任务，在 @AI/需求 目录下的最新版本的{xx}_prd.md文档，这里面描述的是本次需要完整的需求
- 读取了需求之后，给出实现方案，实现方案相应的写入到 @AI/方案设计/vx.x.x./{xx}_solution.md中


## 任务分解
- 请先了解我的整体架构  @AI/infrastructure.md
- 在 @AI/方案设计/vx.x.x/{xxx}_solution.md中，是我们需求的方案设计，请将这些方案分解成执行计划，放到 @AI/执行计划/vx.x.x/{xx}_plan.md

## 执行任务
- 请先了解我的整体架构  @AI/infrastructure.md
- 在 @AI/执行计划/vx.x.x/{xxx}_plan.md中，是我们本次需要执行的任务，请将这些任务分解成功执行计划，放到 @AI/执行计划/vx.x.x/{xx}_plan.md

## 单元测试
- 请针对 service 层的public 方法做单元测试

## 集成测试
- 请模拟启动 grpc 服务，然后写一个 grpc client请求 grpc 接口，校验结果是否符合预期

