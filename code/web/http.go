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

func main() {
	// http.HandleFunc("/hello", middleFoo(middleBar(hello)))
	// http.HandleFunc("/hello", chainRecursion(middleFoo, middleBar)(hello))
	http.HandleFunc("/hello", ChainCirculate(middleFoo, middleBar)(hello))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}
