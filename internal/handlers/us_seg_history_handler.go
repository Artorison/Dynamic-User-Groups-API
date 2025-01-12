package handlers

import (
	"API/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserSegmentHistoryHandler struct {
	Service *services.UserSegmentHistoryService
}

func NewUserSegmentHistoryHandler(service *services.UserSegmentHistoryService) *UserSegmentHistoryHandler {
	return &UserSegmentHistoryHandler{Service: service}
}

// GenerateHistoryCSV handles the request to generate a user history CSV file.
// @Summary Generate User History CSV
// @Description Generate a CSV file containing the user's segment history for a specific month.
// @Tags UserSegmentHistory
// @Param user_id path int true "User ID"
// @Param date query string true "Year-Month in YYYY-MM format"
// @Produce text/plain
// @Success 200 {string} string "URL to the generated CSV file"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user_segments/history/{user_id} [get]
func (h *UserSegmentHistoryHandler) GenerateHistoryCSV(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid user ID")
	}

	date := c.QueryParam("date")
	if date == "" {
		return c.JSON(http.StatusBadRequest, "date parameter is required")
	}

	fileName, err := h.Service.GenerateUserHistoryCSV(userID, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	url := fmt.Sprintf("http://%s/csv_reports/%s", c.Request().Host, fileName)
	return c.String(http.StatusOK, url)
}
