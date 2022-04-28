package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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
		clientIP := ctx.ClientIP()
		statusCode := ctx.Writer.Status()
		responseSize := ctx.Writer.Size()

		msg := "[GIN] [%s] [%s] %s %s (%d | %dms | %dB)"
		args := []any{
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
