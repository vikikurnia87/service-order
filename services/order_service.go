package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/vikikurnia87/service-order/models"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/structs"

	"github.com/google/uuid"
	"github.com/vikikurnia87/service-utils/dbutil"
	"github.com/vikikurnia87/service-utils/monitoring"
)

var (
	// ErrNotFound = order tidak ditemukan / bukan milik company.
	ErrNotFound = errors.New("order not found")
	// ErrPriorityNotFound = priority_uuid tidak dikenal.
	ErrPriorityNotFound = errors.New("priority not found")
	// ErrStatusNotFound = status_uuid tidak dikenal.
	ErrStatusNotFound = errors.New("status not found")
	// ErrCategoryNotFound = salah satu category_uuid tidak dikenal / bukan milik company.
	ErrCategoryNotFound = errors.New("category not found")
)

type OrderService interface {
	List(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Order, int, error)
	GetByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Order, error)
	Create(ctx context.Context, companyUUID uuid.UUID, req structs.OrderCreateRequest, actor *uuid.UUID) (*models.Order, error)
}

type orderService struct {
	repo     repositories.OrderRepository
	txHelper dbutil.TransactionHelper
	logger   *slog.Logger
}

func NewOrderService(repo repositories.OrderRepository, txHelper dbutil.TransactionHelper, logger *slog.Logger) OrderService {
	return &orderService{repo: repo, txHelper: txHelper, logger: logger}
}

func (s *orderService) List(ctx context.Context, companyUUID uuid.UUID, search string, limit, offset int) ([]models.Order, int, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.orderService.List", monitoring.SpanTypeApp)
	defer span.End()
	rows, count, err := s.repo.FindAll(ctx, companyUUID, search, limit, offset)
	if err != nil {
		s.recordErr(ctx, span, "list orders", err)
		return nil, 0, err
	}
	return rows, count, nil
}

