package handlers

import (
	"API/internal/models"
	"API/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type UserSegmentHandler struct {
	service *services.UserSegmentService
}

func NewUserSegmentHandler(service *services.UserSegmentService) *UserSegmentHandler {
	return &UserSegmentHandler{service: service}
}

// GetUserSegments retrieves segments for a specific user.
// @Summary Get segments for a user
// @Description Retrieves all segments associated with a user by user ID.
// @Tags UserSegments
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} models.UserSegment "List of user segments"
// @Failure 400 {object} models.ResponseError "Invalid user ID"
// @Failure 500 {object} models.ResponseError "Failed to retrieve user segments"
// @Router /user_segments/{user_id} [get]
func (h *UserSegmentHandler) GetUserSegments(c echo.Context) error {
	userIDParam := c.Param("user_id")

	// Преобразуем его в int64
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseErr("Invalid user ID"))
	}

	segments, err := h.service.GetUserSegments(int64(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseErr("failed to get user segments", err))
	}

	return c.JSON(http.StatusOK, segments)
}

// GetAllUserSegments retrieves all user-segment relationships.
// @Summary Retrieve all user-segment relationships
// @Description Fetches all user-to-segment mappings stored in the database.
// @Tags UserSegments
// @Produce json
// @Success 200 {array} models.UserSegment "List of user-segment relationships"
// @Failure 500 {object} models.ResponseError "Failed to retrieve user segments"
// @Router /user_segments [get]
func (h *UserSegmentHandler) GetAllUserSegments(c echo.Context) error {
	segments, err := h.service.GetAllUserSegments()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseErr("failed to get all user segments", err))
	}

	return c.JSON(http.StatusOK, segments)
}

// UpdateUserSegments modifies a user's segments.
// @Summary Update a user's segments
// @Description Adds or removes segments associated with a user.
// @Tags UserSegments
// @Accept json
// @Produce json
// @Param userSegment body models.UpdateSegmentsRequest true "Segments to add or remove"
// @Success 200 {object} models.Response "User segments updated successfully"
// @Failure 400 {object} models.ResponseError "Invalid request payload"
// @Failure 500 {object} models.ResponseError "Failed to update user segments"
// @Router /user_segments [patch]
func (h *UserSegmentHandler) UpdateUserSegments(c echo.Context) error {
	var req models.UpdateSegmentsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseErr("invalid request body"))
	}

	var ttl *time.Time

	if req.TTL != nil {
		parsedTTL, err := time.Parse(time.RFC3339, *req.TTL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ResponseErr("invalid TTL format"))
		}

		ttl = &parsedTTL
	}

	if err := h.service.UpdateUserSegments(req.UserID, req.AddSegments, req.DeleteSegments, ttl); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseErr("failed to update user segments", err))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "User segments updated successfully",
	})
}
