package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// MasterDate = kalender (1 baris per tanggal) untuk mesin penjadwalan: memetakan
// tanggal kalender ke hari/minggu/bulan/tahun. Di-generate programatik per tahun.
type MasterDate struct {
	bun.BaseModel `bun:"table:t_master_date,alias:md"`

	ID             int64      `bun:"id,pk,autoincrement"`
	MasterDateUUID uuid.UUID  `bun:"master_date_uuid,type:uuid,notnull"`
	FullDate       time.Time  `bun:"full_date,type:date,notnull"` // tanggal kalender, mis. 2025-06-15
	DayNo          int16      `bun:"day_no,notnull"`              // 1=Minggu … 7=Sabtu
	Day            string     `bun:"day,notnull"`                 // sunday … saturday
	DayOfMonth     int16      `bun:"day_of_month,notnull"`        // 1 … 31
	Week           int16      `bun:"week,notnull"`                // ISO week 1 … 53
	Month          int16      `bun:"month,notnull"`               // 1 … 12
	Year           int16      `bun:"year,notnull"`
	Status         int16      `bun:"status,notnull,default:1"`
	CreatedAt      *time.Time `bun:"created_at"`
	UpdatedAt      *time.Time `bun:"updated_at"`
}
