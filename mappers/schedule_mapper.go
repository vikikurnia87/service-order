package mappers

import "github.com/vikikurnia87/service-order/models"

// MapSchedule — representasi publik pola jadwal.
func MapSchedule(m *models.Schedule) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":            m.ID,
		"schedule_uuid": m.ScheduleUUID.String(),
		"name":          m.Name,
		"status":        m.Status,
	}
}
