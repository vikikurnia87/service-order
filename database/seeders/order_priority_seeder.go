package seeders

import (
	"context"
	"time"

	"github.com/vikikurnia87/service-order/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrderPrioritySeeder struct{}

func NewOrderPrioritySeeder() *OrderPrioritySeeder { return &OrderPrioritySeeder{} }

func (s *OrderPrioritySeeder) Name() string { return "OrderPrioritySeeder" }

func (s *OrderPrioritySeeder) Run(ctx context.Context, db *bun.DB) error {
	records, err := ReadCSV("order_priority.csv")
	if err != nil {
		return err
	}
	var existing []string
	if err := db.NewSelect().Model((*models.OrderPriority)(nil)).Column("name").Scan(ctx, &existing); err != nil {
		return err
	}
	seen := nameSet(existing)

	now := time.Now()
	rows := make([]models.OrderPriority, 0, len(records))
	for _, r := range records {
		name := GetString(r, "orderpriorityname", "")
		if name == "" || seen[name] {
			continue
		}
		id := GetUUID(r, "orderpriorityuuid")
		if id == uuid.Nil {
			id = uuid.New()
		}
		rows = append(rows, models.OrderPriority{OrderPriorityUUID: id, Name: name, Status: 1, CreatedAt: &now, UpdatedAt: &now})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err = db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
