package routes

import (
	"github.com/vikikurnia87/service-order/handlers"

	"github.com/labstack/echo/v5"
)

// registerOrderPriorityRoutes mendaftarkan endpoint prioritas order pada group /api/v1/order-priorities.
func registerOrderPriorityRoutes(g *echo.Group, h *handlers.OrderPriorityHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}
