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

// DateRepository menangani referensi tanggal dalam bulan (t_date). Read-only dari handler.
type DateRepository interface {
	FindAll(ctx context.Context) ([]models.Date, error)
	FindByID(ctx context.Context, id int64) (*models.Date, error)
}

type dateRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewDateRepository(db bun.IDB, logger *slog.Logger) DateRepository {
	return &dateRepository{db: db, logger: logger}
}

func (r *dateRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *dateRepository) FindAll(ctx context.Context) ([]models.Date, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.dateRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.Date
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Where("dt.status = 1").
		Order("dt.id ASC")
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		monitoring.RecordRepositoryError(ctx, span, err, "t_date", monitoring.OpDatabaseOperation, nil)
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find all dates", "t_date", nil, err)
		return nil, err
	}
	return rows, nil
}

func (r *dateRepository) FindByID(ctx context.Context, id int64) (*models.Date, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.dateRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.Date)
	query := r.getDB(ctx).NewSelect().Model(m).Where("dt.id = ?", id)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		monitoring.RecordRepositoryError(ctx, span, err, "t_date", monitoring.OpDatabaseOperation, map[string]any{"id": id})
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find date by id", "t_date", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}
