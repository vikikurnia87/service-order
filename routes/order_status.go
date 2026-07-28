package routes

import (
	"github.com/vikikurnia87/service-order/handlers"

	"github.com/labstack/echo/v5"
)

// registerOrderStatusRoutes mendaftarkan endpoint status order pada group /api/v1/order-statuses.
func registerOrderStatusRoutes(g *echo.Group, h *handlers.OrderStatusHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}
