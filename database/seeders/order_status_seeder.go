package seeders

import (
	"context"
	"time"

	"github.com/vikikurnia87/service-order/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrderStatusSeeder struct{}

func NewOrderStatusSeeder() *OrderStatusSeeder { return &OrderStatusSeeder{} }

func (s *OrderStatusSeeder) Name() string { return "OrderStatusSeeder" }

func (s *OrderStatusSeeder) Run(ctx context.Context, db *bun.DB) error {
	records, err := ReadCSV("order_status.csv")
	if err != nil {
		return err
	}
	var existing []string
	if err := db.NewSelect().Model((*models.OrderStatus)(nil)).Column("name").Scan(ctx, &existing); err != nil {
		return err
	}
	seen := nameSet(existing)

	now := time.Now()
	rows := make([]models.OrderStatus, 0, len(records))
	for _, r := range records {
		name := GetString(r, "orderstatusname", "")
		if name == "" || seen[name] {
			continue
		}
		id := GetUUID(r, "orderstatusuuid")
		if id == uuid.Nil {
			id = uuid.New()
		}
		rows = append(rows, models.OrderStatus{OrderStatusUUID: id, Name: name, Status: 1, CreatedAt: &now, UpdatedAt: &now})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err = db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
