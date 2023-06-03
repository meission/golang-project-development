# grpc

<https://github.com/grpc/grpc-go>


简介
gRPC  是一个高性能、开源和通用的 RPC 框架，面向移动和 HTTP/2 设计。目前提供 C、Java 和 Go 语言版本，分别是：grpc, grpc-java, grpc-go. 其中 C 版本支持 C, C++, Node.js, Python, Ruby, Objective-C, PHP 和 C# 支持.

gRPC 基于 HTTP/2 标准设计，带来诸如双向流、流控、头部压缩、单 TCP 连接上的多复用请求等特。这些特性使得其在移动设备上表现更好，更省电和节省空间占用。


安装protoc

这个工具也称为proto编译器，可以用来生成各种开发语言使用proto协议的代码。
https://github.com/protocolbuffers/protobuf/releases



安装protoc的Golang gRPC插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

```


项目目录结构如下

```bash
.
├── api
│   ├── hello_grpc.pb.go
│   ├── hello.pb.go
│   ├── hello.proto
│   └── proto.sh
├── client
│   └── main.go
└── main.go

2 directories, 6 files

```

hello.proto 文件内容

```proto

syntax = "proto3";

option go_package="/api";

package api;

service Hello {
  rpc Say (SayRequest) returns (SayResponse);
}

message SayResponse {
  string Message = 1;
}

message SayRequest {
  string Name = 1;
}

```


proto.sh 文件内容

```bash

#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative hello.proto


```
在上述api 目录下执行 ./proto.sh 就会生成 hello_grpc.pb.go hello.pb.go 文件。如果没有执行权限，执行 chmod +x ./proto.sh 增加执行权限


### 简单实例

server 端代码

```golang

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "code/grpc/api"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) Say(ctx context.Context, in *pb.SayRequest) (*pb.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.SayResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &server{})
	defer func() {
		s.Stop()
		lis.Close()
	}()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```

client 端代码


```golang

package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "code/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Say(ctx, &pb.SayRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}


```

###  中间件


写完中间件 目录结构如下，middleware 下为中间件写法实例

```bash
.
├── api
│   ├── hello_grpc.pb.go
│   ├── hello.pb.go
│   ├── hello.proto
│   └── proto.sh
├── client
│   └── main.go
├── main.go
└── middleware
    ├── ratelimit.go
    ├── recovery.go
    └── timeout.go

3 directories, 9 files

```

#### 写法一


Recovery 写法 较为简单,普通模式

```golang
package middleware

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/grpc"
)

const size = 64 << 10

// recovery returns a new unary server interceptor for panic recovery.
func Recovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("grpc server panic: %v\n%v\n%s\n", req, r, buf)
				fmt.Println(pl)
			}
		}()
		return handler(ctx, req)
	}
}


```

#### 写法二 


Limiter 写法，定义了Limiter 接口形式，可以实现各种各样的算法

```golang

package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Limit(ctx context.Context) error
}

// Ratelimit returns a new unary server interceptors that performs request rate limiting.
func Ratelimit(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := limiter.Limit(ctx); err != nil {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later. %s", info.FullMethod, err)
		}
		return handler(ctx, req)
	}
}

type AlwaysPassLimiter struct{}

func (*AlwaysPassLimiter) Limit(_ context.Context) error {
	fmt.Println("AlwaysPassLimiter")
	return nil
}

```


#### 写法三

Timeout 写法实现了 option 模式，适合固定的options struct ，其中包含里多个可自定义的参数

```golang

package middleware

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type options struct {
	callTimeout time.Duration
}

type CallOption struct {
	grpc.EmptyCallOption // make sure we implement private after() and before() fields so we don't panic.
	applyFunc            func(opt *options)
}

func WithTimeout(timeout time.Duration) CallOption {
	return CallOption{applyFunc: func(o *options) {
		o.callTimeout = timeout
	}}
}

// Timeout returns a new unary server interceptor for Timeout .
func Timeout(callOptions ...CallOption) grpc.UnaryServerInterceptor {
	fmt.Println("Timeout")
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		if len(callOptions) == 0 {
			return handler(ctx, req)
		}
		optCopy := &options{}
		for _, f := range callOptions {
			f.applyFunc(optCopy)
		}
		timedCtx, cancel := context.WithTimeout(ctx, optCopy.callTimeout)
		defer cancel()
		return handler(timedCtx, req)

	}
}

```


此时server端代码做一些调整引入中间件

```golang

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "code/grpc/api"
	"code/grpc/middleware"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) Say(ctx context.Context, in *pb.SayRequest) (*pb.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.SayResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := grpc.ChainUnaryInterceptor(
		middleware.Recovery(),
		middleware.Ratelimit(&middleware.AlwaysPassLimiter{}),
		middleware.Timeout(middleware.WithTimeout(100*time.Millisecond)),
	)

	s := grpc.NewServer(interceptor)
	pb.RegisterHelloServer(s, &server{})
	defer func() {
		s.Stop()
		lis.Close()
	}()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```

其中重点代码改动如下

```golang

    interceptor := grpc.ChainUnaryInterceptor(
        middleware.Recovery(),
        middleware.Ratelimit(&middleware.AlwaysPassLimiter{}),
        middleware.Timeout(middleware.WithTimeout(100*time.Millisecond)),
    )

    s := grpc.NewServer(interceptor)

```





