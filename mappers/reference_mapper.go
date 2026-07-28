package mappers

import "github.com/vikikurnia87/service-order/models"

// MapCategory — representasi publik master kategori order.
func MapCategory(m *models.Category) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":            m.ID,
		"category_uuid": m.CategoryUUID.String(),
		"name":          m.Name,
		"description":   m.Description,
		"company_uuid":  m.CompanyUUID.String(),
		"status":        m.Status,
		"created_at":    m.CreatedAt,
		"updated_at":    m.UpdatedAt,
	}
}

// MapOrderStatus — representasi publik status order.
func MapOrderStatus(m *models.OrderStatus) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":                m.ID,
		"order_status_uuid": m.OrderStatusUUID.String(),
		"name":              m.Name,
		"status":            m.Status,
	}
}

// MapOrderPriority — representasi publik prioritas order.
func MapOrderPriority(m *models.OrderPriority) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":                  m.ID,
		"order_priority_uuid": m.OrderPriorityUUID.String(),
		"name":                m.Name,
		"status":              m.Status,
	}
}

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
