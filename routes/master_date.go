package routes

import (
	"github.com/vikikurnia87/service-order/handlers"

	"github.com/labstack/echo/v5"
)

// registerMasterDateRoutes mendaftarkan endpoint kalender pada group /api/v1/master-dates.
func registerMasterDateRoutes(g *echo.Group, h *handlers.MasterDateHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}
