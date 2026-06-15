package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/vikikurnia87/service-order/auth"
	"github.com/vikikurnia87/service-order/mappers"
	"github.com/vikikurnia87/service-order/services"
	"github.com/vikikurnia87/service-order/structs"
	"github.com/vikikurnia87/service-order/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-utils/monitoring"
	"github.com/vikikurnia87/service-utils/pagination"
	"github.com/vikikurnia87/service-utils/response"
	"github.com/vikikurnia87/service-utils/types"
)

type OrderHandler struct {
	svc       services.OrderService
	validator *validator.Validate
	logger    *slog.Logger
}

func NewOrderHandler(svc services.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{svc: svc, validator: utils.NewValidator(), logger: logger}
}

func (h *OrderHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.OrderHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	p := pagination.GetPaginationParams(c, 10)
	rows, total, err := h.svc.List(ctx, auth.CompanyUUID(c), c.QueryParam("search"), p.Limit, p.Offset)
	if err != nil {
		return h.fail(c, span, ctx, "list orders", "failed to list orders", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapOrderList(&rows[i]))
	}
	meta := pagination.BuildPaginationMeta(len(data), total, p.Page, p.Limit)
	return response.JSONSuccess(c, http.StatusOK, "ok", data, &types.SuccessOption{Meta: meta})
}

func (h *OrderHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.OrderHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, auth.CompanyUUID(c), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return response.JSONBadRequest(c, "order not found")
		}
		return h.fail(c, span, ctx, "get order", "failed to get order", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapOrder(m), nil)
}

func (h *OrderHandler) Create(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.OrderHandler.Create", monitoring.SpanTypeHTTP)
	defer span.End()

	var req structs.OrderCreateRequest
	if err := c.Bind(&req); err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidBody)
	}
	if err := h.validator.Struct(req); err != nil {
		return response.JSONValidationError(c, utils.MsgValidation, utils.TranslateErrorMessage(err))
	}
	m, err := h.svc.Create(ctx, auth.CompanyUUID(c), req, auth.UserUUIDPtr(c))
	if err != nil {
		return h.businessErr(c, span, ctx, "create order", "failed to create order", err)
	}
	return response.JSONSuccess(c, http.StatusCreated, "order created", mappers.MapOrder(m), nil)
}

// businessErr memetakan error bisnis ke 4xx (selain itu 500).
func (h *OrderHandler) businessErr(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	switch {
	case errors.Is(err, services.ErrPriorityNotFound):
		return response.JSONValidationError(c, "priority not found", nil)
	case errors.Is(err, services.ErrStatusNotFound):
		return response.JSONValidationError(c, "status not found", nil)
	case errors.Is(err, services.ErrCategoryNotFound):
		return response.JSONValidationError(c, "category not found", nil)
	case errors.Is(err, services.ErrNotFound):
		return response.JSONBadRequest(c, "order not found")
	default:
		return h.fail(c, span, ctx, action, clientMsg, err)
	}
}

func (h *OrderHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
