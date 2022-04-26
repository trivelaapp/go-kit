package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

const (
	requestCounterName           = "http.server.request_counter"
	requestBytesCounterName      = "http.server.request_bytes_counter"
	responseBytesCounterName     = "http.server.response_bytes_counter"
	responseLatencyHistogramName = "http.server.duration"
)

// Meter creates a new Meter middleware that measures application performance.
func Meter(applicationName string) func(ctx *gin.Context) {
	meter := global.MeterProvider().Meter(applicationName)

	requestCounter, _ := meter.SyncInt64().Counter(requestCounterName)
	requestBytesCounter, _ := meter.SyncInt64().Counter(requestBytesCounterName)
	responseBytesCounter, _ := meter.SyncInt64().Counter(responseBytesCounterName)
	responseLatency, _ := meter.SyncInt64().Histogram(responseLatencyHistogramName)

	counters := map[string]syncint64.Counter{
		requestCounterName:       requestCounter,
		requestBytesCounterName:  requestBytesCounter,
		responseBytesCounterName: responseBytesCounter,
	}

	histograms := map[string]syncint64.Histogram{
		responseLatencyHistogramName: responseLatency,
	}

	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		attributes := []attribute.KeyValue{
			attribute.String("path", ctx.FullPath()),
			attribute.String("method", ctx.Request.Method),
			attribute.Int("status_code", ctx.Writer.Status()),
		}

		counters[requestCounterName].Add(ctx, 1, attributes...)
		counters[requestBytesCounterName].Add(ctx, ctx.Request.ContentLength, attributes...)
		counters[responseBytesCounterName].Add(ctx, int64(ctx.Writer.Size()))

		latency := time.Now().Sub(start).Milliseconds()
		histograms[responseLatencyHistogramName].Record(ctx, latency, attributes...)
	}
}

// Logger creates a new Logger middleware that uses Trivela's GoKit default logger.
func Logger(logger logger) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		start := time.Now()

		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		ctx.Next()

		now := time.Now()
		latency := now.Sub(start).Milliseconds()
		clientIP := ctx.ClientIP()
		statusCode := ctx.Writer.Status()
		responseSize := ctx.Writer.Size()

		msg := "[GIN] [%s] [%s] %s %s (%d | %dms | %dB)"
		args := []interface{}{
			now.UTC().Format(time.RFC3339),
			clientIP,
			method,
			path,
			statusCode,
			latency,
			responseSize,
		}

		switch {
		case statusCode < 300:
			logger.Info(ctx, msg, args...)
			break
		case statusCode < 500:
			logger.Warning(ctx, msg, args...)
			break
		default:
			logger.Error(ctx, errors.New(fmt.Sprintf(msg, args...)))
		}
	}
}

// ErrorHandler handles request errors, standardizing how error responses payloads should be served.
func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) == 0 {
		return
	}

	if len(ctx.Errors) == 1 {
		res := newErrorResponse(ctx, ctx.Errors[0].Err)
		ctx.JSON(res.StatusCode(), res)
		return
	}

	errs := []error{}
	for _, err := range ctx.Errors {
		errs = append(errs, err.Err)
	}

	res := newErrorListResponse(ctx, errs...)
	ctx.JSON(res.StatusCode(), res)
}
