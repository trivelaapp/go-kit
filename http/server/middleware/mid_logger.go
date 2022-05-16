package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	// HTTPResponseLatencyKey is the amount of time needed to produce a response to a request.
	// Time measured in milliseconds.
	HTTPResponseLatencyKey = "http.response_latency_in_milli"
)

// LogProvider defines how a logger should behave.
type LogProvider interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warning(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error)
	Critical(ctx context.Context, err error)
}

// Logger creates a new Logger middleware that uses Trivela's GoKit default logger.
func Logger(logger LogProvider) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		start := time.Now()

		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		ctx.Next()

		now := time.Now()
		latency := now.Sub(start).Milliseconds()
		statusCode := ctx.Writer.Status()

		lctx := context.WithValue(ctx, string(semconv.NetPeerIPKey), ctx.ClientIP())
		lctx = context.WithValue(lctx, string(semconv.HTTPMethodKey), method)
		lctx = context.WithValue(lctx, string(semconv.HTTPRouteKey), path)
		lctx = context.WithValue(lctx, string(semconv.HTTPStatusCodeKey), statusCode)
		lctx = context.WithValue(lctx, string(semconv.HTTPResponseContentLengthKey), ctx.Writer.Size())
		lctx = context.WithValue(lctx, string(HTTPResponseLatencyKey), latency)

		msg := fmt.Sprintf("[GIN] %s %s", method, path)

		switch {
		case statusCode >= 500:
			logger.Error(lctx, errors.New(msg))
			break
		case statusCode >= 400:
			logger.Warning(lctx, msg)
			break
		default:
			logger.Info(lctx, msg)
		}
	}
}
