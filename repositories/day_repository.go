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

// DayRepository menangani referensi hari (t_day). Read-only dari handler.
type DayRepository interface {
	FindAll(ctx context.Context) ([]models.Day, error)
	FindByID(ctx context.Context, id int64) (*models.Day, error)
}

type dayRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewDayRepository(db bun.IDB, logger *slog.Logger) DayRepository {
	return &dayRepository{db: db, logger: logger}
}

func (r *dayRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *dayRepository) FindAll(ctx context.Context) ([]models.Day, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.dayRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.Day
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Where("dy.status = 1").
		Order("dy.id ASC")
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		monitoring.RecordRepositoryError(ctx, span, err, "t_day", monitoring.OpDatabaseOperation, nil)
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find all days", "t_day", nil, err)
		return nil, err
	}
	return rows, nil
}

func (r *dayRepository) FindByID(ctx context.Context, id int64) (*models.Day, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.dayRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.Day)
	query := r.getDB(ctx).NewSelect().Model(m).Where("dy.id = ?", id)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		monitoring.RecordRepositoryError(ctx, span, err, "t_day", monitoring.OpDatabaseOperation, map[string]any{"id": id})
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find day by id", "t_day", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}
