package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// OrderStatus = referensi status order (Open, On Hold, In Progress, Done).
type OrderStatus struct {
	bun.BaseModel `bun:"table:t_order_status,alias:os"`

	ID              int64      `bun:"id,pk,autoincrement"`
	OrderStatusUUID uuid.UUID  `bun:"order_status_uuid,type:uuid,notnull"`
	Name            string     `bun:"name,notnull"`
	Status          int16      `bun:"status,notnull,default:1"`
	CreatedAt       *time.Time `bun:"created_at"`
	UpdatedAt       *time.Time `bun:"updated_at"`
}
