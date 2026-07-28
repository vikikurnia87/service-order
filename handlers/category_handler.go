package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/vikikurnia87/service-order/auth"
	"github.com/vikikurnia87/service-order/mappers"
	"github.com/vikikurnia87/service-order/services"
	"github.com/vikikurnia87/service-order/structs"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
	"github.com/vikikurnia87/service-utils/pagination"
	"github.com/vikikurnia87/service-utils/response"
	"github.com/vikikurnia87/service-utils/types"
)

// CategoryHandler menangani endpoint master kategori order (scoped per company).
type CategoryHandler struct {
	svc       services.CategoryService
	validator *validator.Validate
	logger    *slog.Logger
}

func NewCategoryHandler(svc services.CategoryService, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{svc: svc, validator: utils.NewValidator(), logger: logger}
}

// List GET /api/v1/categories — daftar kategori milik company aktif.
func (h *CategoryHandler) List(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.CategoryHandler.List", monitoring.SpanTypeHTTP)
	defer span.End()

	p := pagination.GetPaginationParams(c, 10)
	search := c.QueryParam("search")

	rows, total, err := h.svc.List(ctx, auth.CompanyUUID(c), search, p.Limit, p.Offset)
	if err != nil {
		return h.fail(c, span, ctx, "list categories", "failed to list categories", err)
	}
	data := make([]map[string]any, 0, len(rows))
	for i := range rows {
		data = append(data, mappers.MapCategory(&rows[i]))
	}
	meta := pagination.BuildPaginationMeta(len(data), total, p.Page, p.Limit)
	return response.JSONSuccess(c, http.StatusOK, "ok", data, &types.SuccessOption{Meta: meta})
}

// Get GET /api/v1/categories/:id — detail kategori.
func (h *CategoryHandler) Get(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.CategoryHandler.Get", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	m, err := h.svc.GetByID(ctx, auth.CompanyUUID(c), id)
	if err != nil {
		if errors.Is(err, utils.ErrCategoryEntityNotFound) {
			return response.JSONBadRequest(c, "category not found")
		}
		return h.fail(c, span, ctx, "get category", "failed to get category", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "ok", mappers.MapCategory(m), nil)
}

// Create POST /api/v1/categories — buat kategori baru dalam company aktif.
func (h *CategoryHandler) Create(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.CategoryHandler.Create", monitoring.SpanTypeHTTP)
	defer span.End()

	var req structs.CategoryCreateRequest
	if err := c.Bind(&req); err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidBody)
	}
	if err := h.validator.Struct(req); err != nil {
		return response.JSONValidationError(c, utils.MsgValidation, utils.TranslateErrorMessage(err))
	}
	m, err := h.svc.Create(ctx, auth.CompanyUUID(c), req, auth.UserUUIDPtr(c))
	if err != nil {
		if errors.Is(err, utils.ErrCategoryNameExists) {
			return response.JSONValidationError(c, "category name already exists", nil)
		}
		return h.fail(c, span, ctx, "create category", "failed to create category", err)
	}
	return response.JSONSuccess(c, http.StatusCreated, "category created", mappers.MapCategory(m), nil)
}

// Update PUT /api/v1/categories/:id — ubah nama/deskripsi kategori.
func (h *CategoryHandler) Update(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.CategoryHandler.Update", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	var req structs.CategoryUpdateRequest
	if err := c.Bind(&req); err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidBody)
	}
	if err := h.validator.Struct(req); err != nil {
		return response.JSONValidationError(c, utils.MsgValidation, utils.TranslateErrorMessage(err))
	}
	m, err := h.svc.Update(ctx, auth.CompanyUUID(c), id, req, auth.UserUUIDPtr(c))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrCategoryEntityNotFound):
			return response.JSONBadRequest(c, "category not found")
		case errors.Is(err, utils.ErrCategoryNameExists):
			return response.JSONValidationError(c, "category name already exists", nil)
		}
		return h.fail(c, span, ctx, "update category", "failed to update category", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "category updated", mappers.MapCategory(m), nil)
}

// Delete DELETE /api/v1/categories/:id — soft-delete kategori.
func (h *CategoryHandler) Delete(c *echo.Context) error {
	span, ctx := monitoring.StartSpan(c.Request().Context(), "handler.CategoryHandler.Delete", monitoring.SpanTypeHTTP)
	defer span.End()

	id, err := utils.ParamID(c)
	if err != nil {
		return response.JSONError(c, http.StatusBadRequest, utils.MsgInvalidID)
	}
	if err := h.svc.Delete(ctx, auth.CompanyUUID(c), id, auth.UserUUIDPtr(c)); err != nil {
		if errors.Is(err, utils.ErrCategoryEntityNotFound) {
			return response.JSONBadRequest(c, "category not found")
		}
		return h.fail(c, span, ctx, "delete category", "failed to delete category", err)
	}
	return response.JSONSuccess(c, http.StatusOK, "category deleted", nil, nil)
}

func (h *CategoryHandler) fail(c *echo.Context, span *monitoring.APMSpan, ctx context.Context, action, clientMsg string, err error) error {
	monitoring.RecordSpanError(span, monitoring.LogErrorHandler, monitoring.SpanLayerHandler, monitoring.OpServiceCall, err)
	monitoring.ContextLogger(ctx, h.logger).ErrorContext(ctx, action+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerHandler),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
	return response.JSONError(c, http.StatusInternalServerError, clientMsg)
}
