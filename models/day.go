package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Day = referensi hari (sunday … saturday) untuk jadwal mingguan.
type Day struct {
	bun.BaseModel `bun:"table:t_day,alias:dy"`

	ID        int64      `bun:"id,pk,autoincrement"`
	DayUUID   uuid.UUID  `bun:"day_uuid,type:uuid,notnull"`
	Name      string     `bun:"name,notnull"`
	Status    int16      `bun:"status,notnull,default:1"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
}
