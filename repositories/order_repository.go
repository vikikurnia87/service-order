package repositories

import (
	"context"
	"log/slog"

	"github.com/vikikurnia87/service-order/configs"
	"github.com/vikikurnia87/service-order/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/vikikurnia87/service-utils/dbutil"
	"github.com/vikikurnia87/service-utils/monitoring"
)

// OrderRepository menangani order inti (t_order) + anak (kategori/vendor/assign/komentar)
// dan resolve referensi (priority/status/category UUID → id). Orkestrasi transaksi di service.
type OrderRepository interface {
	FindAll(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Order, int, error)
	FindByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Order, error)

	ResolvePriorityID(ctx context.Context, priorityUUID uuid.UUID) (int64, error)
	ResolveStatusID(ctx context.Context, statusUUID uuid.UUID) (int64, error)
	ResolveCategoryIDs(ctx context.Context, companyUUID uuid.UUID, categoryUUIDs []uuid.UUID) (map[uuid.UUID]int64, error)

	Create(ctx context.Context, m *models.Order) (int64, error)
	InsertCategories(ctx context.Context, rows []models.OrderCategory) error
	InsertVendors(ctx context.Context, rows []models.OrderVendor) error
	InsertAssignees(ctx context.Context, rows []models.OrderAssign) error
	InsertComments(ctx context.Context, rows []models.OrderComment) error
}

type orderRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewOrderRepository(db bun.IDB, logger *slog.Logger) OrderRepository {
	return &orderRepository{db: db, logger: logger}
}

func (r *orderRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *orderRepository) recordErr(ctx context.Context, span *monitoring.APMSpan, msg, table string, params map[string]any, err error) {
	monitoring.RecordRepositoryError(ctx, span, err, table, monitoring.OpDatabaseOperation, params)
	monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), msg, table, params, err)
}

func (r *orderRepository) FindAll(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Order, int, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.Order
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Relation("OrderStatus").Relation("OrderPriority").
		Where("o.company_uuid = ?", companyUUID).
		Order("o.id DESC").Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("(o.name ILIKE ? OR o.order_number ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		params := map[string]any{"company_uuid": companyUUID, "limit": limit, "offset": offset}
		r.recordErr(ctx, span, "find all orders", "t_order", params, err)
		return nil, 0, err
	}
	return rows, count, nil
}

func (r *orderRepository) FindByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Order, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.Order)
	query := r.getDB(ctx).NewSelect().Model(m).
		Relation("OrderStatus").Relation("OrderPriority").
		Relation("Categories.Category").Relation("Vendors").
		Relation("Assignees").Relation("Comments").
		Where("o.id = ?", id).Where("o.company_uuid = ?", companyUUID)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		r.recordErr(ctx, span, "find order by id", "t_order", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}

func (r *orderRepository) ResolvePriorityID(ctx context.Context, priorityUUID uuid.UUID) (int64, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.ResolvePriorityID", monitoring.SpanTypeDB)
	defer span.End()

	var id int64
	err := r.getDB(ctx).NewSelect().Model((*models.OrderPriority)(nil)).
		Column("id").Where("op.order_priority_uuid = ?", priorityUUID).Scan(ctx, &id)
	if err != nil {
		if err.Error() == "bun: no rows in result set" {
			return 0, nil
		}
		r.recordErr(ctx, span, "resolve priority id", "t_order_priority", map[string]any{"uuid": priorityUUID}, err)
		return 0, err
	}
	return id, nil
}

func (r *orderRepository) ResolveStatusID(ctx context.Context, statusUUID uuid.UUID) (int64, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.ResolveStatusID", monitoring.SpanTypeDB)
	defer span.End()

	var id int64
	err := r.getDB(ctx).NewSelect().Model((*models.OrderStatus)(nil)).
		Column("id").Where("os.order_status_uuid = ?", statusUUID).Scan(ctx, &id)
	if err != nil {
		if err.Error() == "bun: no rows in result set" {
			return 0, nil
		}
		r.recordErr(ctx, span, "resolve status id", "t_order_status", map[string]any{"uuid": statusUUID}, err)
		return 0, err
	}
	return id, nil
}

// ResolveCategoryIDs memetakan category_uuid → id untuk kategori milik company.
// UUID yang tak ditemukan tidak akan muncul di map (service yang menentukan error).
func (r *orderRepository) ResolveCategoryIDs(ctx context.Context, companyUUID uuid.UUID, categoryUUIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	out := make(map[uuid.UUID]int64, len(categoryUUIDs))
	if len(categoryUUIDs) == 0 {
		return out, nil
	}
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.ResolveCategoryIDs", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.Category
	err := r.getDB(ctx).NewSelect().Model(&rows).
		Column("id", "category_uuid").
		Where("cat.company_uuid = ?", companyUUID).
		Where("cat.category_uuid IN (?)", bun.In(categoryUUIDs)).
		Scan(ctx)
	if err != nil {
		r.recordErr(ctx, span, "resolve category ids", "t_category", map[string]any{"count": len(categoryUUIDs)}, err)
		return nil, err
	}
	for i := range rows {
		out[rows[i].CategoryUUID] = rows[i].ID
	}
	return out, nil
}

func (r *orderRepository) Create(ctx context.Context, m *models.Order) (int64, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.Create", monitoring.SpanTypeDB)
	defer span.End()

	if _, err := r.getDB(ctx).NewInsert().Model(m).Returning("id").Exec(ctx); err != nil {
		r.recordErr(ctx, span, "insert order", "t_order", map[string]any{"name": m.Name}, err)
		return 0, err
	}
	return m.ID, nil
}

func (r *orderRepository) InsertCategories(ctx context.Context, rows []models.OrderCategory) error {
	if len(rows) == 0 {
		return nil
	}
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.InsertCategories", monitoring.SpanTypeDB)
	defer span.End()

	if _, err := r.getDB(ctx).NewInsert().Model(&rows).Exec(ctx); err != nil {
		r.recordErr(ctx, span, "insert order categories", "t_order_category", map[string]any{"count": len(rows)}, err)
		return err
	}
	return nil
}

func (r *orderRepository) InsertVendors(ctx context.Context, rows []models.OrderVendor) error {
	if len(rows) == 0 {
		return nil
	}
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.InsertVendors", monitoring.SpanTypeDB)
	defer span.End()

	if _, err := r.getDB(ctx).NewInsert().Model(&rows).Exec(ctx); err != nil {
		r.recordErr(ctx, span, "insert order vendors", "t_order_vendor", map[string]any{"count": len(rows)}, err)
		return err
	}
	return nil
}

func (r *orderRepository) InsertAssignees(ctx context.Context, rows []models.OrderAssign) error {
	if len(rows) == 0 {
		return nil
	}
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.InsertAssignees", monitoring.SpanTypeDB)
	defer span.End()

	if _, err := r.getDB(ctx).NewInsert().Model(&rows).Exec(ctx); err != nil {
		r.recordErr(ctx, span, "insert order assignees", "t_order_assign", map[string]any{"count": len(rows)}, err)
		return err
	}
	return nil
}

func (r *orderRepository) InsertComments(ctx context.Context, rows []models.OrderComment) error {
	if len(rows) == 0 {
		return nil
	}
	span, ctx := monitoring.StartSpan(ctx, "repository.orderRepository.InsertComments", monitoring.SpanTypeDB)
	defer span.End()

	if _, err := r.getDB(ctx).NewInsert().Model(&rows).Exec(ctx); err != nil {
		r.recordErr(ctx, span, "insert order comments", "t_order_comment", map[string]any{"count": len(rows)}, err)
		return err
	}
	return nil
}
