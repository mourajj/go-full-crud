package httpserver

import (
	"codelit/internal/repositories"

	"github.com/labstack/echo"
)

type API struct {
	dbRepo repositories.MemberRepository
}

func RegisterRoutes(e *echo.Echo, dbRepo repositories.MemberRepository) {
	api := &API{
		dbRepo: dbRepo,
	}

	e.GET("/members", httpGetMembers)
	e.GET("/members/:id", api.GetMemberByID)
	e.POST("/members", api.CreateMember)
	e.PUT("/members/:id", api.UpdateMember)
	e.DELETE("/members/:id", api.DeleteMember)
}
œ
