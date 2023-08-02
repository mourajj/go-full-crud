package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Route defining an HTTP endpoint
type Route struct {
	HttpMethod     string
	HttpPath       string
	HttpHandler    HTTPHandler
	HttpMiddleware []HttpMiddleware
}

func NewRoute(httpMethod string, httpPath string, httpHandler HTTPHandler, middleware ...HttpMiddleware) *Route {
	return &Route{
		HttpMethod:     httpMethod,
		HttpPath:       httpPath,
		HttpHandler:    httpHandler,
		HttpMiddleware: middleware,
	}
}

type RouteGroup struct {
	Group          string
	Routes         []*Route
	HttpMiddleware []HttpMiddleware
}

// Adds a Route to the RouteGroup
func (rg *RouteGroup) AddRoute(route *Route) {
	rg.Routes = append(rg.Routes, route)
}

// Factory to create a RouteGroup
func NewRouteGroup(group string, middleware ...HttpMiddleware) *RouteGroup {
	return &RouteGroup{
		Group:          group,
		Routes:         make([]*Route, 0),
		HttpMiddleware: middleware,
	}
}

// The EchoServer is a wrapper around the Echo library that adds start/stop behavior
// and internally configures routes
type EchoServer struct {
	shutdownTimeout time.Duration
	address         string
	server          *echo.Echo
}

// Creates a new EchoServer
func New(port int, middleware ...HttpMiddleware) *EchoServer {
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

	return &echoServer
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

// Adds a HTTP Route
func (e *EchoServer) AddRoute(route *Route) {
	e.server.Add(
		route.HttpMethod,
		route.HttpPath,
		route.HttpHandler.Handler(),
		mapMiddleware(route.HttpMiddleware)...)
}

// Adds a HTTP RouteGroup
func (e *EchoServer) AddRouteGroup(routeGroup *RouteGroup) {
	group := e.server.Group(routeGroup.Group, mapMiddleware(routeGroup.HttpMiddleware)...)
	for _, route := range routeGroup.Routes {
		group.Add(
			route.HttpMethod,
			route.HttpPath,
			route.HttpHandler.Handler(),
			mapMiddleware(route.HttpMiddleware)...)
	}
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
