# urfave/cli


官方代码 <https://github.com/urfave/cli>

很多用Go写的命令行程序都用了urfave/cli这个库。urfave/cli是一个命令行的框架。

urfave/cli把这个过程做了一下封装，抽象出flag/command/subcommand这些模块，用户只需要提供一些模块的配置，参数的解析和关联在库内部完成，帮助信息也可以自动生成。

总体来说，urfave/cli这个库还是很好用的，完成了很多routine的工作，程序员只需要专注于具体业务逻辑的实现。


### 安装使用


#### 使用v2版本(推荐使用)

```
$ go get github.com/urfave/cli/v2

...
import (
  "github.com/urfave/cli/v2" // imports as package "cli"
)
...
```

#### 使用v1版本

```
$ go get github.com/urfave/cli
 
...
import (
  "github.com/urfave/cli"
)
...
```


#### 简单用法


```
package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gocli",
		Usage: "gocli is test cli ",
		Commands: []*cli.Command{ //命令
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "file list test",
				Action: func(ctx *cli.Context) error {
					fmt.Println("list: \n main.go \n cmd.go", ctx.Args().First())
					return nil
				},
			},
			{
				Name:    "proto",
				Aliases: []string{"p"},
				Usage:   "protobuf file list",
				Action: func(ctx *cli.Context) error {
					fmt.Println("list: \n api.proto \n cmd.proto", ctx.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

```
$ go run main.go
```

运行结果展示
``` 
NAME:
   gocli - gocli is test cli 

USAGE:
   gocli [global options] command [command options] [arguments...]

COMMANDS:
   list, l   file list test
   proto, p  protobuf file list
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

```

```
$ go run  main.go l
```

运行结果展示

```
list: 
 main.go 
 cmd.go 

```


```
$ go run  main.go p
```

```
list: 
 api.proto 
 cmd.proto 

```


#### 带子命令用法


```
package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gocli",
		Usage: "gocli is test cli ",
		Commands: []*cli.Command{ //命令
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "file list test",
				Action: func(ctx *cli.Context) error {
					fmt.Println("list: \n main.go \n cmd.go", ctx.Args().First())
					return nil
				},
			},
			{
				Name:    "proto",
				Aliases: []string{"p"},
				Usage:   "protobuf file list",
				// Action: func(ctx *cli.Context) error {
				// 	fmt.Println("list: \n api.proto \n cmd.proto", ctx.Args().First())
				// 	return nil
				// },
				Subcommands: []*cli.Command{ //子命令
					{
						Name:  "grpc",
						Usage: "gen grpc template",
						Action: func(ctx *cli.Context) error {
							fmt.Println("gen grpc template: ", ctx.Args().First())
							return nil
						},
					},
					{
						Name:  "http",
						Usage: "gen http template",
						Action: func(ctx *cli.Context) error {
							fmt.Println("gen http template: ", ctx.Args().First())
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

```


```
$ go run  main.go p
$ go run  main.go p h
```


运行结果

```bash
NAME:
   gocli proto - protobuf file list

USAGE:
   gocli proto command [command options] [arguments...]

COMMANDS:
   grpc     gen grpc template
   http     gen http template
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

```
$ go run  main.go p grpc
```

运行结果

```
gen grpc template:
```

