package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trivelaapp/go-kit/errors"
	"github.com/trivelaapp/go-kit/http/server"
	"github.com/trivelaapp/go-kit/log"
	"github.com/trivelaapp/go-kit/metric"
	"github.com/trivelaapp/go-kit/trace"
)

const applicationName = "trivela-http-server"

func main() {
	ctx := context.Background()

	logger := log.NewLogger(log.LoggerParams{
		Level: "INFO",
	})

	trace := trace.MustNewJaegerTracerProvider(trace.JaegerTracerProviderParams{
		ApplicationName:    applicationName,
		ApplicationVersion: "v0.0.0",
		Endpoint:           "http://localhost:14268/api/traces",
		TraceRatio:         1,
	})
	_, flush, err := trace.Tracer(ctx)
	if err != nil {
		logger.Panic(ctx, err)
	}
	defer flush(ctx)

	metric := metric.MustNewPrometheusMeterProvider(metric.PrometheusMeterProviderParams{
		ApplicationName:     applicationName,
		ApplicationVersion:  "v0.0.0",
		MetricsServerPort:   8081,
		HistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	})
	if _, _, err := metric.Meter(ctx); err != nil {
		logger.Panic(ctx, err)
	}

	srv := server.MustNew(server.ServerParams{
		ApplicationName: applicationName,
		Port:            8080,
		Logger:          logger,
	})

	router := srv.RouterDefault()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, server.NewMessageResponse("pong"))
	})

	router.GET("/error", func(ctx *gin.Context) {
		err := errors.New("something went wrong").
			WithKind(errors.KindNotFound).
			WithCode("TEST")

		ctx.Error(err)
		ctx.Abort()
	})

	router.Run()
}
