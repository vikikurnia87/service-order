package seeders

import (
	"context"
	"time"

	"github.com/vikikurnia87/service-order/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DaySeeder struct{}

func NewDaySeeder() *DaySeeder { return &DaySeeder{} }

func (s *DaySeeder) Name() string { return "DaySeeder" }

func (s *DaySeeder) Run(ctx context.Context, db *bun.DB) error {
	records, err := ReadCSV("day.csv")
	if err != nil {
		return err
	}
	var existing []string
	if err := db.NewSelect().Model((*models.Day)(nil)).Column("name").Scan(ctx, &existing); err != nil {
		return err
	}
	seen := nameSet(existing)

	now := time.Now()
	rows := make([]models.Day, 0, len(records))
	for _, r := range records {
		name := GetString(r, "dayname", "")
		if name == "" || seen[name] {
			continue
		}
		id := GetUUID(r, "dayuuid")
		if id == uuid.Nil {
			id = uuid.New()
		}
		rows = append(rows, models.Day{DayUUID: id, Name: name, Status: 1, CreatedAt: &now, UpdatedAt: &now})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err = db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
