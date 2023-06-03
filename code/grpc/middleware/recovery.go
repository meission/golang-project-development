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
