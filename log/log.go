package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/trivelaapp/go-kit/errors"
	"github.com/trivelaapp/go-kit/log/format"
)

// LoggerParams defines the dependencies of a Logger.
type LoggerParams struct {
	Level      string
	Formatter  LogFormatter
	Attributes format.LogAttributeSet
}

// Logger is the structure responsible for log data.
type Logger struct {
	level      Level
	formatter  LogFormatter
	attributes format.LogAttributeSet
	now        func() time.Time
}

// NewLogger constructs a new Logger instance.
func NewLogger(params LoggerParams) *Logger {
	logger := &Logger{
		level:      levelStringValueMap[params.Level],
		attributes: params.Attributes,
		formatter:  params.Formatter,
		now:        time.Now,
	}

	if logger.level < LevelCritical || logger.level > LevelDebug {
		logger.level = LevelInfo
	}

	if logger.formatter == nil {
		logger.formatter = format.NewDefault()
	}

	return logger
}

// Debug logs debug data.
func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelDebug {
		l.printMsg(ctx, fmt.Sprintf(msg, args...), LevelDebug)
	}
}

// Info logs info data.
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelInfo {
		l.printMsg(ctx, fmt.Sprintf(msg, args...), LevelInfo)
	}
}

// Warning logs warning data.
func (l Logger) Warning(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelWarning {
		l.printMsg(ctx, fmt.Sprintf(msg, args...), LevelWarning)
	}
}

// Error logs error data. It increases error counter metrics.
func (l Logger) Error(ctx context.Context, err error) {
	if l.level >= LevelError {
		l.printError(ctx, err, LevelError)
	}
}

// Critical logs critical data. It increases error counter metrics.
func (l Logger) Critical(ctx context.Context, err error) {
	if l.level >= LevelCritical {
		l.printError(ctx, err, LevelCritical)
	}
}

// Fatal logs critical data and exists current program execution.
func (l Logger) Fatal(ctx context.Context, err error) {
	if l.level >= LevelCritical {
		l.printError(ctx, err, LevelCritical)
		os.Exit(1)
	}
}

// JSON logs JSON data.
// By default, it logs with Debug level, but it can be overwritten with the optional logLevel parameter.
func (l Logger) JSON(ctx context.Context, data any, logLevel ...Level) {
	level := LevelDebug
	if len(logLevel) > 0 {
		level = logLevel[0]
	}

	if l.level < level {
		return
	}

	if _, err := json.Marshal(data); err != nil {
		l.Error(ctx, errors.New("could not marshal payload to JSON format").WithRootError(err))
		return
	}

	l.printJSON(ctx, data, level)
}

func (l Logger) printMsg(ctx context.Context, msg string, level Level) {
	payload := l.formatter.Format(ctx, format.LogInput{
		Level:      level.String(),
		Message:    msg,
		Attributes: l.attributes,
		Timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))
}

func (l Logger) printJSON(ctx context.Context, jsonData any, level Level) {
	payload := l.formatter.Format(ctx, format.LogInput{
		Level:      level.String(),
		Message:    "JSON data logged",
		Payload:    jsonData,
		Attributes: l.attributes,
		Timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))
}

func (l Logger) printError(ctx context.Context, err error, level Level) {
	payload := l.formatter.Format(ctx, format.LogInput{
		Level:      level.String(),
		Message:    err.Error(),
		Err:        err,
		Attributes: l.attributes,
		Timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))

	counter := errorCounter()
	if counter != nil {
		attrs := []attribute.KeyValue{
			attribute.String("level", level.String()),
		}

		if service := ctx.Value(string(semconv.ServiceNameKey)); service != nil {
			attribute.String(string(semconv.ServiceNameKey), service.(string))
		}

		if version := ctx.Value(string(semconv.ServiceVersionKey)); version != nil {
			attribute.String(string(semconv.ServiceVersionKey), version.(string))
		}

		counter.Add(ctx, 1, attrs...)
	}
}

var errCounter syncint64.Counter

func errorCounter() syncint64.Counter {
	if errCounter != nil {
		return errCounter
	}

	counter, err := global.Meter("trivelaapp.go-kit.errors").SyncInt64().Counter(
		"app.error_counter",
		instrument.WithDescription("Counts errors logged by the application"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil
	}

	errCounter = counter
	return counter
}
