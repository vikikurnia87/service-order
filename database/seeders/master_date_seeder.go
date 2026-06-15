package seeders

import (
	"context"
	"strings"
	"time"

	"github.com/vikikurnia87/service-order/configs"
	"github.com/vikikurnia87/service-order/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// MasterDateSeeder mengisi kalender t_master_date secara programatik (bukan CSV
// raksasa) dari MASTERDATE_START_YEAR s/d MASTERDATE_END_YEAR (default: tahun
// berjalan s/d +2). Idempotent: full_date UNIQUE → ON CONFLICT DO NOTHING.
type MasterDateSeeder struct{}

func NewMasterDateSeeder() *MasterDateSeeder { return &MasterDateSeeder{} }

func (s *MasterDateSeeder) Name() string { return "MasterDateSeeder" }

func (s *MasterDateSeeder) Run(ctx context.Context, db *bun.DB) error {
	startYear := configs.GetEnvAsInt("MASTERDATE_START_YEAR", time.Now().Year())
	endYear := configs.GetEnvAsInt("MASTERDATE_END_YEAR", time.Now().Year()+2)
	if endYear < startYear {
		endYear = startYear
	}

	start := time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(endYear, time.December, 31, 0, 0, 0, 0, time.UTC)

	now := time.Now()
	rows := make([]models.MasterDate, 0, 512)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		_, isoWeek := d.ISOWeek()
		rows = append(rows, models.MasterDate{
			MasterDateUUID: uuid.New(),
			FullDate:       d,
			DayNo:          int16(d.Weekday()) + 1, // Sunday=0 → 1 … Saturday=6 → 7
			Day:            strings.ToLower(d.Weekday().String()),
			DayOfMonth:     int16(d.Day()),
			Week:           int16(isoWeek),
			Month:          int16(d.Month()),
			Year:           int16(d.Year()),
			Status:         1,
			CreatedAt:      &now,
			UpdatedAt:      &now,
		})
	}

	const batch = 500
	for i := 0; i < len(rows); i += batch {
		j := i + batch
		if j > len(rows) {
			j = len(rows)
		}
		chunk := rows[i:j]
		if _, err := db.NewInsert().Model(&chunk).On("CONFLICT (full_date) DO NOTHING").Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}
