package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// OrderAssign = penugasan order ke tim/user. team_uuid & user_uuid referensi
// lintas-service (service-user) TANPA FK.
type OrderAssign struct {
	bun.BaseModel `bun:"table:t_order_assign,alias:oa"`

	ID              int64      `bun:"id,pk,autoincrement"`
	OrderAssignUUID uuid.UUID  `bun:"order_assign_uuid,type:uuid,notnull"`
	OrderID         int64      `bun:"order_id,notnull"`
	TeamUUID        *uuid.UUID `bun:"team_uuid,type:uuid"`
	UserUUID        *uuid.UUID `bun:"user_uuid,type:uuid"`
	Status          int16      `bun:"status,notnull,default:1"`
	CreatedAt       *time.Time `bun:"created_at"`
	UpdatedAt       *time.Time `bun:"updated_at"`
}
