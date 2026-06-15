package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Category = master kategori order (ter-scope tenant via company_uuid).
type Category struct {
	bun.BaseModel `bun:"table:t_category,alias:cat"`

	ID           int64      `bun:"id,pk,autoincrement"`
	CategoryUUID uuid.UUID  `bun:"category_uuid,type:uuid,notnull"`
	Name         string     `bun:"name,notnull"`
	Description  *string    `bun:"description"`
	CompanyUUID  uuid.UUID  `bun:"company_uuid,type:uuid,notnull"`
	Status       int16      `bun:"status,notnull,default:1"`
	CreatedAt    *time.Time `bun:"created_at"`
	UpdatedAt    *time.Time `bun:"updated_at"`
	CreatedBy    *uuid.UUID `bun:"created_by,type:uuid"`
	UpdatedBy    *uuid.UUID `bun:"updated_by,type:uuid"`
}
