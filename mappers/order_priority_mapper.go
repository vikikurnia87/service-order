package mappers

import "github.com/vikikurnia87/service-order/models"

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
