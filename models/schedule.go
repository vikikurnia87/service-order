package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Schedule = referensi pola jadwal berulang (None, Daily, Weekly, Monthly, Yearly).
type Schedule struct {
	bun.BaseModel `bun:"table:t_schedule,alias:sc"`

	ID           int64      `bun:"id,pk,autoincrement"`
	ScheduleUUID uuid.UUID  `bun:"schedule_uuid,type:uuid,notnull"`
	Name         string     `bun:"name,notnull"`
	Status       int16      `bun:"status,notnull,default:1"`
	CreatedAt    *time.Time `bun:"created_at"`
	UpdatedAt    *time.Time `bun:"updated_at"`
}
