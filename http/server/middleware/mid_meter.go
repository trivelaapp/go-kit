package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
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

	requestCounter, _ := meter.SyncInt64().Counter(
		requestCounterName,
		instrument.WithDescription("Counts requests to a an specific (service, http_method, path, status_code) set"),
		instrument.WithUnit(unit.Dimensionless),
	)
	requestBytesCounter, _ := meter.SyncInt64().Counter(
		requestBytesCounterName,
		instrument.WithDescription("Counts the total bytes present in requests to an specific (service, http_method, path, status_code) set"),
		instrument.WithUnit(unit.Bytes),
	)
	responseBytesCounter, _ := meter.SyncInt64().Counter(
		responseBytesCounterName,
		instrument.WithDescription("Counts the total bytes present in response from an specific (service, http_method, path, status_code) set"),
		instrument.WithUnit(unit.Bytes),
	)
	responseLatency, _ := meter.SyncInt64().Histogram(
		responseLatencyHistogramName,
		instrument.WithDescription("Measures response latencies to an specific (service, http_method, path, status_code) set"),
		instrument.WithUnit(unit.Milliseconds),
	)

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
			attribute.String("application_name", applicationName),
			attribute.String("path", ctx.FullPath()),
			attribute.String("method", ctx.Request.Method),
			attribute.Int("status_code", ctx.Writer.Status()),
		}

		var requestSize, responseSize int64
		if ctx.Request.ContentLength > 0 {
			requestSize = ctx.Request.ContentLength
		}
		if ctx.Writer.Size() > 0 {
			responseSize = int64(ctx.Writer.Size())
		}

		counters[requestCounterName].Add(ctx, 1, attributes...)
		counters[requestBytesCounterName].Add(ctx, requestSize, attributes...)
		counters[responseBytesCounterName].Add(ctx, responseSize, attributes...)

		latency := time.Now().Sub(start).Milliseconds()
		histograms[responseLatencyHistogramName].Record(ctx, latency, attributes...)
	}
}
