package services

import (
	"context"
	"log/slog"

	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/monitoring"
)

type MasterDateService interface {
	List(ctx context.Context, year int, limit, offset int) ([]models.MasterDate, int, error)
	GetByID(ctx context.Context, id int64) (*models.MasterDate, error)
}

type masterDateService struct {
	repo   repositories.MasterDateRepository
	logger *slog.Logger
}

func NewMasterDateService(repo repositories.MasterDateRepository, logger *slog.Logger) MasterDateService {
	return &masterDateService{repo: repo, logger: logger}
}

func (s *masterDateService) List(ctx context.Context, year int, limit, offset int) ([]models.MasterDate, int, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.masterDateService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, count, err := s.repo.FindAll(ctx, year, limit, offset)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "list master dates failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()))
		return nil, 0, err
	}
	return rows, count, nil
}

func (s *masterDateService) GetByID(ctx context.Context, id int64) (*models.MasterDate, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.masterDateService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
		monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, "get master date by id failed",
			slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
			slog.String(monitoring.LabelErrorMessage, err.Error()),
			slog.Int64("id", id))
		return nil, err
	}
	if m == nil {
		return nil, utils.ErrMasterDateNotFound
	}
	return m, nil
}
