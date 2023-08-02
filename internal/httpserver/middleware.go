package httpserver

import "github.com/labstack/echo"

// HttpMiddleware defines a method to retrieve configurable middleware processors
type HttpMiddleware interface {
	Middleware(next echo.HandlerFunc) echo.HandlerFunc
}
