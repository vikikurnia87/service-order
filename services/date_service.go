package services

import (
	"context"
	"log/slog"

	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
)

type DateService interface {
	List(ctx context.Context) ([]models.Date, error)
	GetByID(ctx context.Context, id int64) (*models.Date, error)
}

type dateService struct {
	repo   repositories.DateRepository
	logger *slog.Logger
}

func NewDateService(repo repositories.DateRepository, logger *slog.Logger) DateService {
	return &dateService{repo: repo, logger: logger}
}

func (s *dateService) List(ctx context.Context) ([]models.Date, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.dateService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, err := s.repo.FindAll(ctx)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "list dates failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()))
		return nil, err
	}
	return rows, nil
}

func (s *dateService) GetByID(ctx context.Context, id int64) (*models.Date, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.dateService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "get date by id failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	if m == nil {
		return nil, utils.ErrDateNotFound
	}
	return m, nil
}
