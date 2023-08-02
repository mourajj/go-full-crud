package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// The EchoServer is a wrapper around the Echo library that adds start/stop behavior
// and internally configures routes
type EchoServer struct {
	shutdownTimeout time.Duration
	address         string
	server          *echo.Echo
}

type Route struct {
	HttpPath       string
	HttpHandler    echo.HandlerFunc
	HttpMiddleware echo.MiddlewareFunc
}

// Creates a new EchoServer
func New(port int, middleware ...HttpMiddleware) (*EchoServer, *echo.Echo) {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.HTTPErrorHandler = errorHandler

	// Set middleware that is applied to all routes
	middlewareFuncs := make([]echo.MiddlewareFunc, 0)
	for i := range middleware {
		middlewareFuncs = append(middlewareFuncs, middleware[i].Middleware)
	}
	server.Use(middlewareFuncs...)

	echoServer := EchoServer{
		address:         fmt.Sprintf(":%d", port),
		shutdownTimeout: 10 * time.Second,
		server:          server,
	}

	return &echoServer, server
}

// Asynchronously starts the Echo server.
// This returns an error channel that receives a message if an error occurs starting the Echo server.
func (e *EchoServer) Start() <-chan error {
	errorCh := make(chan error, 1)
	go func() {
		defer close(errorCh)

		if err := e.server.Start(e.address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorCh <- err
		}
	}()
	return errorCh
}

// Gracefully stops the Echo server
func (e *EchoServer) Stop() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), e.shutdownTimeout)
	defer cancelFunc()
	return e.server.Shutdown(ctx)
}

func mapMiddleware(middleware []HttpMiddleware) []echo.MiddlewareFunc {
	middlewareFuncs := make([]echo.MiddlewareFunc, 0)
	for i := range middleware {
		middlewareFuncs = append(middlewareFuncs, middleware[i].Middleware)
	}
	return middlewareFuncs
}
