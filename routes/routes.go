package routes

import (
	"net/http"

	"github.com/vikikurnia87/service-order/clients"
	"github.com/vikikurnia87/service-order/container"
	"github.com/vikikurnia87/service-order/middlewares"

	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-utils/response"
)

// SetupRouter membangun Echo + middleware global + route. userClient dipakai
// auth middleware untuk ValidateToken ke service-user.
func SetupRouter(c *container.Container, userClient *clients.UserClient) http.Handler {
	e := echo.New()

	e.Use(middlewares.RequestContextMiddleware())
	e.Use(middlewares.APMTransactionMiddleware())
	e.Use(middlewares.SlogMiddleware(c.Logger))
	e.Use(middlewares.RecoverWithAPM())

	e.GET("/health", func(ctx *echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	// Konvensi: path tak match route apa pun -> 404 (envelope standar).
	e.RouteNotFound("/*", response.RouteNotFoundHandler)

	registerAPIRoutes(e, c, userClient)
	return e
}

// registerAPIRoutes orchestrator: group /api/v1 (auth) lalu delegasi ke registrar
// per-entitas. Registrar ditambahkan pada Fase 1+ (order_status/priority/.../order).
func registerAPIRoutes(e *echo.Echo, c *container.Container, userClient *clients.UserClient) {
	api := e.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware(userClient, c.Logger))
	_ = api
}
