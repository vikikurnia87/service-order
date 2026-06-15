// Package seeders berisi daftar seeder service-order. Runner generik (RunAll/RunOne)
// dari service-utils/database (DRY).
package seeders

import (
	"github.com/uptrace/bun"
	sudb "github.com/vikikurnia87/service-utils/database"
)

// Seeder = kontrak seeder generik (Name + Run).
type Seeder = sudb.Seeder

// NewDatabaseSeeder membangun runner. Seeder referensi (status/priority/schedule/
// day/date/master_date) ditambahkan pada Fase 1.
func NewDatabaseSeeder(db *bun.DB) *sudb.SeederRunner {
	return sudb.NewSeederRunner(db)
}
