package httpserver

import (
	"codelit/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// HTTPHandler Wrapper API around an echo.Handler function to support dependency injection
type HTTPHandler interface {
	Handler() echo.HandlerFunc
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

	notValid, err := checkMemberType(member, c)
	if notValid {
		return err
	}

	err = api.dbRepo.CreateMember(member)
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

	notValid, err := checkMemberType(member, c)
	if notValid {
		return err
	}
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

func checkMemberType(member *models.Member, c echo.Context) (bool, error) {

	if member.Name == "" {
		return true, c.JSON(http.StatusBadRequest, "Members must have a name")
	}

	if member.Type == "contractor" {
		if member.Duration == 0 {
			return true, c.JSON(http.StatusBadRequest, "Contractors must have a duration")
		}
		if member.Role != "" {
			return true, c.JSON(http.StatusBadRequest, "Contractors must not have a role")
		}
	} else if member.Type == "employee" {
		if member.Role == "" {
			return true, c.JSON(http.StatusBadRequest, "Employees must have a role")
		}
		if member.Duration != 0 {
			return true, c.JSON(http.StatusBadRequest, "Employees must not have a duration")
		}

		member.Duration = 0 //setting a default value for employees
	} else {
		return true, c.JSON(http.StatusBadRequest, "Invalid member type, please use 'contractor' or 'employee'")
	}

	return false, nil
}
