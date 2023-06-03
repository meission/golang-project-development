# toml


安装

```
$ go get github.com/BurntSushi/toml@latest

```

配置文件 config.toml
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


main.go 代码

```golang
package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
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
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
	}

	fmt.Printf("config: %+v\n", config)
}

```




