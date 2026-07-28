package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-order/mappers"
	"github.com/vikikurnia87/service-order/services"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
	"github.com/vikikurnia87/service-utils/response"
)

// OrderStatusHandler menangani endpoint referensi status order (read-only).
type OrderStatusHandler struct {
	svc    services.OrderStatusService
	logger *slog.Logger
}

func NewOrderStatusHandler(svc services.OrderStatusService, logger *slog.Logger) *OrderStatusHandler {
	return &OrderStatusHandler{svc: svc, logger: logger}
}

// List GET /api/v1/order-statuses — daftar semua status order.
func (h *OrderStatusHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.OrderStatusHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	rows, err := h.svc.List(ctx)
	if err != nil {
		return h.fail(c, span, ctx, "list order statuses", "failed to list order statuses", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapOrderStatus(&rows[i]))
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", data, nil)
}

// Get GET /api/v1/order-statuses/:id — detail status order.
func (h *OrderStatusHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.OrderStatusHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrOrderStatusNotFound) {
			return response.JSONBadRequest(c, "order status not found")
		}
		return h.fail(c, span, ctx, "get order status", "failed to get order status", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapOrderStatus(m), nil)
}

func (h *OrderStatusHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