func (s *orderService) GetByID(ctx context.Context, companyUUID uuid.UUID, id int64) (*models.Order, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.orderService.GetByID", monitoring.SpanTypeApp)
	defer span.End()
	m, err := s.repo.FindByID(ctx, companyUUID, id)
	if err != nil {
		s.recordErr(ctx, span, "get order", err)
		return nil, err
	}
	if m == nil {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *orderService) Create(ctx context.Context, companyUUID uuid.UUID, req structs.OrderCreateRequest, actor *uuid.UUID) (*models.Order, error) {
	span, ctx := monitoring.StartSpan(ctx, "service.orderService.Create", monitoring.SpanTypeApp)
	defer span.End()

	// Resolve referensi WAJIB di luar transaksi (read-only).
	priorityID, err := s.repo.ResolvePriorityID(ctx, uuid.MustParse(req.PriorityUUID))
	if err != nil {
		s.recordErr(ctx, span, "resolve priority", err)
		return nil, err
	}
	if priorityID == 0 {
		return nil, ErrPriorityNotFound
	}
	statusID, err := s.repo.ResolveStatusID(ctx, uuid.MustParse(req.StatusUUID))
	if err != nil {
		s.recordErr(ctx, span, "resolve status", err)
		return nil, err
	}
	if statusID == 0 {
		return nil, ErrStatusNotFound
	}

	catUUIDs := parseUUIDs(req.Categories)
	catIDByUUID, err := s.repo.ResolveCategoryIDs(ctx, companyUUID, catUUIDs)
	if err != nil {
		s.recordErr(ctx, span, "resolve categories", err)
		return nil, err
	}
	if len(catIDByUUID) != len(catUUIDs) {
		return nil, ErrCategoryNotFound
	}

	now := time.Now()
	m := &models.Order{
		OrderUUID:       uuid.New(),
		CompanyUUID:     companyUUID,
		OrderNumber:     strPtrOrNil(req.OrderNumber),
		Name:            req.Name,
		Description:     strPtrOrNil(req.Description),
		OrderPriorityID: priorityID,
		OrderStatusID:   statusID,
		AssetUUID:       parseUUIDPtr(req.AssetUUID),
		LocationUUID:    parseUUIDPtr(req.LocationUUID),
		StartDate:       req.StartDate,
		DueDate:         req.DueDate,
		PeriodDate:      derivePeriodDate(req),
		IsOrder:         true,
		IsSystem:        false,
		Status:          1,
		CreatedAt:       &now,
		UpdatedAt:       &now,
		CreatedBy:       actor,
	}

	err = s.txHelper.WithTx(ctx, func(txCtx context.Context) error {
		orderID, err := s.repo.Create(txCtx, m)
		if err != nil {
			return err
		}

		cats := make([]models.OrderCategory, 0, len(catUUIDs))
		for _, cu := range catUUIDs {
			cats = append(cats, models.OrderCategory{
				OrderCategoryUUID: uuid.New(),
				OrderID:           orderID,
				CategoryID:        catIDByUUID[cu],
				Status:            1,
				CreatedAt:         &now,
				UpdatedAt:         &now,
			})
		}
		if err := s.repo.InsertCategories(txCtx, cats); err != nil {
			return err
		}

		vendors := make([]models.OrderVendor, 0, len(req.Vendors))
		for _, vu := range parseUUIDs(req.Vendors) {
			vendors = append(vendors, models.OrderVendor{
				OrderVendorUUID: uuid.New(),
				OrderID:         orderID,
				VendorUUID:      vu,
				Status:          1,
				CreatedAt:       &now,
				UpdatedAt:       &now,
			})
		}
		if err := s.repo.InsertVendors(txCtx, vendors); err != nil {
			return err
		}

		assignees := make([]models.OrderAssign, 0, len(req.Assignees))
		for _, a := range req.Assignees {
			team := parseUUIDPtr(a.TeamUUID)
			user := parseUUIDPtr(a.UserUUID)
			if team == nil && user == nil {
				continue // baris kosong diabaikan
			}
			assignees = append(assignees, models.OrderAssign{
				OrderAssignUUID: uuid.New(),
				OrderID:         orderID,
				TeamUUID:        team,
				UserUUID:        user,
				Status:          1,
				CreatedAt:       &now,
				UpdatedAt:       &now,
			})
		}
		if err := s.repo.InsertAssignees(txCtx, assignees); err != nil {
			return err
		}

		comments := make([]models.OrderComment, 0, len(req.Comments))
		for _, txt := range req.Comments {
			if txt == "" {
				continue
			}
			comments = append(comments, models.OrderComment{
				OrderCommentUUID: uuid.New(),
				OrderID:          orderID,
				Name:             txt,
				Status:           1,
				CreatedAt:        &now,
				UpdatedAt:        &now,
				CreatedBy:        actor,
			})
		}
		return s.repo.InsertComments(txCtx, comments)
	})
	if err != nil {
		s.recordErr(ctx, span, "create order", err)
		return nil, err
	}

	// Muat ulang lengkap (relasi + anak) untuk respons.
	full, err := s.repo.FindByID(ctx, companyUUID, m.ID)
	if err != nil {
		s.recordErr(ctx, span, "reload order", err)
		return nil, err
	}
	if full != nil {
		return full, nil
	}
	return m, nil
}

// derivePeriodDate menentukan tanggal periode (partition key): dari PeriodDate eksplisit,
// lalu StartDate, lalu hari ini — dinormalkan ke tengah malam UTC.
func derivePeriodDate(req structs.OrderCreateRequest) time.Time {
	switch {
	case req.PeriodDate != nil:
		return truncateDate(*req.PeriodDate)
	case req.StartDate != nil:
		return truncateDate(*req.StartDate)
	default:
		return truncateDate(time.Now())
	}
}

func truncateDate(t time.Time) time.Time {
	y, mo, d := t.UTC().Date()
	return time.Date(y, mo, d, 0, 0, 0, 0, time.UTC)
}

func parseUUIDs(ss []string) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(ss))
	for _, s := range ss {
		if id, err := uuid.Parse(s); err == nil {
			out = append(out, id)
		}
	}
	return out
}

func parseUUIDPtr(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (s *orderService) recordErr(ctx context.Context, span *monitoring.APMSpan, msg string, err error) {
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrPriorityNotFound) ||
		errors.Is(err, ErrStatusNotFound) || errors.Is(err, ErrCategoryNotFound) {
		return
	}
	monitoring.RecordSpanError(span, monitoring.LogErrorService, monitoring.SpanLayerService, monitoring.OpRepositoryCall, err)
	monitoring.ContextLogger(ctx, s.logger).ErrorContext(ctx, msg+" failed",
		slog.String(monitoring.LabelErrorLayer, monitoring.SpanLayerService),
		slog.String(monitoring.LabelErrorMessage, err.Error()))
}
