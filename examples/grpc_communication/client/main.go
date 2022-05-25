package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/trivelaapp/go-kit/errors"
	"github.com/trivelaapp/go-kit/log"
	"github.com/trivelaapp/go-kit/metric"
	"github.com/trivelaapp/go-kit/trace"

	pb "go.buf.build/grpc/go/trivela/go-kit-example/example/v1"
)

const applicationName = "example-grpc-communication-client"

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, string(semconv.ServiceNameKey), applicationName)
	ctx = context.WithValue(ctx, string(semconv.ServiceVersionKey), "v0.0.0")

	logger := log.NewLogger(log.LoggerParams{
		Level: "INFO",
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

	ctx, span := tracer.Start(ctx, "client")
	defer span.End()

	conn, err := grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		logger.Fatal(ctx, errors.New("did not connect to server").WithRootError(err))
	}
	defer conn.Close()

	client := pb.NewExampleServiceClient(conn)

	res, err := client.Hello(ctx, &pb.HelloRequest{
		Name: "Unary Client",
	})
	if err != nil {
		logger.Fatal(ctx, errors.New("could not request data to server").WithRootError(err))
	}

	logger.Info(ctx, res.Message)

	time.Sleep(time.Second * 2)

	stream, err := client.StreamHello(ctx)
	if err != nil {
		logger.Fatal(ctx, errors.New("could not open server stream").WithRootError(err))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				logger.Error(ctx, errors.New("server error received").WithRootError(err))
				return
			}

			logger.Info(ctx, res.Data.Message)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 1; i <= 5; i++ {
			err := stream.Send(&pb.StreamHelloRequest{
				Data: &pb.HelloRequest{
					Name: fmt.Sprintf("Stream Client %d", i),
				},
			})
			if err != nil {
				logger.Error(ctx, errors.New("could not send data to server").WithRootError(err))
				return
			}

			time.Sleep(time.Second * 1)
		}

		if err := stream.CloseSend(); err != nil {
			logger.Error(ctx, errors.New("could not close send stream").WithRootError(err))
		}
	}()

	wg.Wait()
}
