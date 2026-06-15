// Package database (service-vendor) adalah facade tipis di atas
// service-utils/database agar logika koneksi/migrasi tidak diduplikasi (DRY).
package database

import (
	"log/slog"

	sudb "github.com/vikikurnia87/service-utils/database"

	"github.com/uptrace/bun"
)

func NewPostgresDB(cfg sudb.Config, env string, logger *slog.Logger) (*bun.DB, error) {
	return sudb.NewPostgresDB(cfg, env, logger)
}

func InitPostgresDatabase(cfg sudb.Config, env string, logger *slog.Logger) *bun.DB {
	return sudb.InitPostgresDatabase(cfg, env, logger)
}

func Close(db *bun.DB) error {
	return sudb.Close(db)
}
