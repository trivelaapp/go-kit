package server

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

type logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warning(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, err error)
	Critical(ctx context.Context, err error)
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
	default:
		return http.StatusInternalServerError
	}
}

type messageResponse struct {
	Message string `json:"message"`
}

// NewMessageResponse creates a a generic message that should be sent to a client of HTTP Server.
func NewMessageResponse(msg string, params ...interface{}) messageResponse {
	return messageResponse{Message: fmt.Sprintf(msg, params...)}
}

type resourceCreatedResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// NewResourceCreatedResponse creates response designed to be sent when a new resource is created in the system.
func NewResourceCreatedResponse(id string) resourceCreatedResponse {
	return resourceCreatedResponse{ID: id, Message: "Resource created successfully!"}
}
