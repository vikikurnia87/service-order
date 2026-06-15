package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// OrderCategory = junction order ↔ category. order_id TANPA FK (t_order terpartisi);
// category_id FK ke t_category.
type OrderCategory struct {
	bun.BaseModel `bun:"table:t_order_category,alias:oc"`

	ID                int64      `bun:"id,pk,autoincrement"`
	OrderCategoryUUID uuid.UUID  `bun:"order_category_uuid,type:uuid,notnull"`
	OrderID           int64      `bun:"order_id,notnull"`
	CategoryID        int64      `bun:"category_id,notnull"`
	Status            int16      `bun:"status,notnull,default:1"`
	CreatedAt         *time.Time `bun:"created_at"`
	UpdatedAt         *time.Time `bun:"updated_at"`

	Category *Category `bun:"rel:belongs-to,join:category_id=id"`
}
