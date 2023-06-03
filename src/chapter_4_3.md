# 配置库viper


介绍一个轻便好用的Golang配置库viper

### 介绍Viper

Viper是一个完整的Go应用程序的配置解决方案，包含12因素应用程序。它被设计为在应用程序中工作，并能处理所有类型的配置需求和格式。
Viper主要是用于处理各种格式的配置文件，简化程序配置的读取问题。


设置默认值
从JSON，TOML，YAML，HCL和Java属性配置文件中读取
实时观看和重新读取配置文件（可选）
从环境变量中读取
从远程配置系统（etcd或Consul）读取，并观察变化
从命令行标志读取（flag）
从缓冲区读取
设置显式值

当构建一个现代的应用程序时，你不想担心配置文件的格式；你想专注于构建令人敬畏的软件。Viper就是来帮助你的。

安装
```
go get github.com/spf13/viper

```

### 读取toml 配置文件

#### 配置文件 config.toml

```
# 全局信息
title = "TOML示例"

# 应用信息
[app]
    author = "史布斯"
    organization = "Mafool"
    mark = "第一行\n第二行."            # 换行
    release = 2020-05-27T07:32:00Z   # 时间

# 数据库配置
[mysql]
    server = "192.168.1.1"
    ports = [ 8001, 8001, 8002 ]     # 数组
    connection_max = 5000
    enabled = true

# Redis主从                           # 字典对象
[redis]
    [redis.master]
        host = "10.0.0.1"
        port = 6379
    [redis.slave]
        host = "10.0.0.1"
        port = 6380

# 二维数组
[releases]
release = ["dev", "test", "stage", "prod"]
tags = [["dev", "stage", "prod"],[2.2, 2.1]]


# 公司信息                             #对象嵌套
[company]
    name = "xx科技"
[company.detail]
    type = "game"
    addr = "北京朝阳"
    icp = "030173"

```


#### main.go 代码
```golang

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Title    string
	App      app
	DB       mysql `toml:"mysql"`
	Redis    map[string]redis
	Releases releases
	Company  Company
}

type app struct {
	Author  string
	Org     string `toml:"organization"`
	Mark    string
	Release time.Time
}

type mysql struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type redis struct {
	Host string
	Port int
}

type releases struct {
	Release []string
	Tags    [][]interface{}
}

type Company struct {
	Name   string
	Detail detail
}

type detail struct {
	Type string
	Addr string
	ICP  string
}

func main() {
	var config Config
	viper.SetConfigFile("./config.toml")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}
	// fmt.Println("all settings: ", viper.AllSettings())
	viper.Unmarshal(&config)
	fmt.Printf("config %+v\n", config)
}

```



### 读取yaml 文件



#### config.yaml 

```yaml
database:
  server: 8090
  ports:
    - 8001
    - 8001
    - 8002

servers:
  alpha:
    ip: 10.0.0.1

```

##### main.go

```golang
package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Database struct {
	Server string
	Ports  []int
}

type Servers struct {
	Alpha detail
}

type detail struct {
	IP string
}

type Config struct {
	Database Database
	Servers  Servers
}

func main() {
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}

	var config Config
	viper.Unmarshal(&config)
	fmt.Printf("config %+v\n", config)
}


````


###  监听配置变化
作用：热更新时，可能会更改配置文件，方便重新读取。

```
//监听配置变化
viper.WatchConfig()
//配置改变时的回调
viper.OnConfigChange(func(in fsnotify.Event) {
   switch in.Op {
   case fsnotify.Create:
      fmt.Println("监听到Create")
   case fsnotify.Remove:
      fmt.Println("监听到Remove")
   case fsnotify.Rename:
      fmt.Println("监听到Rename")
   case fsnotify.Write:
      fmt.Println("监听到Write")
   case fsnotify.Chmod:
      fmt.Println("监听到Chmod")
   }
})

```
 
viper 还有一些功能值得探索
