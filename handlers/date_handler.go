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

// DateHandler menangani endpoint referensi tanggal dalam bulan (read-only).
type DateHandler struct {
	svc    services.DateService
	logger *slog.Logger
}

func NewDateHandler(svc services.DateService, logger *slog.Logger) *DateHandler {
	return &DateHandler{svc: svc, logger: logger}
}

// List GET /api/v1/dates — daftar semua tanggal dalam bulan.
func (h *DateHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.DateHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	rows, err := h.svc.List(ctx)
	if err != nil {
		return h.fail(c, span, ctx, "list dates", "failed to list dates", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapDate(&rows[i]))
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", data, nil)
}

// Get GET /api/v1/dates/:id — detail tanggal dalam bulan.
func (h *DateHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.DateHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrDateNotFound) {
			return response.JSONBadRequest(c, "date not found")
		}
		return h.fail(c, span, ctx, "get date", "failed to get date", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapDate(m), nil)
}

func (h *DateHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
