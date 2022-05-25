package meter

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	requestCounterName           = "rpc.grpc.server.request_counter"
	responseLatencyHistogramName = "rpc.grpc.server.duration"
)

// UnaryServerInterceptor returns a new unary interceptor suitable for request metrics handling.
func UnaryServerInterceptor(applicationName string) grpc.UnaryServerInterceptor {
	counters, histograms := loadMeters(applicationName)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		start := time.Now()

		resp, err = handler(ctx, req)

		attributes := []attribute.KeyValue{
			attribute.String("application_name", applicationName),
			attribute.String("method", info.FullMethod),
			attribute.String("status_code", status.Code(err).String()),
		}

		counters[requestCounterName].Add(ctx, 1, attributes...)

		latency := time.Now().Sub(start).Milliseconds()
		histograms[responseLatencyHistogramName].Record(ctx, latency, attributes...)

		return
	}
}

// StreamServerInterceptor returns a new stream interceptor suitable for request metrics handling.
func StreamServerInterceptor(applicationName string) grpc.StreamServerInterceptor {
	counters, histograms := loadMeters(applicationName)

	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {

		ctx := ss.Context()
		start := time.Now()

		err = handler(srv, ss)

		attributes := []attribute.KeyValue{
			attribute.String("application_name", applicationName),
			attribute.String("method", info.FullMethod),
			attribute.String("status_code", status.Code(err).String()),
		}

		counters[requestCounterName].Add(ctx, 1, attributes...)

		latency := time.Now().Sub(start).Milliseconds()
		histograms[responseLatencyHistogramName].Record(ctx, latency, attributes...)

		return
	}
}

func loadMeters(applicationName string) (map[string]syncint64.Counter, map[string]syncint64.Histogram) {
	meter := global.MeterProvider().Meter(applicationName)

	requestCounter, _ := meter.SyncInt64().Counter(
		requestCounterName,
		instrument.WithDescription("Counts requests to a an specific (service, rpc_method, status_code) set"),
		instrument.WithUnit(unit.Dimensionless),
	)
	responseLatency, _ := meter.SyncInt64().Histogram(
		responseLatencyHistogramName,
		instrument.WithDescription("Measures response latencies to an specific (service, rpc_method, status_code) set"),
		instrument.WithUnit(unit.Milliseconds),
	)

	counters := map[string]syncint64.Counter{
		requestCounterName: requestCounter,
	}

	histograms := map[string]syncint64.Histogram{
		responseLatencyHistogramName: responseLatency,
	}

	return counters, histograms
}
