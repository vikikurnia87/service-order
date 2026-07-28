// Package container memegang shared dependency (repos + services) service-order.
package container

import (
	"log/slog"

	"github.com/vikikurnia87/service-order/database"
	"github.com/vikikurnia87/service-order/repositories"
	"github.com/vikikurnia87/service-order/services"

	"github.com/uptrace/bun"
	sucache "github.com/vikikurnia87/service-utils/cache"
	"github.com/vikikurnia87/service-utils/dbutil"
)

// Deps adalah parameter untuk New().
type Deps struct {
	DB     *bun.DB
	Logger *slog.Logger
}

// Container memegang shared dependency antar handler.
type Container struct {
	DB     *bun.DB
	Logger *slog.Logger

	// Cache (tenant-aware pagination).
	Cache *sucache.Cache

	// Transaction helper — wrapping multi-step CUD atomik.
	TxHelper dbutil.TransactionHelper

	// Order (inti).
	OrderRepo    repositories.OrderRepository
	OrderService services.OrderService

	// Kategori (ter-scope tenant).
	CategoryRepo    repositories.CategoryRepository
	CategoryService services.CategoryService

	// Referensi global (read-only).
	OrderStatusRepo    repositories.OrderStatusRepository
	OrderStatusService services.OrderStatusService

	OrderPriorityRepo    repositories.OrderPriorityRepository
	OrderPriorityService services.OrderPriorityService

	ScheduleRepo    repositories.ScheduleRepository
	ScheduleService services.ScheduleService

	DayRepo    repositories.DayRepository
	DayService services.DayService

	DateRepo    repositories.DateRepository
	DateService services.DateService

	MasterDateRepo    repositories.MasterDateRepository
	MasterDateService services.MasterDateService
}

func New(deps Deps) *Container {
	db, lg := deps.DB, deps.Logger
	txHelper := dbutil.NewTransactionHelper(db, lg)

	// Order (inti).
	orderRepo := repositories.NewOrderRepository(db, lg)
	orderService := services.NewOrderService(orderRepo, txHelper, lg)

	// Kategori.
	categoryRepo := repositories.NewCategoryRepository(db, lg)
	categoryService := services.NewCategoryService(categoryRepo, txHelper, lg)

	// Referensi global.
	orderStatusRepo := repositories.NewOrderStatusRepository(db, lg)
	orderStatusService := services.NewOrderStatusService(orderStatusRepo, lg)

	orderPriorityRepo := repositories.NewOrderPriorityRepository(db, lg)
	orderPriorityService := services.NewOrderPriorityService(orderPriorityRepo, lg)

	scheduleRepo := repositories.NewScheduleRepository(db, lg)
	scheduleService := services.NewScheduleService(scheduleRepo, lg)

	dayRepo := repositories.NewDayRepository(db, lg)
	dayService := services.NewDayService(dayRepo, lg)

	dateRepo := repositories.NewDateRepository(db, lg)
	dateService := services.NewDateService(dateRepo, lg)

	masterDateRepo := repositories.NewMasterDateRepository(db, lg)
	masterDateService := services.NewMasterDateService(masterDateRepo, lg)

	return &Container{
		DB:       db,
		Logger:   lg,
		Cache:    database.GetCache(),
		TxHelper: txHelper,

		OrderRepo:    orderRepo,
		OrderService: orderService,

		CategoryRepo:    categoryRepo,
		CategoryService: categoryService,

		OrderStatusRepo:    orderStatusRepo,
		OrderStatusService: orderStatusService,

		OrderPriorityRepo:    orderPriorityRepo,
		OrderPriorityService: orderPriorityService,

		ScheduleRepo:    scheduleRepo,
		ScheduleService: scheduleService,

		DayRepo:    dayRepo,
		DayService: dayService,

		DateRepo:    dateRepo,
		DateService: dateService,

		MasterDateRepo:    masterDateRepo,
		MasterDateService: masterDateService,
	}
}
