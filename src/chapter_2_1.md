# log4go

项目路径 <https://github.com/meission/log4go>



这里介绍一个轻量级的log模块——log4go. 源于google的一项log工程，官方已经停止维护更新，这里对他进行了稍微调整，使用起来也特别简单，就像自身的log模块一样。

简洁，核心代码只有10K左右，但是功能并不简单。

在自己的go代码中，只需要配置简单的log的路径，以及需要打印的日志级相关信息，即可使用日志工具。该日志工具支持将日志文件按时间、文件大小、日志级别进行文件切分。


```golang

package log

import (
	"math"
	"path"

	log "github.com/meission/log4go"
)

var (
	logger log.Logger
)

type Config struct {
	Dir string
}

func Init(c *Config) {
	if c == nil || c.Dir == "" {
		c = &Config{
			Dir: "./",
		}
	}
	logger = log.Logger{}
	log.LogBufferLength = 10240
	// new info writer
	iw := log.NewFileLogWriter(path.Join(c.Dir, "info.log"), false)
	iw.SetRotateDaily(true)
	iw.SetRotateSize(math.MaxInt32)
	iw.SetRotate(true)
	iw.SetFormat("[%D %T] [%L] [%S] %M")
	logger["info"] = &log.Filter{
		Level:     log.INFO,
		LogWriter: iw,
	}
	// new warning writer
	ww := log.NewFileLogWriter(path.Join(c.Dir, "warning.log"), false)
	ww.SetRotateDaily(true)
	ww.SetRotateSize(math.MaxInt32)
	ww.SetRotate(true)
	ww.SetFormat("[%D %T] [%L] [%S] %M")
	logger["warning"] = &log.Filter{
		Level:     log.WARNING,
		LogWriter: ww,
	}
	// new error writer
	ew := log.NewFileLogWriter(path.Join(c.Dir, "error.log"), false)
	ew.SetRotateDaily(true)
	ew.SetRotateSize(math.MaxInt32)
	ew.SetRotate(true)
	ew.SetFormat("[%D %T] [%L] [%S] %M")
	logger["error"] = &log.Filter{
		Level:     log.ERROR,
		LogWriter: ew,
	}
}

// Close close resource.
func Close() {
	if logger != nil {
		logger.Close()
	}
}

// Info write info log .
func Info(format string, args ...interface{}) {
	if logger != nil {
		logger.Info(format, args...)
	}
}

// Warn write warn log .
func Warn(format string, args ...interface{}) {
	if logger != nil {
		logger.Warn(format, args...)
	}
}

// Error write error log .
func Error(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(format, args...)
	}
}



```

#### 使用样例 
```

package main

import (
	"code/log"
	"time"
)

func main() {

	log.Init(nil)
	defer log.Close()
	log.Info("%s", "this is test info log ")
	log.Error("%s", "this is test error log ")
	time.Sleep(time.Second * 2)
}

```


#### 输出结果

##### info 
```
[2023/05/25 17:52:48 CST] [INFO] [main.main:9] this is test info log 
[2023/05/25 17:52:48 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:52:54 CST] [INFO] [main.main:9] this is test info log 
[2023/05/25 17:52:54 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:54:56 CST] [INFO] [main.main:12] this is test info log 
[2023/05/25 17:54:56 CST] [EROR] [main.main:13] this is test error log 

```

##### error
```
[2023/05/25 17:52:35 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:52:48 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:52:54 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:54:56 CST] [EROR] [main.main:13] this is test error log 

```

##### warning 

```
[2023/05/25 17:51:05 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:52:48 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:52:54 CST] [EROR] [main.main:10] this is test error log 
[2023/05/25 17:54:56 CST] [EROR] [main.main:13] this is test error log 

```



