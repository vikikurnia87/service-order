package mappers

import (
	"github.com/google/uuid"
	"github.com/vikikurnia87/service-order/models"
)

// MapOrderList = ringkasan order untuk list (tanpa anak, dengan nama status/priority
// bila relasi di-load).
func MapOrderList(m *models.Order) map[string]any {
	out := map[string]any{
		"id":            m.ID,
		"order_uuid":    m.OrderUUID.String(),
		"order_number":  m.OrderNumber,
		"name":          m.Name,
		"period_date":   m.PeriodDate,
		"start_date":    m.StartDate,
		"due_date":      m.DueDate,
		"complete_date": m.CompleteDate,
		"priority":      priorityInfo(m),
		"status":        statusInfo(m),
		"status_code":   m.Status,
		"created_at":    m.CreatedAt,
		"updated_at":    m.UpdatedAt,
	}
	return out
}

// MapOrder = detail order lengkap dengan anak (kategori/vendor/assignee/komentar).
func MapOrder(m *models.Order) map[string]any {
	out := MapOrderList(m)
	out["company_uuid"] = m.CompanyUUID.String()
	out["description"] = m.Description
	out["asset_uuid"] = uuidPtrStr(m.AssetUUID)
	out["location_uuid"] = uuidPtrStr(m.LocationUUID)

	cats := make([]map[string]any, 0, len(m.Categories))
	for i := range m.Categories {
		c := &m.Categories[i]
		item := map[string]any{
			"order_category_uuid": c.OrderCategoryUUID.String(),
			"category_id":         c.CategoryID,
		}
		if c.Category != nil {
			item["category_uuid"] = c.Category.CategoryUUID.String()
			item["category_name"] = c.Category.Name
		}
		cats = append(cats, item)
	}
	out["categories"] = cats

	vendors := make([]map[string]any, 0, len(m.Vendors))
	for i := range m.Vendors {
		v := &m.Vendors[i]
		vendors = append(vendors, map[string]any{
			"order_vendor_uuid": v.OrderVendorUUID.String(),
			"vendor_uuid":       v.VendorUUID.String(),
		})
	}
	out["vendors"] = vendors

	assignees := make([]map[string]any, 0, len(m.Assignees))
	for i := range m.Assignees {
		a := &m.Assignees[i]
		assignees = append(assignees, map[string]any{
			"order_assign_uuid": a.OrderAssignUUID.String(),
			"team_uuid":         uuidPtrStr(a.TeamUUID),
			"user_uuid":         uuidPtrStr(a.UserUUID),
		})
	}
	out["assignees"] = assignees

	comments := make([]map[string]any, 0, len(m.Comments))
	for i := range m.Comments {
		cm := &m.Comments[i]
		comments = append(comments, map[string]any{
			"order_comment_uuid": cm.OrderCommentUUID.String(),
			"name":               cm.Name,
			"created_at":         cm.CreatedAt,
		})
	}
	out["comments"] = comments

	return out
}

func statusInfo(m *models.Order) map[string]any {
	out := map[string]any{"id": m.OrderStatusID}
	if m.OrderStatus != nil {
		out["uuid"] = m.OrderStatus.OrderStatusUUID.String()
		out["name"] = m.OrderStatus.Name
	}
	return out
}

func priorityInfo(m *models.Order) map[string]any {
	out := map[string]any{"id": m.OrderPriorityID}
	if m.OrderPriority != nil {
		out["uuid"] = m.OrderPriority.OrderPriorityUUID.String()
		out["name"] = m.OrderPriority.Name
	}
	return out
}

func uuidPtrStr(u *uuid.UUID) any {
	if u == nil {
		return nil
	}
	return u.String()
}
