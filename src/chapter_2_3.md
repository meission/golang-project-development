# logrus

logrus完全兼容标准的log库，还支持文本、JSON 两种日志输出格式。很多知名的开源项目都使用了这个库
logrus : Logrus 是 Go (golang) 的结构化日志记录器，与标准库 log 完全 API 兼容。（19.1 k stars）

logrus具有以下特性：

1. 完全兼容golang标准库日志模块：logrus拥有七种日志级别：debug、info、warn、error、fatal和panic，trace，这是golang标准库日志模块的API的超集。
2. 可扩展的Hook机制：允许使用者通过hook的方式将日志分发到任意地方，如本地文件系统、标准输出、logstash、elasticsearch或者mq等，或者通过hook定义日志内容和格式等。
3. 可选的日志输出格式：logrus内置了两种日志格式，JSONFormatter和TextFormatter，如果这两个格式不满足需求，可以自己动手实现接口Formatter，来定义自己的日志格式。
4. Field机制：logrus鼓励通过Field机制进行精细化的、结构化的日志记录，而不是通过冗长的消息来记录日志。
5. logrus是一个可插拔的、结构化的日志框架。
6. 线程安全


logrus支持更多的日志级别：

0. Panic：记录日志，然后panic。
1. Fatal：致命错误，出现错误时程序无法正常运转。输出日志后，程序退出；
2. Error：错误日志；
3. Warn：警告信息；
4. Info：信息级别日志；
5. Debug：调试级别；
6. Trace：追踪级别，验证逻辑流程；

线程安全
默认情况下，logrus的api都是线程安全的，其内部通过互斥锁来保护并发写。互斥锁工作于调用hooks或者写日志的时候，如果不需要锁，可以调用logger.SetNoLock()来关闭之。可以关闭logrus互斥锁的情形包括：

没有设置hook，或者所有的hook都是线程安全的实现。
写日志到logger.Out已经是线程安全的了，如logger.Out已经被锁保护，或者写文件时，文件是以O_APPEND方式打开的，并且每次写操作都小于4k。
 

简单使用


```
package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// 在测试环境中设置低等级级别，在生产环境中需要考虑性能，level 设置高一点logrus.InfoLevel
	logrus.SetLevel(logrus.TraceLevel)
	// 调用者文件名与位置
	logrus.SetReportCaller(true)
	// 日志格式设置成json
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// 日志样例
	logrus.Traceln("trace 日志")
	logrus.Debugln("debug 日志")
	logrus.Infoln("Info 日志")
	logrus.Warnln("warn 日志")
	logrus.Errorln("error msg")
	logrus.Fatalf("fatal 日志")
	logrus.Panicln("panic 日志")
	time.Sleep(time.Second * 2)
}


```



日志级别从上向下依次增加，Trace最大，Panic最小。logrus有一个日志级别，高于这个级别的日志不会输出。默认的级别为InfoLevel。所以为了能看到Trace和Debug日志，我们在main函数第一行设置日志级别为TraceLevel。


考虑到易用性，库一般会使用默认值创建一个对象，包最外层的方法一般都是操作这个默认对象，logrus 与 log 都是这样

结构化前：

```
TRAC[0000]/home/meission/gopath/src/golang-project-development/code/cmd/main.go:18 main.main() trace 日志                                     
DEBU[0000]/home/meission/gopath/src/golang-project-development/code/cmd/main.go:19 main.main() debug 日志                                     
INFO[0000]/home/meission/gopath/src/golang-project-development/code/cmd/main.go:20 main.main() Info 日志                                      
WARN[0000]/home/meission/gopath/src/golang-project-development/code/cmd/main.go:21 main.main() warn 日志                                      
ERRO[0000]/home/meission/gopath/src/golang-project-development/code/cmd/main.go:22 main.main() error msg                                    
FATA[0000]/home/meission/gopath/src/golang-project-development/code/cmd/main.go:23 main.main() fatal 日志                                     
exit status 1


```

结构化后：


