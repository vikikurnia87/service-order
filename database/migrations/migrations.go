// Package migrations berisi registry migrasi service-vendor. Runner generik
// (Migrate/Rollback/Reset/Status/Fresh) berasal dari service-utils/database (DRY).
package migrations

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	sudb "github.com/vikikurnia87/service-utils/database"
)

// Migrations adalah registry domain — tiap file migration register via init().
var Migrations = migrate.NewMigrations()

// NewMigrator mengembalikan runner service-utils atas registry domain ini.
func NewMigrator(db *bun.DB) *sudb.Migrator {
	return sudb.NewMigrator(db, Migrations)
}
