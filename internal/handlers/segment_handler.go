package handlers

import (
	"API/internal/models"
	"API/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SegmentHandler struct {
	segmentService *services.SegmentService
}

func NewSegmentHandler(segmentService *services.SegmentService) *SegmentHandler {
	return &SegmentHandler{segmentService: segmentService}
}

// GetAllSegments retrieves all segments from the database.
// @Summary Retrieve all segments
// @Description Fetches a list of all segments stored in the database.
// @Tags Segments
// @Produce json
// @Success 200 {array} models.Segments "List of segments"
// @Failure 500 {object} models.ResponseError "Failed to retrieve segments"
// @Router /segments [get]
func (h *SegmentHandler) GetAllSegments(c echo.Context) error {
	segments, err := h.segmentService.GetAllSegments()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseErr("could not retrieve segments"))
	}

	return c.JSON(http.StatusOK, segments)
}

// CreateSegment creates a new segment
// @Summary Create a new segment
// @Description Adds a new segment to the database.
// @Tags Segments
// @Accept json
// @Produce json
// @Param segment body models.SegmentRequest true "Segment data"
// @Success 200 {object} models.Response "Segment created successfully"
// @Failure 400 {object} models.ResponseError "Invalid slug"
// @Failure 500 {object} models.ResponseError "Failed to create segment"
// @Router /segments [post]
func (h *SegmentHandler) CreateSegment(c echo.Context) error {
	var segment models.SegmentRequest
	if err := c.Bind(&segment); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseErr("invalid slug"))
	}

	if err := h.segmentService.CreateSegment(segment.Slug); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Segment created successfully",
	})
}

// DeleteSegment removes a segment by its slug.
// @Summary Delete a segment
// @Description Deletes a segment from the system using the provided slug.
// @Tags Segments
// @Accept json
// @Produce json
// @Param segment body models.SegmentRequest true "Segment data"
// @Success 200 {object} models.Response "Segment deleted successfully"
// @Failure 400 {object} models.ResponseError "Invalid slug"
// @Failure 500 {object} models.ResponseError "Failed to delete segment"
// @Router /segments [delete]
func (h *SegmentHandler) DeleteSegment(c echo.Context) error {
	var segment models.SegmentRequest
	if err := c.Bind(&segment); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseErr("invalid slug"))
	}

	if err := h.segmentService.DeleteSegment(segment.Slug); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Segment deleted successfully",
	})
}
