package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

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

type errorPayload struct {
	Code    errors.CodeType `json:"code,omitempty"`
	Message string          `json:"message"`
}

type errorResponse struct {
	status int

	TraceID string       `json:"trace_id,omitempty"`
	Err     errorPayload `json:"error"`
}

func newErrorResponse(ctx context.Context, err error) errorResponse {
	return errorResponse{
		TraceID: getTraceID(trace.SpanFromContext(ctx)),
		status:  kindToHTTPStatusCode(errors.Kind(err)),
		Err: errorPayload{
			Code:    errors.Code(err),
			Message: err.Error(),
		},
	}
}

func (e errorResponse) StatusCode() int {
	return e.status
}

type errorListResponse struct {
	status int

	TraceID string         `json:"trace_id,omitempty"`
	Errs    []errorPayload `json:"errors"`
}

func newErrorListResponse(ctx context.Context, errs ...error) errorListResponse {
	if len(errs) == 0 {
		return errorListResponse{
			status: 500,
		}
	}

	errsPayload := []errorPayload{}
	for _, err := range errs {
		errsPayload = append(errsPayload, errorPayload{
			Code:    errors.Code(err),
			Message: err.Error(),
		})
	}

	return errorListResponse{
		TraceID: getTraceID(trace.SpanFromContext(ctx)),
		status:  kindToHTTPStatusCode(errors.Kind(errs[0])),
		Errs:    errsPayload,
	}
}

func (e errorListResponse) StatusCode() int {
	return e.status
}

func getTraceID(span trace.Span) string {
	if !span.SpanContext().HasTraceID() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

func kindToHTTPStatusCode(kind errors.KindType) int {
	switch kind {
	case errors.KindInvalidInput:
		return http.StatusBadRequest
	case errors.KindUnauthenticated:
		return http.StatusUnauthorized
	case errors.KindUnauthorized:
		return http.StatusForbidden
	case errors.KindNotFound:
		return http.StatusNotFound
	case errors.KindConflict:
		return http.StatusConflict
	case errors.KindUnexpected:
		return http.StatusInternalServerError
	case errors.KindInternal:
		return http.StatusInternalServerError
	case errors.KindResourceExhausted:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
