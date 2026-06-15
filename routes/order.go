package routes

import (
	"github.com/vikikurnia87/service-order/handlers"

	"github.com/labstack/echo/v5"
)

// registerOrderRoutes mendaftarkan endpoint order pada group /api/v1/order.
// GET (list/detail) + POST (create) dulu; PUT/DELETE & snapshot/jadwal menyusul.
func registerOrderRoutes(g *echo.Group, h *handlers.OrderHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
}
