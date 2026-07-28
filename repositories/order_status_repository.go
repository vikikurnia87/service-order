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

// OrderStatusRepository menangani referensi status order (t_order_status). Read-only dari handler;
// data di-seed sekali lewat migrasi.
type OrderStatusRepository interface {
	FindAll(ctx context.Context) ([]models.OrderStatus, error)
	FindByID(ctx context.Context, id int64) (*models.OrderStatus, error)
}

type orderStatusRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewOrderStatusRepository(db bun.IDB, logger *slog.Logger) OrderStatusRepository {
	return &orderStatusRepository{db: db, logger: logger}
}

func (r *orderStatusRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *orderStatusRepository) FindAll(ctx context.Context) ([]models.OrderStatus, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderStatusRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.OrderStatus
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Where("os.status = 1").
		Order("os.id ASC")
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		monitoring.RecordRepositoryError(ctx, span, err, "t_order_status", monitoring.OpDatabaseOperation, nil)
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find all order statuses", "t_order_status", nil, err)
		return nil, err
	}
	return rows, nil
}

func (r *orderStatusRepository) FindByID(ctx context.Context, id int64) (*models.OrderStatus, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderStatusRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.OrderStatus)
	query := r.getDB(ctx).NewSelect().Model(m).Where("os.id = ?", id)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		monitoring.RecordRepositoryError(ctx, span, err, "t_order_status", monitoring.OpDatabaseOperation, map[string]any{"id": id})
		monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), "find order status by id", "t_order_status", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}
