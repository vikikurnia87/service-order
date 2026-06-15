package structs

import "time"

// OrderAssignInput = satu penugasan (tim dan/atau user). Minimal salah satu diisi
// (divalidasi di service).
type OrderAssignInput struct {
	TeamUUID string `json:"team_uuid" validate:"omitempty,uuid4"`
	UserUUID string `json:"user_uuid" validate:"omitempty,uuid4"`
}

// OrderCreateRequest = payload buat order baru.
//
// Mengacu order lama (ordername/orderpriorityuuid/orderstatusuuid/category/vendor/...)
// tetapi memakai snake_case + UUID v4. Snapshot prosedur (procedure_uuid) & penjadwalan
// (schedule/day/date/repeat_every) DITUNDA ke Fase 3/4.
type OrderCreateRequest struct {
	Name         string `json:"name" validate:"required,max=255"`
	Description  string `json:"description"`
	OrderNumber  string `json:"order_number" validate:"omitempty,max=100"`
	PriorityUUID string `json:"priority_uuid" validate:"required,uuid4"`
	StatusUUID   string `json:"status_uuid" validate:"required,uuid4"`
	AssetUUID    string `json:"asset_uuid" validate:"omitempty,uuid4"`
	LocationUUID string `json:"location_uuid" validate:"omitempty,uuid4"`

	StartDate *time.Time `json:"start_date"`
	DueDate   *time.Time `json:"due_date"`
	// PeriodDate = tanggal periode (partition key). Opsional; default dari StartDate
	// atau hari ini bila kosong.
	PeriodDate *time.Time `json:"period_date"`

	Categories []string           `json:"categories" validate:"omitempty,dive,uuid4"`
	Vendors    []string           `json:"vendors" validate:"omitempty,dive,uuid4"`
	Assignees  []OrderAssignInput `json:"assignees" validate:"omitempty,dive"`
	Comments   []string           `json:"comments" validate:"omitempty,dive,max=2000"`
}
