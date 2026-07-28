package routes

import (
	"github.com/vikikurnia87/service-order/handlers"

	"github.com/labstack/echo/v5"
)

// registerScheduleRoutes mendaftarkan endpoint jadwal pada group /api/v1/schedules.
func registerScheduleRoutes(g *echo.Group, h *handlers.ScheduleHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}
