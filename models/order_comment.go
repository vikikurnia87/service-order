package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// OrderComment = catatan/komentar pada order. order_id TANPA FK (t_order terpartisi).
type OrderComment struct {
	bun.BaseModel `bun:"table:t_order_comment,alias:ocm"`

	ID               int64      `bun:"id,pk,autoincrement"`
	OrderCommentUUID uuid.UUID  `bun:"order_comment_uuid,type:uuid,notnull"`
	OrderID          int64      `bun:"order_id,notnull"`
	Name             string     `bun:"name,notnull"` // isi komentar
	Status           int16      `bun:"status,notnull,default:1"`
	CreatedAt        *time.Time `bun:"created_at"`
	UpdatedAt        *time.Time `bun:"updated_at"`
	CreatedBy        *uuid.UUID `bun:"created_by,type:uuid"`
}
