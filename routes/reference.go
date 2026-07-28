package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-order/handlers"
)

// registerOrderStatusRoutes mendaftarkan endpoint status order pada group /api/v1/order-statuses.
func registerOrderStatusRoutes(g *echo.Group, h *handlers.OrderStatusHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}

// registerOrderPriorityRoutes mendaftarkan endpoint prioritas order pada group /api/v1/order-priorities.
func registerOrderPriorityRoutes(g *echo.Group, h *handlers.OrderPriorityHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}

// registerScheduleRoutes mendaftarkan endpoint jadwal pada group /api/v1/schedules.
func registerScheduleRoutes(g *echo.Group, h *handlers.ScheduleHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}

// registerDayRoutes mendaftarkan endpoint hari pada group /api/v1/days.
func registerDayRoutes(g *echo.Group, h *handlers.DayHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}

// registerDateRoutes mendaftarkan endpoint tanggal pada group /api/v1/dates.
func registerDateRoutes(g *echo.Group, h *handlers.DateHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}

// registerMasterDateRoutes mendaftarkan endpoint kalender pada group /api/v1/master-dates.
func registerMasterDateRoutes(g *echo.Group, h *handlers.MasterDateHandler) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
}
