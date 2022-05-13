package server

import (
	"fmt"

	"github.com/trivelaapp/go-kit/errors"
)

var (
	// ErrRequestQueryStringValidation indicates a failure during request query string binding.
	ErrRequestQueryStringValidation errors.CustomError = errors.
					New("request query string validation failed").
					WithKind(errors.KindInvalidInput).
					WithCode("ERR_REQUEST_QUERY_STRING_VALIDATION")

	// ErrRequestBodyValidation indicates a failure during request body binding.
	ErrRequestBodyValidation errors.CustomError = errors.
					New("request query string validation failed").
					WithKind(errors.KindInvalidInput).
					WithCode("ERR_REQUEST_BODY_VALIDATION")
)

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
