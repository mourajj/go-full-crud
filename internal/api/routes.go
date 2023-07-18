package api

import (
	"codelit/internal/models"
	"codelit/internal/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type API struct {
	dbRepo repositories.MemberRepository
}

func RegisterRoutes(e *echo.Echo, dbRepo repositories.MemberRepository) {
	api := &API{
		dbRepo: dbRepo,
	}

	e.GET("/members", api.GetMembers)
	e.GET("/members/:id", api.GetMemberByID)
	e.POST("/members", api.CreateMember)
	e.PUT("/members/:id", api.UpdateMember)
	e.DELETE("/members/:id", api.DeleteMember)
}

func (api *API) GetMembers(c echo.Context) error {
	members, err := api.dbRepo.GetAllMembers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, members)
}

func (api *API) GetMemberByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid member ID")
	}

	member, err := api.dbRepo.GetMemberByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Member not found")
	}
	return c.JSON(http.StatusOK, member)
}

func (api *API) CreateMember(c echo.Context) error {
	member := new(models.Member)
	if err := c.Bind(member); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid member data")
	}

	checkMemberType(member, c)

	err := api.dbRepo.CreateMember(member)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, member)
}

func (api *API) UpdateMember(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid member ID")
	}

	member := new(models.Member)
	if err := c.Bind(member); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid member data")
	}

	checkMemberType(member, c)
	member.ID = id

	_, err = api.dbRepo.GetMemberByID(member.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Member does not exist")
	} else {
		api.dbRepo.UpdateMember(member)
		return c.JSON(http.StatusOK, member)
	}
}

func (api *API) DeleteMember(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid member ID")
	}

	_, err = api.dbRepo.GetMemberByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	} else {
		api.dbRepo.DeleteMember(id)
		return c.NoContent(http.StatusNoContent)
	}
}

func checkMemberType(member *models.Member, c echo.Context) error {

	if member.Name == "" {
		return c.JSON(http.StatusBadRequest, "Members must have a name")
	}

	if member.Type == "contractor" {
		if member.Duration == 0 {
			return c.JSON(http.StatusBadRequest, "Contractors must have a duration")
		}
		if member.Role != "" {
			return c.JSON(http.StatusBadRequest, "Contractors must not have a role")
		}
	} else if member.Type == "employee" {
		if member.Role == "" {
			return c.JSON(http.StatusBadRequest, "Employees must have a role")
		}
		if member.Duration != 0 {
			return c.JSON(http.StatusBadRequest, "Employees must not have a duration")
		}

		member.Duration = 0 //setting a default value for employees
	} else {
		return c.JSON(http.StatusBadRequest, "Invalid member type, please use 'contractor' or 'employee'")
	}

	return nil
}
