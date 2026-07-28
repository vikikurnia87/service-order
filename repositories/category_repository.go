package repositories

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/vikikurnia87/service-order/configs"
	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-utils/dbutil"
	"github.com/vikikurnia87/service-utils/monitoring"
)

// CategoryRepository menangani master kategori order (t_category) ter-scope per company_uuid.
type CategoryRepository interface {
	FindAll(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Category, int, error)
	FindByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Category, error)
	ExistsByID(ctx context.Context, companyUUID uuid.UUID, id int64) (bool, error)
	ExistsByName(ctx context.Context, companyUUID uuid.UUID, name string, excludeID int64) (bool, error)
	Create(ctx context.Context, m *models.Category) (int64, error)
	Update(ctx context.Context, m *models.Category) error
	SoftDelete(ctx context.Context, companyUUID uuid.UUID, id int64, actor *uuid.UUID) error
}

type categoryRepository struct {
	db     bun.IDB
	logger *slog.Logger
}

func NewCategoryRepository(db bun.IDB, logger *slog.Logger) CategoryRepository {
	return &categoryRepository{db: db, logger: logger}
}

func (r *categoryRepository) getDB(ctx context.Context) bun.IDB {
	if tx, ok := dbutil.FromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *categoryRepository) recordErr(ctx context.Context, span *monitoring.APMSpan, msg string, params map[string]any, err error) {
	monitoring.RecordRepositoryError(ctx, span, err, "t_category", monitoring.OpDatabaseOperation, params)
	monitoring.LogRepositoryError(ctx, monitoring.ContextLogger(ctx, r.logger), msg, "t_category", params, err)
}

func (r *categoryRepository) FindAll(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Category, int, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.categoryRepository.FindAll", monitoring.SpanTypeDB)
	defer span.End()

	var rows []models.Category
	query := r.getDB(ctx).NewSelect().Model(&rows).
		Where("cat.company_uuid = ?", companyUUID).
		Where("cat.status = 1").
		Order("cat.id DESC").Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("cat.name ILIKE ?", "%"+search+"%")
	}
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		r.recordErr(ctx, span, "find all categories", map[string]any{"company_uuid": companyUUID, "limit": limit, "offset": offset}, err)
		return nil, 0, err
	}
	return rows, count, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Category, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.categoryRepository.FindByID", monitoring.SpanTypeDB)
	defer span.End()

	m := new(models.Category)
	query := r.getDB(ctx).NewSelect().Model(m).
		Where("cat.id = ?", id).
		Where("cat.company_uuid = ?", companyUUID)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "bun: no rows in result set" {
			return nil, nil
		}
		r.recordErr(ctx, span, "find category by id", map[string]any{"id": id}, err)
		return nil, err
	}
	return m, nil
}

func (r *categoryRepository) ExistsByID(ctx context.Context, companyUUID uuid.UUID, id int64) (bool, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.categoryRepository.ExistsByID", monitoring.SpanTypeDB)
	defer span.End()

	query := r.getDB(ctx).NewSelect().Model((*models.Category)(nil)).
		Where("cat.id = ?", id).
		Where("cat.company_uuid = ?", companyUUID).
		Where("cat.status = 1")
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	count, err := query.Count(ctx)
	if err != nil {
		r.recordErr(ctx, span, "check category exists by id", map[string]any{"id": id}, err)
		return false, err
	}
	return count > 0, nil
}

// ExistsByName mengecek duplikasi nama dalam company (excludeID 0 = tidak ada pengecualian).
func (r *categoryRepository) ExistsByName(ctx context.Context, companyUUID uuid.UUID, name string, excludeID int64) (bool, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.categoryRepository.ExistsByName", monitoring.SpanTypeDB)
	defer span.End()

	query := r.getDB(ctx).NewSelect().Model((*models.Category)(nil)).
		Where("LOWER(cat.name) = LOWER(?)", name).
		Where("cat.company_uuid = ?", companyUUID).
		Where("cat.status = 1")
	if excludeID > 0 {
		query = query.Where("cat.id != ?", excludeID)
	}
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	count, err := query.Count(ctx)
	if err != nil {
		r.recordErr(ctx, span, "check category exists by name", map[string]any{"name": name}, err)
		return false, err
	}
	return count > 0, nil
}

func (r *categoryRepository) Create(ctx context.Context, m *models.Category) (int64, error) {
	span, ctx := monitoring.StartSpan(ctx, "repository.categoryRepository.Create", monitoring.SpanTypeDB)
	defer span.End()

	if _, err := r.getDB(ctx).NewInsert().Model(m).Returning("id").Exec(ctx); err != nil {
		r.recordErr(ctx, span, "insert category", map[string]any{"name": m.Name}, err)
		return 0, err
	}
	r.logger.DebugContext(ctx, "category inserted", slog.Int64("category_id", m.ID))
	return m.ID, nil
}

func (r *categoryRepository) Update(ctx context.Context, m *models.Category) error {
	span, ctx := monitoring.StartSpan(ctx, "repository.categoryRepository.Update", monitoring.SpanTypeDB)
	defer span.End()

	if _, err := r.getDB(ctx).NewUpdate().Model(m).WherePK().Exec(ctx); err != nil {
		r.recordErr(ctx, span, "update category", map[string]any{"id": m.ID}, err)
		return err
	}
	return nil
}

// SoftDelete meng-nonaktifkan kategori (status=0).
func (r *categoryRepository) SoftDelete(ctx context.Context, companyUUID uuid.UUID, id int64, actor *uuid.UUID) error {
	span, ctx := monitoring.StartSpan(ctx, "repository.categoryRepository.SoftDelete", monitoring.SpanTypeDB)
	defer span.End()

	now := time.Now()
	query := r.getDB(ctx).NewUpdate().Model((*models.Category)(nil)).
		Set("status = 0").
		Set("updated_at = ?", now).
		Set("updated_by = ?", actor).
		Where("id = ?", id).
		Where("company_uuid = ?", companyUUID)
	monitoring.SpanSetSQLString(span, query, configs.DatabaseName, configs.DatabaseUser, monitoring.SpanSubtypePostgreSQL)

	if _, err := query.Exec(ctx); err != nil {
		r.recordErr(ctx, span, "soft delete category", map[string]any{"id": id}, err)
		return err
	}
	return nil
}
