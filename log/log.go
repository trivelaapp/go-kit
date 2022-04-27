package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

// Level indicates the severity of the data being logged.
type Level int

const (
	// LevelCritical alerts about severe problems. Most of the time, needs some human intervention ASAP.
	LevelCritical Level = iota + 1
	// LevelError alerts about events that are likely to cause problems.
	LevelError
	// LevelWarning warns about events the might cause problems to the system.
	LevelWarning
	// LevelInfo are routine information.
	LevelInfo
	// LevelDebug are debug or trace information.
	LevelDebug
)

var levelStringValueMap = map[string]Level{
	"CRITICAL": LevelCritical,
	"ERROR":    LevelError,
	"WARNING":  LevelWarning,
	"INFO":     LevelInfo,
	"DEBUG":    LevelDebug,
}

// String returns the name of the LogLevel.
func (l Level) String() string {
	return []string{
		"CRITICAL",
		"ERROR",
		"WARNING",
		"INFO",
		"DEBUG",
	}[l-1]
}

// LoggerParams defines the dependencies of a Logger.
type LoggerParams struct {
	Level      string
	Attributes LogAttributeSet
}

// Logger is the structure responsible for log data.
type Logger struct {
	level          Level
	payloadBuilder payloadBuilder
	attributes     LogAttributeSet
	now            func() time.Time
}

// NewLogger constructs a new Logger instance.
func NewLogger(in LoggerParams) *Logger {
	logger := &Logger{
		level:          levelStringValueMap[in.Level],
		payloadBuilder: defaultLogBuilder{},
		attributes:     in.Attributes,
		now:            time.Now,
	}

	if logger.level < LevelCritical || logger.level > LevelDebug {
		logger.level = LevelInfo
	}

	return logger
}

func (l Logger) WithTraceFormat() Logger {
	l.payloadBuilder = defaultWithTraceLogBuilder{}
	return l
}

func (l Logger) WithGCPCloudLoggingFormat(projectID string) Logger {
	l.payloadBuilder = gcpCloudLoggingWithTraceLogBuilder{projectID: projectID}
	return l
}

// Debug logs debug data.
func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelDebug {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelDebug)
	}
}

// Info logs info data.
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelInfo {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelInfo)
	}
}

// Warning logs warning data.
func (l Logger) Warning(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= LevelWarning {
		l.print(ctx, fmt.Sprintf(msg, args...), LevelWarning)
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

type logInput struct {
	level      Level
	message    string
	err        error
	attributes LogAttributeSet
	timestamp  time.Time
}

func (l Logger) print(ctx context.Context, msg string, level Level) {
	payload := l.payloadBuilder.build(ctx, logInput{
		level:      level,
		message:    msg,
		attributes: l.attributes,
		timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))
}

func (l Logger) printError(ctx context.Context, err error, level Level) {
	payload := l.payloadBuilder.build(ctx, logInput{
		level:      level,
		message:    err.Error(),
		err:        err,
		attributes: l.attributes,
		timestamp:  l.now(),
	})

	data, _ := json.Marshal(payload)
	fmt.Println(string(data))

	counter := errorCounter()
	if counter != nil {
		counter.Add(ctx, 1, attribute.String("level", level.String()))
	}
}

var errCounter syncint64.Counter

func errorCounter() syncint64.Counter {
	if errCounter != nil {
		return errCounter
	}

	counter, err := global.Meter("trivelaapp.go-kit.errors").SyncInt64().Counter("app.error_counter")
	if err != nil {
		return nil
	}

	errCounter = counter
	return counter
}
