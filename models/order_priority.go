package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// OrderPriority = referensi prioritas order (None, Low, Medium, High).
type OrderPriority struct {
	bun.BaseModel `bun:"table:t_order_priority,alias:op"`

	ID                int64      `bun:"id,pk,autoincrement"`
	OrderPriorityUUID uuid.UUID  `bun:"order_priority_uuid,type:uuid,notnull"`
	Name              string     `bun:"name,notnull"`
	Status            int16      `bun:"status,notnull,default:1"`
	CreatedAt         *time.Time `bun:"created_at"`
	UpdatedAt         *time.Time `bun:"updated_at"`
}
