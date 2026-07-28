package services

import (
	"context"
	"log/slog"

	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
)

type ScheduleService interface {
	List(ctx context.Context) ([]models.Schedule, error)
	GetByID(ctx context.Context, id int64) (*models.Schedule, error)
}

type scheduleService struct {
	repo   repositories.ScheduleRepository
	logger *slog.Logger
}

func NewScheduleService(repo repositories.ScheduleRepository, logger *slog.Logger) ScheduleService {
	return &scheduleService{repo: repo, logger: logger}
}

func (s *scheduleService) List(ctx context.Context) ([]models.Schedule, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.scheduleService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, err := s.repo.FindAll(ctx)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "list schedules failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()))
		return nil, err
	}
	return rows, nil
}

func (s *scheduleService) GetByID(ctx context.Context, id int64) (*models.Schedule, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.scheduleService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "get schedule by id failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	if m == nil {
		return nil, utils.ErrScheduleNotFound
	}
	return m, nil
}
