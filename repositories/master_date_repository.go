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

// MasterDateRepository menangani tabel kalender (t_master_date). Read-only dari handler;
// baris di-generate programatik per tahun.
type MasterDateRepository interface {
	FindAll(ctx context.Context, year int, limit, offset int) ([]models.MasterDate, int, error)
	FindByID(ctx context.Context, id int64) (*models.MasterDate, error)
}

type masterDateRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewMasterDateRepository(db bun.IDB, logger *slog.Logger) MasterDateRepository {
	return &masterDateRepository{db: db, logger: logger}
}

func (r *masterDateRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *masterDateRepository) FindAll(ctx context.Context, year int, limit, offset int) ([]models.MasterDate, int, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.masterDateRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.MasterDate
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Where("md.status = 1").
		Order("md.full_date ASC").
		Limit(limit).Offset(offset)
	if year > 0 {
		query = query.Where("md.year = ?", year)
	}
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		monitoring.RecordRepositoryError(ctx, span, err, "t_master_date", monitoring.OpDatabaseOperation, map[string]any{"year": year, "limit": limit, "offset": offset})
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find all master dates", "t_master_date", map[string]any{"year": year}, err)
		return nil, 0, err
	}
	return rows, count, nil
}

func (r *masterDateRepository) FindByID(ctx context.Context, id int64) (*models.MasterDate, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.masterDateRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.MasterDate)
	query := r.getDB(ctx).NewSelect().Model(m).Where("md.id = ?", id)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		monitoring.RecordRepositoryError(ctx, span, err, "t_master_date", monitoring.OpDatabaseOperation, map[string]any{"id": id})
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find master date by id", "t_master_date", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}
