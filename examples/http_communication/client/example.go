package main

import (
	"context"
	"time"

	"github.com/trivelaapp/go-kit/errors"
	"github.com/trivelaapp/go-kit/http/client"
	"github.com/trivelaapp/go-kit/log"
	"github.com/trivelaapp/go-kit/log/format"
	"github.com/trivelaapp/go-kit/trace"
)

const applicationName = "trivela-http-client"

func main() {
	ctx := context.Background()

	logger := log.NewLogger(log.LoggerParams{
		Level: "INFO",
		Attributes: format.LogAttributeSet{
			format.LogAttribute("foo"): true,
		},
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

	cli := client.New(time.Second)

	ctx, span := tracer.Start(ctx, "test-client")
	defer span.End()

	ctx = context.WithValue(ctx, "foo", "bar")

	res, err := cli.Get(ctx, client.HTTPRequest{
		URL: "http://localhost:3000/error",
	})
	if err != nil {
		logger.Error(ctx, errors.New("can't call Server").WithRootError(err))
		return
	}

	logger.Info(ctx, string(res.Response))
}
