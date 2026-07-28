package mappers

import "github.com/vikikurnia87/service-order/models"

// MapMasterDate — representasi publik baris kalender.
func MapMasterDate(m *models.MasterDate) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":               m.ID,
		"master_date_uuid": m.MasterDateUUID.String(),
		"full_date":        m.FullDate,
		"day_no":           m.DayNo,
		"day":              m.Day,
		"day_of_month":     m.DayOfMonth,
		"week":             m.Week,
		"month":            m.Month,
		"year":             m.Year,
		"status":           m.Status,
	}
}
