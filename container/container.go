// Package container memegang shared dependency (repos + services) service-order.
package container

import (
	"log/slog"

	"github.com/vikikurnia87/service-order/database"

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
// Entitas (order/procedure-snapshot/schedule/...) ditambahkan pada Fase 1+.
type Container struct {
	DB     *bun.DB
	Logger *slog.Logger

	// Cache (tenant-aware pagination).
	Cache *sucache.Cache

	// Transaction helper — wrapping multi-step CUD atomik.
	TxHelper dbutil.TransactionHelper
}

func New(deps Deps) *Container {
	db, lg := deps.DB, deps.Logger
	return &Container{
		DB:       db,
		Logger:   lg,
		Cache:    database.GetCache(),
		TxHelper: dbutil.NewTransactionHelper(db, lg),
	}
}
