package services

import (
	"context"
	"log/slog"

	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
)

type OrderPriorityService interface {
	List(ctx context.Context) ([]models.OrderPriority, error)
	GetByID(ctx context.Context, id int64) (*models.OrderPriority, error)
}

type orderPriorityService struct {
	repo   repositories.OrderPriorityRepository
	logger *slog.Logger
}

func NewOrderPriorityService(repo repositories.OrderPriorityRepository, logger *slog.Logger) OrderPriorityService {
	return &orderPriorityService{repo: repo, logger: logger}
}

func (s *orderPriorityService) List(ctx context.Context) ([]models.OrderPriority, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.orderPriorityService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, err := s.repo.FindAll(ctx)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "list order priorities failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()))
		return nil, err
	}
	return rows, nil
}

func (s *orderPriorityService) GetByID(ctx context.Context, id int64) (*models.OrderPriority, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.orderPriorityService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "get order priority by id failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	if m == nil {
		return nil, utils.ErrOrderPriorityNotFound
	}
	return m, nil
}
