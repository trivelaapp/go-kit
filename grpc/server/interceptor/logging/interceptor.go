package logging

import (
	"context"
	"fmt"
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/trivelaapp/go-kit/errors"
)

const (
	// GRPCResponseLatencyKey is the amount of time needed to produce a response to a request.
	GRPCResponseLatencyKey = "rpc.grpc.response_latency"
)

// UnaryServerInterceptor returns a new unary interceptor suitable for request logging.
func UnaryServerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		now := time.Now()
		latency := now.Sub(start).String()

		lctx := context.WithValue(ctx, string(semconv.RPCMethodKey), info.FullMethod)
		lctx = context.WithValue(lctx, string(semconv.RPCGRPCStatusCodeKey), status.Code(err).String())
		lctx = context.WithValue(lctx, GRPCResponseLatencyKey, latency)

		peer, ok := peer.FromContext(ctx)
		if ok {
			lctx = context.WithValue(lctx, string(semconv.NetPeerIPKey), peer.Addr.String())
		}

		msg := fmt.Sprintf("[gRPC] %s", info.FullMethod)
		if err != nil {
			logger.Error(lctx, errors.New(msg).WithRootError(err))
		} else {
			logger.Info(lctx, msg)
		}

		return resp, err
	}
}

// StreamServerInterceptor returns a new stream interceptor suitable for request logging.
func StreamServerInterceptor(logger Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {

		ctx := ss.Context()
		start := time.Now()

		err = handler(srv, ss)

		now := time.Now()
		latency := now.Sub(start).String()

		lctx := context.WithValue(ctx, string(semconv.RPCMethodKey), info.FullMethod)
		lctx = context.WithValue(lctx, string(semconv.RPCGRPCStatusCodeKey), status.Code(err).String())
		lctx = context.WithValue(lctx, GRPCResponseLatencyKey, latency)

		peer, ok := peer.FromContext(ctx)
		if ok {
			lctx = context.WithValue(lctx, string(semconv.NetPeerIPKey), peer.Addr.String())
		}

		msg := fmt.Sprintf("[gRPC] %s", info.FullMethod)
		if err != nil {
			logger.Error(lctx, errors.New(msg).WithRootError(err))
		} else {
			logger.Info(lctx, msg)
		}

		return
	}
}
