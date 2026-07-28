package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-order/mappers"
	"github.com/vikikurnia87/service-order/services"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
	"github.com/vikikurnia87/service-utils/pagination"
	"github.com/vikikurnia87/service-utils/response"
	"github.com/vikikurnia87/service-utils/types"
)

// MasterDateHandler menangani endpoint kalender (read-only).
type MasterDateHandler struct {
	svc    services.MasterDateService
	logger *slog.Logger
}

func NewMasterDateHandler(svc services.MasterDateService, logger *slog.Logger) *MasterDateHandler {
	return &MasterDateHandler{svc: svc, logger: logger}
}

// List GET /api/v1/master-dates — daftar kalender dengan filter opsional ?year=YYYY.
func (h *MasterDateHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.MasterDateHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	p := pagination.GetPaginationParams(c, 31)
	year := 0
	if y := c.QueryParam("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}

	rows, total, err := h.svc.List(ctx, year, p.Limit, p.Offset)
	if err != nil {
		return h.fail(c, span, ctx, "list master dates", "failed to list master dates", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapMasterDate(&rows[i]))
	}
	meta := pagination.BuildPaginationMeta(len(data), total, p.Page, p.Limit)
	return response.JSONSuccess(c, http.StatusOK, "ok", data, &types.SuccessOption{Meta: meta})
}

// Get GET /api/v1/master-dates/:id — detail baris kalender.
func (h *MasterDateHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.MasterDateHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrMasterDateNotFound) {
			return response.JSONBadRequest(c, "master date not found")
		}
		return h.fail(c, span, ctx, "get master date", "failed to get master date", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapMasterDate(m), nil)
}

func (h *MasterDateHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
