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

// DayHandler menangani endpoint referensi hari (read-only).
type DayHandler struct {
	svc    services.DayService
	logger *slog.Logger
}

func NewDayHandler(svc services.DayService, logger *slog.Logger) *DayHandler {
	return &DayHandler{svc: svc, logger: logger}
}

// List GET /api/v1/days — daftar semua hari.
func (h *DayHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.DayHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	rows, err := h.svc.List(ctx)
	if err != nil {
		return h.fail(c, span, ctx, "list days", "failed to list days", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapDay(&rows[i]))
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", data, nil)
}

// Get GET /api/v1/days/:id — detail hari.
func (h *DayHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.DayHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrDayNotFound) {
			return response.JSONBadRequest(c, "day not found")
		}
		return h.fail(c, span, ctx, "get day", "failed to get day", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapDay(m), nil)
}

func (h *DayHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
