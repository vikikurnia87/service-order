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

// ScheduleHandler menangani endpoint referensi pola jadwal (read-only).
type ScheduleHandler struct {
	svc    services.ScheduleService
	logger *slog.Logger
}

func NewScheduleHandler(svc services.ScheduleService, logger *slog.Logger) *ScheduleHandler {
	return &ScheduleHandler{svc: svc, logger: logger}
}

// List GET /api/v1/schedules — daftar semua pola jadwal.
func (h *ScheduleHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.ScheduleHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	rows, err := h.svc.List(ctx)
	if err != nil {
		return h.fail(c, span, ctx, "list schedules", "failed to list schedules", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapSchedule(&rows[i]))
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", data, nil)
}

// Get GET /api/v1/schedules/:id — detail pola jadwal.
func (h *ScheduleHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.ScheduleHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrScheduleNotFound) {
			return response.JSONBadRequest(c, "schedule not found")
		}
		return h.fail(c, span, ctx, "get schedule", "failed to get schedule", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapSchedule(m), nil)
}

func (h *ScheduleHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
