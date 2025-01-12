package handlers

import (
	"API/internal/models"
	"API/internal/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// GetAllUsers retrieves all users.
// @Summary Get all users
// @Description Retrieves all users from the database.
// @Tags Users
// @Produce json
// @Success 200 {array} models.Users "List of users"
// @Failure 400 {object} models.ResponseError "Failed to retrieve users"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseErr("could not select users"))
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser creates a new user.
// @Summary Create a new user
// @Description Adds a new user to the database.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.Users true "User data"
// @Success 200 {object} models.Response "User created successfully"
// @Failure 400 {object} models.ResponseError "Invalid JSON payload"
// @Failure 500 {object} models.ResponseError "Failed to create user"
// @Router /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.Users
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseErr("Invalid JSON payload"))
	}

	if err := h.Service.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseErr("Failed to create user"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Created user",
		Data:    user,
	})
}

// @Summary Delete a user
// @Description Removes a user from the database by ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Response "User deleted successfully"
// @Failure 400 {object} models.ResponseError "Invalid user ID"
// @Failure 500 {object} models.ResponseError "Failed to delete user"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	// Извлекаем ID из параметров пути
	idParam := c.Param("id")

	// Преобразуем ID в int64
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseErr("Invalid user ID"))
	}

	if err := h.Service.DeleteUser(userID); err != nil {
		// return c.JSON(http.StatusInternalServerError, models.ResponseErr("Failed to delete user"))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "successfuly delete user",
		Data:    userID,
	})
}
