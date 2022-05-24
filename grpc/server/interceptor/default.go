package interceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Default returns a list of middlewares usually used by most applications. It includes:
// - Disaster Recovery (from panics)
// - Tracing
// - Request Logging
// - Error Handling
func Default() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	}
}
