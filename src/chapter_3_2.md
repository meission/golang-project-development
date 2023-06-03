# Cobra
 

Cobra既是一个用来创建强大的现代命令行应用的库，又是一个用来生成应用和命令文件的脚手架。


cobra用于许多著名的开源项目，如Kubernetes，Hugo和Github Cli 很多流行的Go项目都使用Cobra，例如Kubernetes, Hugo, rkt, etcd, Moby (former Docker), Docker (distribution), OpenShift, Delve, GopherJS, CockroachDB, Bleve, ProjectAtomic (enterprise), Giant Swarm’s gsctl, Nanobox/Nanopack, rclone, nehm, Pouch, Istio, Prototool, mattermost-server, Gardener, Linkerd等。



### Cobra提供的功能

1. 简易的子命令行模式，如 app server， app fetch等等
2. 完全兼容posix命令行模式 （包括短和长版本）
3. 嵌套子命令subcommand
4. 支持全局，局部，串联flags
5. 通过cobra init appname和cobra add cmdname轻松生成应用程序和命令
6. 智能建议（如果命令输入错误，将提供智能建议 app srver...是app server吗？）
7. 自动生成commands和flags的帮助信息
8. 自动识别-h、--help等标志。
9. 自动生成应用程序在bash下命令自动完成功能
10. 自动生成应用程序的man手册
11. 命令别名，在不破坏它们的情况下进行更改
12. 灵活地自定义help和usage信息
13. 可选的紧密集成的viper apps 用于12个因素的应用程序

 

### 安装

```
$ go get -u github.com/spf13/cobra@latest

import "github.com/spf13/cobra"

``` 

安装 cobra-cli

```
$ go install github.com/spf13/cobra-cli@latest
```

初始化项目

```
cd $HOME/code/myapp
cobra-cli init
go run main.go
```

初始化后项目目录结构

```
.
├── cmd
│   └── root.go
├── LICENSE
└── main.go

1 directory, 3 files

```

增加命令


```
$ cobra-cli add serve
$ cobra-cli add config
$ cobra-cli add create -p 'configCmd' // 在config命令下 下增加create子命令
```

增加完命令之后 目录结构

```
.
├── cmd
│   ├── config.go
│   ├── create.go
│   ├── root.go
│   └── serve.go
├── LICENSE
└── main.go

1 directory, 6 files


```

```
$ go run main.go
```
运行命令之后输出结果：

```
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  myapp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      A brief description of your command
  help        Help about any command
  serve       A brief description of your command

Flags:
  -h, --help     help for myapp
  -t, --toggle   Help message for toggle

Use "myapp [command] --help" for more information about a command.

```

运行
```
$ go run main.go config -h
```
之后结果 

```bash
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  myapp config [flags]
  myapp config [command]

Available Commands:
  create      A brief description of your command

Flags:
  -h, --help   help for config

Use "myapp config [command] --help" for more information about a command.

```

此时使用cobra 模板已经初始化完毕，接下来就是完善自己的功能了。。。。






