package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/structs"
	"github.com/vikikurnia87/service-order/utils"
	"github.com/vikikurnia87/service-utils/dbutil"
	"github.com/vikikurnia87/service-utils/monitoring"
)

type CategoryService interface {
	List(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Category, int, error)
	GetByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Category, error)
	Create(ctx context.Context, companyUUID uuid.UUID, req structs.CategoryCreateRequest, actor *uuid.UUID) (*models.Category, error)
	Update(ctx context.Context, companyUUID uuid.UUID, id int64, req structs.CategoryUpdateRequest, actor *uuid.UUID) (*models.Category, error)
	Delete(ctx context.Context, companyUUID uuid.UUID, id int64, actor *uuid.UUID) error
}

type categoryService struct {
	repo     repositories.CategoryRepository
	txHelper dbutil.TransactionHelper
	logger   *slog.Logger
}

func NewCategoryService(repo repositories.CategoryRepository, txHelper dbutil.TransactionHelper, logger *slog.Logger) CategoryService {
	return &categoryService{repo: repo, txHelper: txHelper, logger: logger}
}

func (s *categoryService) List(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Category, int, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.categoryService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, count, err := s.repo.FindAll(ctx, companyUUID, search, limit, offset)
	if err != nil {
		s.recordErr(ctx, span, "list categories", err)
		return nil, 0, err
	}
	return rows, count, nil
}

func (s *categoryService) GetByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Category, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.categoryService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, companyUUID, id)
	if err != nil {
		s.recordErr(ctx, span, "get category", err)
		return nil, err
	}
	if m == nil {
		return nil, utils.ErrCategoryEntityNotFound
	}
	return m, nil
}

// Create membuat kategori baru. Nama unik per company.
func (s *categoryService) Create(ctx context.Context, companyUUID uuid.UUID, req structs.CategoryCreateRequest, actor *uuid.UUID) (*models.Category, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.categoryService.Create", monitoring.SpanTypeApp)
	defer span.End()

	now := time.Now()
	descPtr := strPtrOrNil(req.Description)
	m := &models.Category{
		CategoryUUID: uuid.New(),
		CompanyUUID:  companyUUID,
		Name:         req.Name,
		Description:  descPtr,
		Status:       1,
		CreatedAt:    &now,
		UpdatedAt:    &now,
		CreatedBy:    actor,
	}

	err := s.txHelper.WithTx(ctx, func(txCtx context.Context) error {
		exists, cerr := s.repo.ExistsByName(txCtx, companyUUID, req.Name, 0)
		if cerr != nil {
			return cerr
		}
		if exists {
			return utils.ErrCategoryNameExists
		}
		_, cerr = s.repo.Create(txCtx, m)
		return cerr
	})
	if err != nil {
		if errors.Is(err, utils.ErrCategoryNameExists) || errors.Is(err, utils.ErrCategoryEntityNotFound) {
			return nil, err
		}
		s.recordErr(ctx, span, "create category", err)
		return nil, err
	}
	return m, nil
}

// Update mengubah nama dan/atau deskripsi kategori.
func (s *categoryService) Update(ctx context.Context, companyUUID uuid.UUID, id int64, req structs.CategoryUpdateRequest, actor *uuid.UUID) (*models.Category, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.categoryService.Update", monitoring.SpanTypeApp)
	defer span.End()

	var updated *models.Category
	now := time.Now()

	err := s.txHelper.WithTx(ctx, func(txCtx context.Context) error {
		current, ferr := s.repo.FindByID(txCtx, companyUUID, id)
		if ferr != nil {
			return ferr
		}
		if current == nil {
			return utils.ErrCategoryEntityNotFound
		}

		// Cek duplikasi nama (kecuali kategori ini sendiri).
		if current.Name != req.Name {
			nameTaken, cerr := s.repo.ExistsByName(txCtx, companyUUID, req.Name, id)
			if cerr != nil {
				return cerr
			}
			if nameTaken {
				return utils.ErrCategoryNameExists
			}
		}

		current.Name = req.Name
		current.Description = strPtrOrNil(req.Description)
		current.UpdatedAt = &now
		current.UpdatedBy = actor
		updated = current
		return s.repo.Update(txCtx, current)
	})
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrCategoryEntityNotFound):
			return nil, err
		case errors.Is(err, utils.ErrCategoryNameExists):
			return nil, err
		}
		s.recordErr(ctx, span, "update category", err)
		return nil, err
	}
	return updated, nil
}

// Delete = soft-delete (status=0).
func (s *categoryService) Delete(ctx context.Context, companyUUID uuid.UUID, id int64, actor *uuid.UUID) error {
	span, ctx := monitoring.StartSpan(ctx, "service.categoryService.Delete", monitoring.SpanTypeApp)
	defer span.End()

	exists, err := s.repo.ExistsByID(ctx, companyUUID, id)
	if err != nil {
		s.recordErr(ctx, span, "check category exists", err)
		return err
	}
	if !exists {
		return utils.ErrCategoryEntityNotFound
	}
	if err := s.repo.SoftDelete(ctx, companyUUID, id, actor); err != nil {
		s.recordErr(ctx, span, "delete category", err)
		return err
	}
	return nil
}

func (s *categoryService) recordErr(ctx context.Context, span *monitoring.APMSpan, msg string, err error) {
	if errors.Is(err, utils.ErrCategoryEntityNotFound) || errors.Is(err, utils.ErrCategoryNameExists) {
		return
	}
	monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
	monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, msg+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
}
