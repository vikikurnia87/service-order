package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Date = referensi tanggal dalam bulan (1st … 31st, value 01 … 31) untuk jadwal bulanan.
type Date struct {
	bun.BaseModel `bun:"table:t_date,alias:dt"`

	ID        int64      `bun:"id,pk,autoincrement"`
	DateUUID  uuid.UUID  `bun:"date_uuid,type:uuid,notnull"`
	Name      string     `bun:"name,notnull"`  // "1st", "2nd", …
	Value     string     `bun:"value,notnull"` // "01", "02", …
	Status    int16      `bun:"status,notnull,default:1"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
}
