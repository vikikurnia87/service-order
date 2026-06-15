package seeders

import (
	"context"
	"time"

	"github.com/vikikurnia87/service-order/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DateSeeder struct{}

func NewDateSeeder() *DateSeeder { return &DateSeeder{} }

func (s *DateSeeder) Name() string { return "DateSeeder" }

func (s *DateSeeder) Run(ctx context.Context, db *bun.DB) error {
	records, err := ReadCSV("date.csv")
	if err != nil {
		return err
	}
	var existing []string
	if err := db.NewSelect().Model((*models.Date)(nil)).Column("name").Scan(ctx, &existing); err != nil {
		return err
	}
	seen := nameSet(existing)

	now := time.Now()
	rows := make([]models.Date, 0, len(records))
	for _, r := range records {
		name := GetString(r, "datename", "")
		if name == "" || seen[name] {
			continue
		}
		id := GetUUID(r, "dateuuid")
		if id == uuid.Nil {
			id = uuid.New()
		}
		rows = append(rows, models.Date{
			DateUUID:  id,
			Name:      name,
			Value:     GetString(r, "value", ""),
			Status:    1,
			CreatedAt: &now,
			UpdatedAt: &now,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err = db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
