# 日志


### 前言
日志服务是应用程序里不可缺少的一个模块， 在golang里也有不少比较优秀的日志服务框架，比如很多项目都是用的Zap，就是比较常用的golang的日志框架，己的项目里都是使用那个自己封装的log4go的日志包， 注意这里并没有从框架的模式上去进行实现， 所以我这里也称之为一个简洁的golang的日志包。


#### 为什么需要日志
调试开发
程序运行日志
用户行为日志
不同的目的决定了日志输出的格式、频率。作为开发人员，调试开发阶段打印日志目的是输出尽可能全的信息（如上下文，变量值...），辅助开发测试，因此日志格式要易读，打印频率要高。而在程序运行时，日志格式倾向于结构化（便于分析与搜索），而且为了性能和聚焦于关键信息（如error ），打印频率更偏低。

##### 以下介绍几种日志

1. log4go
2. zap
3. logrus