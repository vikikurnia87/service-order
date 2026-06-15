package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Order = work order / inspeksi (tabel inti t_order, DIPARTISI RANGE per period_date).
// Anak (Categories/Vendors/Assignees/Comments) di-join via order_id TANPA FK karena
// FK ke tabel terpartisi rumit — integritas dijaga di layer service.
//
// Catatan: OrderProcedureID/OrderScheduleID (snapshot form + jadwal) diisi pada
// Fase 3/4; untuk saat ini create order hanya mengisi inti + anak sederhana.
type Order struct {
	bun.BaseModel `bun:"table:t_order,alias:o"`

	ID                  int64      `bun:"id,pk,autoincrement"`
	OrderUUID           uuid.UUID  `bun:"order_uuid,type:uuid,notnull"`
	CompanyUUID         uuid.UUID  `bun:"company_uuid,type:uuid,notnull"`
	OrderNumber         *string    `bun:"order_number"`
	Name                string     `bun:"name,notnull"`
	Description         *string    `bun:"description"`
	OrderProcedureID    *int64     `bun:"order_procedure_id"`
	OrderScheduleID     *int64     `bun:"order_schedule_id"`
	OrderPriorityID     int64      `bun:"order_priority_id,notnull"`
	OrderStatusID       int64      `bun:"order_status_id,notnull"`
	AssetUUID           *uuid.UUID `bun:"asset_uuid,type:uuid"`
	LocationUUID        *uuid.UUID `bun:"location_uuid,type:uuid"`
	ProcedureUUIDOrigin *uuid.UUID `bun:"procedure_uuid_origin,type:uuid"`
	StartDate           *time.Time `bun:"start_date"`
	DueDate             *time.Time `bun:"due_date"`
	CompleteDate        *time.Time `bun:"complete_date"`
	PeriodDate          time.Time  `bun:"period_date,type:date,notnull"` // partition key (tanggal periode)
	IsOrder             bool       `bun:"is_order,notnull,default:true"`
	IsSystem            bool       `bun:"is_system,notnull,default:false"`
	Status              int16      `bun:"status,notnull,default:1"`
	CreatedAt           *time.Time `bun:"created_at"`
	UpdatedAt           *time.Time `bun:"updated_at"`
	CreatedBy           *uuid.UUID `bun:"created_by,type:uuid"`
	UpdatedBy           *uuid.UUID `bun:"updated_by,type:uuid"`

	// Relasi (di-load saat detail/list).
	OrderStatus   *OrderStatus    `bun:"rel:belongs-to,join:order_status_id=id"`
	OrderPriority *OrderPriority  `bun:"rel:belongs-to,join:order_priority_id=id"`
	Categories    []OrderCategory `bun:"rel:has-many,join:id=order_id"`
	Vendors       []OrderVendor   `bun:"rel:has-many,join:id=order_id"`
	Assignees     []OrderAssign   `bun:"rel:has-many,join:id=order_id"`
	Comments      []OrderComment  `bun:"rel:has-many,join:id=order_id"`
}
