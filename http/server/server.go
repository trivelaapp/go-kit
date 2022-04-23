package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/trivelaapp/go-kit/errors"
)

// ServerParams encapsulates the necessary inputs to initialize a Server.
type ServerParams struct {
	Port            int
	ApplicationName string
	Logger          logger
}

// Server is an HTTPServer object that can serve HTTP requests.
// It is based on httpfast implementation.
type Server struct {
	port            int
	applicationName string
	router          *gin.Engine
	logger          logger
}

// New creates a new instance of Server.
func New(params ServerParams) (*Server, error) {
	if params.Port == 0 {
		return nil, errors.NewMissingRequiredDependency("PORT")
	}

	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.Logger == nil {
		return nil, errors.NewMissingRequiredDependency("Logger")
	}

	return &Server{
		port:            params.Port,
		applicationName: params.ApplicationName,
		router:          gin.New(),
		logger:          params.Logger,
	}, nil
}

// MustNew creates a new instance of Server.
// It panics if any error is found.
func MustNew(params ServerParams) *Server {
	server, err := New(params)
	if err != nil {
		panic(err)
	}

	return server
}

// RouterEmpty returns server's router with no middleware installed.
func (s Server) RouterEmpty() *gin.Engine {
	return s.router
}

// RouterDefault returns server's router with default middlewares already installed.
// The default middlewares implements:
// - Disaster Recovery (from panics)
// - Tracing
// - Request Logging
// - Error Handling
func (s *Server) RouterDefault() *gin.Engine {
	router := s.router

	router.Use(
		Meter(s.applicationName),
		gin.Recovery(),
		otelgin.Middleware(s.applicationName),
		Logger(s.logger),
		ErrorHandler,
	)

	return router
}

// Run listen and serves HTTP requests at the Port specified in Server construction
func (s Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}
