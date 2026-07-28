package routes

import (
	"github.com/vikikurnia87/service-order/handlers"

	"github.com/labstack/echo/v5"
)

// registerDayRoutes mendaftarkan endpoint hari pada group /api/v1/days.
func registerDayRoutes(g *echo.Group, h *handlers.DayHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}
