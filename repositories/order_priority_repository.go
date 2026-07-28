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

// OrderPriorityRepository menangani referensi prioritas order (t_order_priority). Read-only dari handler;
// data di-seed sekali lewat migrasi.
type OrderPriorityRepository interface {
	FindAll(ctx context.Context) ([]models.OrderPriority, error)
	FindByID(ctx context.Context, id int64) (*models.OrderPriority, error)
}

type orderPriorityRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewOrderPriorityRepository(db bun.IDB, logger *slog.Logger) OrderPriorityRepository {
	return &orderPriorityRepository{db: db, logger: logger}
}

func (r *orderPriorityRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *orderPriorityRepository) FindAll(ctx context.Context) ([]models.OrderPriority, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderPriorityRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.OrderPriority
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Where("op.status = 1").
		Order("op.id ASC")
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		monitoring.RecordRepositoryError(ctx, span, err, "t_order_priority", monitoring.OpDatabaseOperation, nil)
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find all order priorities", "t_order_priority", nil, err)
		return nil, err
	}
	return rows, nil
}

func (r *orderPriorityRepository) FindByID(ctx context.Context, id int64) (*models.OrderPriority, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderPriorityRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.OrderPriority)
	query := r.getDB(ctx).NewSelect().Model(m).Where("op.id = ?", id)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		monitoring.RecordRepositoryError(ctx, span, err, "t_order_priority", monitoring.OpDatabaseOperation, map[string]any{"id": id})
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find order priority by id", "t_order_priority", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}
