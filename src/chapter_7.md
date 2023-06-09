# http标准库


### 最简单的 HTTP Server
 net/http 包为我们提供了一些基础实现。


#### http 服务端
```golang
package main

import (
	"fmt"
	"net/http"
)

// 请求
func hello(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "hello\n")
	if err != nil {
		return
	}
}

func main() {
	// 定义请求路径及响应函数
	http.HandleFunc("/hello", hello)
	// 启动一个服务并监听 8080 端口
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

```

启动一个 server 很简单，但是这个 server 能做事情还比较少，我们的server要提供各种服务，诸如注册、登录、查询信息，上传图片等功能，那么一个简单的 main.go 还远远不够，那么如何扩展并提供更好的实现呢？



#### http client

``` golang
package main

import (
   "net/http"
   "fmt"
)

func main()  {
   // 使用Get方法获取服务器响应包数据
   resp, err := http.Get("http://127.0.0.1:8080/hello")
   if err != nil {
      fmt.Println("Get err:", err)
      return
   }
   defer resp.Body.Close() //  不close 会有内存泄漏

   // 获取服务器端读到的数据
   fmt.Println("Status = ", resp.Status)           // 状态
   fmt.Println("StatusCode = ", resp.StatusCode)   // 状态码
   fmt.Println("Header = ", resp.Header)           // 响应头部
   fmt.Println("Body = ", resp.Body)               // 响应包体

   buf := make([]byte, 4096)         // 定义切片缓冲区，存读到的内容
   var result string
   // 获取服务器发送的数据包内容
   for {
      n, err := resp.Body.Read(buf)  // 读body中的内容。
      if n == 0 {
         fmt.Println("Body.Read err:", err)
         break
      }
      result += string(buf[:n])     // 累加读到的数据内容
   }
   // 从body中读到的所有内容
   fmt.Println("result = ", result)
}
```


<!-- {{#playground ../code/web/http.go editable no_run should_panic}} -->

如上，整个server 和client 简单实现了


### HTTP中间件机制
Golang 内置的 net/http 天生就支持 HTTP 中间件机制，我们可以改写出扩展性很好的 Web 应用。
每一个中间件都是一层洋葱皮，其中每一个中间件都可以改变请求和响应，我们可以很自然的把不同的逻辑放到不同的洋葱皮里，更代码更符合单一职责原则：


#### 简单的中间件实现方式
```golang
package main

import (
	"fmt"
	"net/http"
)

func middleFoo(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleFoo")
		next(w, r)
	}
}

func middleBar(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleBar")
		next(w, r)
	}
}

// 请求
func hello(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "hello\n")
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/hello", middleFoo(middleBar(hello)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}
```

type HandlerFunc func(ResponseWriter, *Request) 

middleFoo(middleBar(hello)) 这种包洋葱的方式实现了中间件机制,如果维护起来不容易，如果后续很多个中间件就更难受了。 
这种写法后续在标准库http实现工程化较为难，以下是基于此做改进

#### 另一种中间件写法


```golang
type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(m ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

```

完整写法


```golang

package main

import (
	"fmt"
	"net/http"
)

func middleFoo(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleFoo")
		next(w, r)
	}
}

func middleBar(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleBar")
		next(w, r)
	}
}

// 请求
func hello(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "hello\n")
	if err != nil {
		return
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(m ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

func main() {
	// http.HandleFunc("/hello", middleFoo(middleBar(hello)))
	http.HandleFunc("/hello", Chain(middleFoo, middleBar)(hello))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}


```

我们把手动嵌套改成让代码实现链式嵌套，大大提高了维护性，也提高里复用性。
工程化上方便了很多，Chain(m ...Middleware) 可以放进包装库



#### chain 代替写法 
递归的方式，出自grpc-go

```golang
func chainRecursion(middle ...Middleware) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return middle[0](chainHandler(middle, 0, handler))
	}
}

func chainHandler(middle []Middleware, curr int, finalHandler http.HandlerFunc) http.HandlerFunc {
	if curr == len(middle)-1 {
		return finalHandler
	}
	return middle[curr+1](chainHandler(middle, curr+1, finalHandler))
}
```

#### chain 代替写法 二 
出自 grpc-ecosystem，精简之后就是 上述chain写法

```golang
func ChainCirculate(middle ...Middleware) Middleware {
	n := len(middle)

	if n == 0 {
		return func(handler http.HandlerFunc) http.HandlerFunc {
			return handler
		}
	}
	if n == 1 {
		return middle[0]
	}

	return func(handler http.HandlerFunc) http.HandlerFunc {
		currHandler := handler
		for i := n - 1; i > 0; i-- {
			innerHandler, i := currHandler, i
			currHandler = middle[i](innerHandler)
		}
		return middle[0](currHandler)
	}
}

```












