package repositories

import (
	"context"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/vikikurnia87/service-order/configs"
	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-utils/dbutil"
	"github.com/vikikurnia87/service-utils/monitoring"
)

// ScheduleRepository menangani referensi pola jadwal berulang (t_schedule). Read-only dari handler.
type ScheduleRepository interface {
	FindAll(ctx context.Context) ([]models.Schedule, error)
	FindByID(ctx context.Context, id int64) (*models.Schedule, error)
}

type scheduleRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewScheduleRepository(db bun.IDB, logger *slog.Logger) ScheduleRepository {
	return &scheduleRepository{db: db, logger: logger}
}

func (r *scheduleRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *scheduleRepository) FindAll(ctx context.Context) ([]models.Schedule, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.scheduleRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.Schedule
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Where("sc.status = 1").
		Order("sc.id ASC")
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		monitoring.RecordRepositoryError(ctx, span, err, "t_schedule", monitoring.OpDatabaseOperation, nil)
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find all schedules", "t_schedule", nil, err)
		return nil, err
	}
	return rows, nil
}

func (r *scheduleRepository) FindByID(ctx context.Context, id int64) (*models.Schedule, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.scheduleRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.Schedule)
	query := r.getDB(ctx).NewSelect().Model(m).Where("sc.id = ?", id)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		monitoring.RecordRepositoryError(ctx, span, err, "t_schedule", monitoring.OpDatabaseOperation, map[string]any{"id": id})
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find schedule by id", "t_schedule", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}
