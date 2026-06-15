package middlewares

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"go.elastic.co/apm/v2"
)

// APMTransactionMiddleware membuat APM transaction per request (root span).
func APMTransactionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			if tx := apm.TransactionFromContext(req.Context()); tx != nil {
				return next(c)
			}

			name := fmt.Sprintf("%s %s", req.Method, c.Path())
			tx := apm.DefaultTracer().StartTransaction(name, "request")
			defer tx.End()

			ctx := apm.ContextWithTransaction(req.Context(), tx)
			c.SetRequest(req.WithContext(ctx))

			tx.Context.SetHTTPRequest(req)
			tx.Context.SetLabel("route", c.Path())
			tx.Context.SetLabel("method", req.Method)

			err := next(c)

			status := http.StatusOK
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				} else {
					status = http.StatusInternalServerError
				}
				e := apm.CaptureError(ctx, err)
				e.SetTransaction(tx)
				e.Send()
				tx.Result = "error"
			} else {
				tx.Result = "success"
			}
			tx.Context.SetHTTPStatusCode(status)
			return err
		}
	}
}

// RecoverWithAPM recover dari panic + kirim ke APM.
func RecoverWithAPM() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					if tx := apm.TransactionFromContext(c.Request().Context()); tx != nil {
						e := apm.CaptureError(c.Request().Context(), err)
						e.SetTransaction(tx)
						e.Send()
						tx.Result = "panic"
						tx.Context.SetHTTPStatusCode(http.StatusInternalServerError)
					}
					_ = c.JSON(http.StatusInternalServerError, map[string]any{
						"success": false,
						"message": "internal server error",
					})
				}
			}()
			return next(c)
		}
	}
}
