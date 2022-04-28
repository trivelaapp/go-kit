package server

import (
	"fmt"
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
