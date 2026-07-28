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

// OrderPriorityHandler menangani endpoint referensi prioritas order (read-only).
type OrderPriorityHandler struct {
	svc    services.OrderPriorityService
	logger *slog.Logger
}

func NewOrderPriorityHandler(svc services.OrderPriorityService, logger *slog.Logger) *OrderPriorityHandler {
	return &OrderPriorityHandler{svc: svc, logger: logger}
}

// List GET /api/v1/order-priorities — daftar semua prioritas order.
func (h *OrderPriorityHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.OrderPriorityHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	rows, err := h.svc.List(ctx)
	if err != nil {
		return h.fail(c, span, ctx, "list order priorities", "failed to list order priorities", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapOrderPriority(&rows[i]))
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", data, nil)
}

// Get GET /api/v1/order-priorities/:id — detail prioritas order.
func (h *OrderPriorityHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.OrderPriorityHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrOrderPriorityNotFound) {
			return response.JSONBadRequest(c, "order priority not found")
		}
		return h.fail(c, span, ctx, "get order priority", "failed to get order priority", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapOrderPriority(m), nil)
}

func (h *OrderPriorityHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
