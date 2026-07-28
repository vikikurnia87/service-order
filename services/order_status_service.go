package services

import (
	"context"
	"log/slog"

	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
)

type OrderStatusService interface {
	List(ctx context.Context) ([]models.OrderStatus, error)
	GetByID(ctx context.Context, id int64) (*models.OrderStatus, error)
}

type orderStatusService struct {
	repo   repositories.OrderStatusRepository
	logger *slog.Logger
}

func NewOrderStatusService(repo repositories.OrderStatusRepository, logger *slog.Logger) OrderStatusService {
	return &orderStatusService{repo: repo, logger: logger}
}

func (s *orderStatusService) List(ctx context.Context) ([]models.OrderStatus, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.orderStatusService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, err := s.repo.FindAll(ctx)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "list order statuses failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()))
		return nil, err
	}
	return rows, nil
}

func (s *orderStatusService) GetByID(ctx context.Context, id int64) (*models.OrderStatus, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.orderStatusService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "get order status by id failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	if m == nil {
		return nil, utils.ErrOrderStatusNotFound
	}
	return m, nil
}
