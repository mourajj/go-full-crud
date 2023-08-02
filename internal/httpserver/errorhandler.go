package httpserver

import (
	"net/http"

	"github.com/labstack/echo"
)

func errorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError

	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
	}

	if !ctx.Response().Committed {
		if ctx.Request().Method == http.MethodHead {
			err = ctx.NoContent(code)
		} else {
			var errorResponse string
			switch code {
			case 400:
				errorResponse = "error 400"
			case 403:
				errorResponse = "error 403"
			case 500:
				errorResponse = "error 500"
			case 503:
				errorResponse = "error 503"
			default:
				errorResponse = "Generic Exception"
			}
			err = ctx.JSON(code, errorResponse)
		}
		if err != nil {
			ctx.Logger().Error(err)
		}
	}
}
