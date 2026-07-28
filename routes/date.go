package routes

import (
	"github.com/vikikurnia87/service-order/handlers"

	"github.com/labstack/echo/v5"
)

// registerDateRoutes mendaftarkan endpoint tanggal pada group /api/v1/dates.
func registerDateRoutes(g *echo.Group, h *handlers.DateHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}
