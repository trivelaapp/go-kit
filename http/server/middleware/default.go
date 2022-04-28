package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// DefaultInput encapsulates inputs to a Default call.
type DefaultInput struct {
	ApplicationName string
	Logger          LogProvider
}

// Default returns a list of middlewares usually used by most applications. It includes:
// - Disaster Recovery (from panics)
// - Tracing
// - Request Logging
// - Error Handling
func Default(in DefaultInput) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		Meter(in.ApplicationName),
		gin.Recovery(),
		otelgin.Middleware(in.ApplicationName),
		Logger(in.Logger),
		ErrorHandler,
	}
}
