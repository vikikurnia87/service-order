package mappers

import "github.com/vikikurnia87/service-order/models"

// MapDate — representasi publik referensi tanggal dalam bulan.
func MapDate(m *models.Date) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":        m.ID,
		"date_uuid": m.DateUUID.String(),
		"name":      m.Name,
		"value":     m.Value,
		"status":    m.Status,
	}
}
