// Package middlewares — auth via gRPC ke service-user (DRY: tidak menyimpan
// logika JWT sendiri), + APM/log middleware.
package middlewares

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/vikikurnia87/service-order/auth"
	"github.com/vikikurnia87/service-order/clients"

	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-utils/monitoring"
)

// AuthMiddleware memvalidasi Bearer token ke service-user (gRPC ValidateToken)
// lalu menaruh identitas+tenant ke Echo context.
func AuthMiddleware(userClient *clients.UserClient, logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Request().Context()
			token := extractBearer(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			resp, err := userClient.ValidateToken(ctx, token)
			if err != nil {
				monitoring.ContextLogger(ctx, logger).ErrorContext(ctx, "auth: validate token gRPC failed",
					slog.String("error", err.Error()))
				return echo.NewHTTPError(http.StatusServiceUnavailable, "auth service unavailable")
			}
			if !resp.GetValid() {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			// Tenant (non-system) wajib punya company aktif.
			if !resp.GetSystem() && resp.GetCompanyUuid() == "" {
				return echo.NewHTTPError(http.StatusForbidden, "no active company")
			}

			auth.SetContext(c, resp)
			return next(c)
		}
	}
}

// RequirePermission membatasi akses berdasarkan permission user.
func RequirePermission(codenames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			perms := auth.Permissions(c)
			for _, need := range codenames {
				if slices.Contains(perms, need) {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permission")
		}
	}
}

func extractBearer(c *echo.Context) string {
	h := c.Request().Header.Get("Authorization")
	if h == "" || !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
}
