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
