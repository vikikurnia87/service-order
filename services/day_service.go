package services

import (
	"context"
	"log/slog"

	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
)

type DayService interface {
	List(ctx context.Context) ([]models.Day, error)
	GetByID(ctx context.Context, id int64) (*models.Day, error)
}

type dayService struct {
	repo   repositories.DayRepository
	logger *slog.Logger
}

func NewDayService(repo repositories.DayRepository, logger *slog.Logger) DayService {
	return &dayService{repo: repo, logger: logger}
}

func (s *dayService) List(ctx context.Context) ([]models.Day, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.dayService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, err := s.repo.FindAll(ctx)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "list days failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()))
		return nil, err
	}
	return rows, nil
}

func (s *dayService) GetByID(ctx context.Context, id int64) (*models.Day, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.dayService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "get day by id failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	if m == nil {
		return nil, utils.ErrDayNotFound
	}
	return m, nil
}