```
{"file":"/home/meission/gopath/src/golang-project-development/code/cmd/main.go:18","func":"main.main","level":"trace","msg":"trace 日志","time":"2023-05-25T21:42:31+08:00"}
{"file":"/home/meission/gopath/src/golang-project-development/code/cmd/main.go:19","func":"main.main","level":"debug","msg":"debug 日志","time":"2023-05-25T21:42:31+08:00"}
{"file":"/home/meission/gopath/src/golang-project-development/code/cmd/main.go:20","func":"main.main","level":"info","msg":"Info 日志","time":"2023-05-25T21:42:31+08:00"}
{"file":"/home/meission/gopath/src/golang-project-development/code/cmd/main.go:21","func":"main.main","level":"warning","msg":"warn 日志","time":"2023-05-25T21:42:31+08:00"}
{"file":"/home/meission/gopath/src/golang-project-development/code/cmd/main.go:22","func":"main.main","level":"error","msg":"error msg","time":"2023-05-25T21:42:31+08:00"}
{"file":"/home/meission/gopath/src/golang-project-development/code/cmd/main.go:23","func":"main.main","level":"fatal","msg":"fatal 日志","time":"2023-05-25T21:42:31+08:00"}
exit status 1


``` 

添加字段

有时候需要在输出中添加一些字段，可以通过调用logrus.WithField和logrus.WithFields实现。logrus.WithFields接受一个logrus.Fields类型的参数，其底层实际上为map[string]interface{}


```
package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"app":  "logrus-test",
		"zone": "shanghai",
	})

	logger.Info("info msg")
	logger.Error("error msg")
}

// 输出结果带app 和 zone 关键字。 在微服务架构里可用于在日志中心平台上定位应用名及环境
INFO[0000] info msg                                      app=logrus-test zone=shanghai
ERRO[0000] error msg                                     app=logrus-test zone=shanghai

```


输出到文件

```
package main

import (
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "./log.log", //日志文件位置
		MaxSize:    64,          // 单文件最大容量,单位是MB
		MaxBackups: 10,          // 最大保留过期文件个数
		MaxAge:     1,           // 保留过期文件的最大时间间隔,单位是天
		Compress:   true,        // 是否需要压缩滚动日志, 使用的 gzip 压缩
	})
	logger := logrus.WithFields(logrus.Fields{
		"app":  "logrus-test",
		"zone": "shanghai",
	})

	logger.Info("info msg")
	logger.Error("error msg")
}


time="2023-05-25T22:24:22+08:00" level=info msg="info msg" app=logrus-test zone=shanghai
time="2023-05-25T22:24:22+08:00" level=error msg="error msg" app=logrus-test zone=shanghai
time="2023-05-25T22:25:11+08:00" level=info msg="info msg" app=logrus-test zone=shanghai
time="2023-05-25T22:25:11+08:00" level=error msg="error msg" app=logrus-test zone=shanghai


```
Hook接口扩展


logrus 通过实现 Hook接口扩展 hook 机制,可以根据需求将日志分发到任意的存储介质, 比如 es, mq 或者监控报警系统,及时获取异常日志。可以说极大的提高了日志系统的可扩展性。

```
type Hook interface {
  // 定义哪些等级的日志触发 hook 机制
	Levels() []Level
  // hook 触发器的具体执行操作
  // 如果 Fire 执行失败,错误日志会重定向到标准错误流
	Fire(*Entry) error
}
```


自定义hook

```
package main

import (
	"os"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// reportHook ...
type reportHook struct {
}

// Levels 只定义 error 和 panic 等级的日志触发 hook
func (h *reportHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.PanicLevel,
	}
}

// Fire 将日志写入到指定日志文件中
func (h *reportHook) Fire(entry *logrus.Entry) error {
	f, err := os.OpenFile("./err.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(entry.Message)); err != nil {
		return err
	}
	return nil
}

func main() {

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "./log.log", //日志文件位置
		MaxSize:    64,          // 单文件最大容量,单位是MB
		MaxBackups: 10,          // 最大保留过期文件个数
		MaxAge:     1,           // 保留过期文件的最大时间间隔,单位是天
		Compress:   true,        // 是否需要压缩滚动日志, 使用的 gzip 压缩
	})
	logrus.AddHook(&reportHook{})
	logger := logrus.WithFields(logrus.Fields{
		"app":  "logrus-test",
		"zone": "shanghai",
	})

	logger.Info("info msg")
	logger.Error("error msg")
}

```

已经实现的hook 库，
https://github.com/sirupsen/logrus/wiki/Hooks








