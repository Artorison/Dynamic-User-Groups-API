package app

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(router *echo.Echo, container *DIContainer) {
	usersRoutes(router, container)
	segmentsRoutes(router, container)
	userSegmentsRoutes(router, container)
	RegisterStaticFiles(router)

}

func usersRoutes(router *echo.Echo, container *DIContainer) {
	users := router.Group("users")
	users.GET("", container.UserHandler.GetAllUsers)
	users.POST("", container.UserHandler.CreateUser)
	users.DELETE("/:id", container.UserHandler.DeleteUser)
}

func segmentsRoutes(router *echo.Echo, container *DIContainer) {
	segments := router.Group("/segments")
	segments.GET("", container.SegmentHandler.GetAllSegments)
	segments.POST("", container.SegmentHandler.CreateSegment)
	segments.DELETE("", container.SegmentHandler.DeleteSegment)
}

func userSegmentsRoutes(router *echo.Echo, container *DIContainer) {
	userSegments := router.Group("/user_segments")
	userSegments.GET("/:user_id", container.UserSegmentHandler.GetUserSegments)
	userSegments.GET("", container.UserSegmentHandler.GetAllUserSegments)
	userSegments.PATCH("", container.UserSegmentHandler.UpdateUserSegments)
	userSegments.GET("/history/:user_id", container.UserSegmentHistoryHandler.GenerateHistoryCSV)

}

func RegisterStaticFiles(router *echo.Echo) {
	router.Static("/csv_reports", "./csv_reports")
}
