package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-order/handlers"
)

// registerCategoryRoutes mendaftarkan endpoint kategori pada group /api/v1/categories.
func registerCategoryRoutes(g *echo.Group, h *handlers.CategoryHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
