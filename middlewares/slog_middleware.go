package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-utils/common"
	"github.com/vikikurnia87/service-utils/monitoring"
)

func RequestContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			if reqID != "" {
				ctx := context.WithValue(c.Request().Context(), common.ContextRequestID, reqID)
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	}
}

func SlogMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			err := next(c)
			ctx := c.Request().Context()

			statusCode := http.StatusOK
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					statusCode = he.Code
				} else if strings.Contains(strings.ToLower(err.Error()), "not found") {
					statusCode = http.StatusNotFound
				} else {
					statusCode = http.StatusInternalServerError
				}
			}

			level := slog.LevelInfo
			if statusCode >= 500 {
				level = slog.LevelError
			} else if statusCode >= 400 {
				level = slog.LevelWarn
			}

			monitoring.ContextLogger(ctx, logger).Log(ctx, level, "HTTP request completed",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Path()),
				slog.Int("status", statusCode),
				slog.Duration("duration", time.Since(start)),
			)
			return err
		}
	}
}
