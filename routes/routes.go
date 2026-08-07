package routes

import (
	"net/http"

	"github.com/vikikurnia87/service-order/clients"
	"github.com/vikikurnia87/service-order/container"
	"github.com/vikikurnia87/service-order/handlers"
	"github.com/vikikurnia87/service-order/middlewares"

	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-utils/httpmiddleware"
	"github.com/vikikurnia87/service-utils/response"
)

// SetupRouter membangun Echo + middleware global + route. userClient dipakai
// auth middleware untuk ValidateToken ke service-user.
func SetupRouter(c *container.Container, userClient *clients.UserClient) http.Handler {
	e := echo.New()

	e.Use(httpmiddleware.RequestContextMiddleware())
	e.Use(httpmiddleware.APMTransactionMiddleware())
	e.Use(httpmiddleware.SlogMiddleware(c.Logger))
	e.Use(httpmiddleware.RecoverWithAPM())

	e.GET("/health", func(ctx *echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	// Konvensi: path tak match route apa pun -> 404 (envelope standar).
	e.RouteNotFound("/*", response.RouteNotFoundHandler)

	registerAPIRoutes(e, c, userClient)
	return e
}

// registerAPIRoutes orchestrator: group /api/v1 (auth) lalu delegasi ke registrar per-entitas.
func registerAPIRoutes(e *echo.Echo, c *container.Container, userClient *clients.UserClient) {
	// Inisialisasi semua handler.
	orderH := handlers.NewOrderHandler(c.OrderService, c.Logger)
	categoryH := handlers.NewCategoryHandler(c.CategoryService, c.Logger)
	orderStatusH := handlers.NewOrderStatusHandler(c.OrderStatusService, c.Logger)
	orderPriorityH := handlers.NewOrderPriorityHandler(c.OrderPriorityService, c.Logger)
	scheduleH := handlers.NewScheduleHandler(c.ScheduleService, c.Logger)
	dayH := handlers.NewDayHandler(c.DayService, c.Logger)
	dateH := handlers.NewDateHandler(c.DateService, c.Logger)
	masterDateH := handlers.NewMasterDateHandler(c.MasterDateService, c.Logger)

	api := e.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware(userClient, c.Logger))

	// Order (inti).
	registerOrderRoutes(api.Group("/order"), orderH)

	// Kategori (ter-scope tenant).
	registerCategoryRoutes(api.Group("/categories"), categoryH)

	// Referensi global (read-only).
	registerOrderStatusRoutes(api.Group("/order-statuses"), orderStatusH)
	registerOrderPriorityRoutes(api.Group("/order-priorities"), orderPriorityH)
	registerScheduleRoutes(api.Group("/schedules"), scheduleH)
	registerDayRoutes(api.Group("/days"), dayH)
	registerDateRoutes(api.Group("/dates"), dateH)
	registerMasterDateRoutes(api.Group("/master-dates"), masterDateH)
}
