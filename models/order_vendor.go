package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// OrderVendor = vendor yang ditugaskan ke order. vendor_uuid referensi lintas-service
// (service-vendor) TANPA FK.
type OrderVendor struct {
	bun.BaseModel `bun:"table:t_order_vendor,alias:ov"`

	ID              int64      `bun:"id,pk,autoincrement"`
	OrderVendorUUID uuid.UUID  `bun:"order_vendor_uuid,type:uuid,notnull"`
	OrderID         int64      `bun:"order_id,notnull"`
	VendorUUID      uuid.UUID  `bun:"vendor_uuid,type:uuid,notnull"`
	Status          int16      `bun:"status,notnull,default:1"`
	CreatedAt       *time.Time `bun:"created_at"`
	UpdatedAt       *time.Time `bun:"updated_at"`
}
