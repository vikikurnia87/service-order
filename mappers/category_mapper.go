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
