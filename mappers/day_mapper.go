package mappers

import "github.com/vikikurnia87/service-order/models"

// MapDay — representasi publik referensi hari.
func MapDay(m *models.Day) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":       m.ID,
		"day_uuid": m.DayUUID.String(),
		"name":     m.Name,
		"status":   m.Status,
	}
}
