// Package seeders berisi daftar seeder service-order. Runner generik (RunAll/RunOne)
// dari service-utils/database (DRY).
package seeders

import (
	"github.com/uptrace/bun"
	sudb "github.com/vikikurnia87/service-utils/database"
)

// Seeder = kontrak seeder generik (Name + Run).
type Seeder = sudb.Seeder

// NewDatabaseSeeder membangun runner dengan urutan seeder referensi.
func NewDatabaseSeeder(db *bun.DB) *sudb.SeederRunner {
	return sudb.NewSeederRunner(db,
		NewOrderStatusSeeder(),
		NewOrderPrioritySeeder(),
		NewScheduleSeeder(),
		NewDaySeeder(),
		NewDateSeeder(),
		NewMasterDateSeeder(),
	)
}

// nameSet membangun set nama untuk cek idempotensi seeder.
func nameSet(names []string) map[string]bool {
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return m
}
