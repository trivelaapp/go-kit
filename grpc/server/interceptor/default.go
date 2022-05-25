package interceptor

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/trivelaapp/go-kit/grpc/server/interceptor/error_handler"
	"github.com/trivelaapp/go-kit/grpc/server/interceptor/logging"
	"github.com/trivelaapp/go-kit/grpc/server/interceptor/meter"
)

// DefaultInput encapsulates inputs to a Default call.
type DefaultInput struct {
	ApplicationName string
	Logger          logging.Logger
}

// Default returns a list of interceptors usually used by most applications. It includes:
// - Disaster Recovery (from panics)
// - Traces
// - Metrics
// - Request Logging
// - Error Handling
func Default(in DefaultInput) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(),
			meter.UnaryServerInterceptor(in.ApplicationName),
			logging.UnaryServerInterceptor(in.Logger),
			error_handler.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			otelgrpc.StreamServerInterceptor(),
			meter.StreamServerInterceptor(in.ApplicationName),
			logging.StreamServerInterceptor(in.Logger),
			error_handler.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),
	}
}
