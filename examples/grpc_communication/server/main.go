package main

import (
	"context"
	"fmt"
	"io"
	"net"

	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"github.com/trivelaapp/go-kit/errors"
	tgrpc "github.com/trivelaapp/go-kit/grpc/server/interceptor"
	"github.com/trivelaapp/go-kit/log"
	"github.com/trivelaapp/go-kit/log/format"
	"github.com/trivelaapp/go-kit/metric"
	"github.com/trivelaapp/go-kit/trace"

	pb "go.buf.build/grpc/go/trivela/go-kit-example/example/v1"
)

const applicationName = "example-grpc-communication-server"
const serverPort = 8080

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, string(semconv.ServiceNameKey), applicationName)
	ctx = context.WithValue(ctx, string(semconv.ServiceVersionKey), "v0.0.0")

	logger := log.NewLogger(log.LoggerParams{
		Level:      "INFO",
		Attributes: format.DefaultGRPCServerAttributeSet,
	})

	trace := trace.MustNewJaegerTracerProvider(trace.JaegerTracerProviderParams{
		ApplicationName:    applicationName,
		ApplicationVersion: "v0.0.0",
		Endpoint:           "http://localhost:14268/api/traces",
		TraceRatio:         1,
	})
	tracer, flush, err := trace.Tracer(ctx)
	if err != nil {
		logger.Fatal(ctx, err)
	}
	defer flush(ctx)

	metric := metric.MustNewPrometheusMeterProvider(metric.PrometheusMeterProviderParams{
		ApplicationName:     applicationName,
		ApplicationVersion:  "v0.0.0",
		MetricsServerPort:   8081,
		HistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	})
	if _, _, err := metric.Meter(ctx); err != nil {
		logger.Fatal(ctx, err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		logger.Fatal(ctx, errors.New("Could not initialize net http listener").WithRootError(err))
	}

	interceptors := tgrpc.Default(tgrpc.DefaultInput{
		ApplicationName: applicationName,
		Logger:          logger,
	})
	srv := grpc.NewServer(interceptors...)
	pb.RegisterExampleServiceServer(srv, helloServer{
		tracer: tracer,
		logger: logger,
	})

	logger.Info(ctx, "Server up and running on port :%d", serverPort)
	if err := srv.Serve(lis); err != nil {
		logger.Fatal(ctx, errors.New("Could not start serving incoming requests at listener").WithRootError(err))
	}
}

type helloServer struct {
	tracer otrace.Tracer
	logger *log.Logger
}

// Hello greets the user.
func (s helloServer) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	_, span := s.tracer.Start(ctx, "server")
	defer span.End()

	// return nil, errors.New("test").WithKind(errors.KindUnauthorized)

	return &pb.HelloResponse{
		Message: fmt.Sprintf("What's up, %s?", req.Name),
	}, nil
}

// StreamHello greets users by an stream.
func (s helloServer) StreamHello(stream pb.ExampleService_StreamHelloServer) error {
	ctx := stream.Context()
	// return errors.New("test").WithKind(errors.KindUnauthorized)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			err = errors.New("client request failed").WithRootError(err)
			s.logger.Error(ctx, err)
			return err
		}

		stream.Send(&pb.StreamHelloResponse{
			Data: &pb.HelloResponse{
				Message: fmt.Sprintf("What's up, %s?", req.Data.Name),
			},
		})
	}
}
