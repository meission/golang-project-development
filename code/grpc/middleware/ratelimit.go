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
