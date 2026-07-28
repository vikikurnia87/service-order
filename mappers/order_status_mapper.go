package mappers

import "github.com/vikikurnia87/service-order/models"

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
