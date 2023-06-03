# zap


zap日志库
在许多Go语言项目中，我们需要一个好的日志记录器能够提供下面这些功能:

能够将事件记录到文件中，而不是应用程序控制台;
日志切割-能够根据文件大小、时间或间隔等来切割日志文件;
支持不同的日志级别。例如INFO，DEBUG，ERROR等;
能够打印基本信息，如调用文件/函数名和行号，日志时间等;


比较全的日志级别

支持结构化日志

快速的、结构化的、分层次的日志记录器


#### 使用
```
go get -u go.uber.org/zap
```

Zap提供了两种类型的日志记录器 — Sugared Logger 和 Logger

Sugared Logger 性能与易用性并重，支持结构化和 printf 风格的日志记录。

Logger 非常强调性能，不提供 printf 风格的 api （减少了 interface{} 与 反射的性能损耗）


```

package main

import "go.uber.org/zap"

func main() {
	// printf 风格，易用性
	sugar := zap.NewExample().Sugar()
	sugar.Infof("this is %s test log ! line:%d", "sugar", 0)

	// 强调性能
	logger := zap.NewExample()
	logger.Info("this is not sugar test log !", zap.String("status", "not sugar"), zap.Int("line", 1))
}

// 输出结构化日志 
{"level":"info","msg":"this is sugar test log ! line:0"}
{"level":"info","msg":"this is not sugar test log !","status":"not sugar","line":1}

```


zap 有三种默认配置创建出一个 logger，分别为 example，development，production



```

package main

import "go.uber.org/zap"

func main() {

	logger := zap.NewExample()
	logger.Info("example log")

	logger, _ = zap.NewDevelopment()
	logger.Info("development log")

	logger, _ = zap.NewProduction()
	logger.Info("production log")
}

// 不同的环境下 输出格式略有差异
{"level":"info","msg":"example log"}
2023-05-25T23:14:14.126+0800    INFO    cmd/main.go:11  development log
{"level":"info","ts":1685027654.1270118,"caller":"cmd/main.go:14","msg":"production log"}

```


也可以自定义 logger，如以下例子


```

package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 设置 encoding 的日志格式
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewCore(encoder, writeSync("./log.log"), zapcore.InfoLevel)
	logger := zap.New(core)

	logger.Info("this is test")
	logger.Error("this is test ")
}

// // 日志写入
// func writeSync(path string) zapcore.WriteSyncer {
// 	file, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)

// 	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(file))
// }

func writeSync(path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    64,
		MaxBackups: 10,
		MaxAge:     1,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}



// 输出到log.log 文件中
{"level":"info","ts":1685028437.8890965,"msg":"this is test"}
{"level":"error","ts":1685028437.8891673,"msg":"this is test "}

```


从 New(core zapcore.Core, options ...Option) *Logger 出发，需要构造 zapcore.Core

1. 通过 NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler) Core 方法，又需要传入三个参数
2. Encoder : 负责设置 encoding 的日志格式, 可以设置 json 或者 text结构，也可以自定义json中 key 值，时间格式...
3. WriteSyncer: 负责日志写入的位置，上述例子往 file 与 console 同时写入，这里也可以写入网络。
4. LevelEnabler: 设置日志记录级别
5. options 是一个实现了apply(*Logger) 方法的接口，可以扩展很多功能
6. 覆盖 core 的配置
7. 增加 hook
8. 增加键值信息
9. error 日志单独输出



