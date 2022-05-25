package error_handler

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/trivelaapp/go-kit/errors"
)

// UnaryServerInterceptor returns a new unary interceptor suitable for request error handling.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		resp, err = handler(ctx, req)
		if err != nil {
			err = status.Error(kindToGRPCStatusCode(errors.Kind(err)), err.Error())
		}

		return
	}
}

// StreamServerInterceptor returns a new stream interceptor suitable for request error handling.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {

		err = handler(srv, ss)
		if err != nil {
			err = status.Error(kindToGRPCStatusCode(errors.Kind(err)), err.Error())
		}

		return
	}
}
