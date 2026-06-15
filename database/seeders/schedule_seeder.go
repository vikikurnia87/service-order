package seeders

import (
	"context"
	"time"

	"github.com/vikikurnia87/service-order/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ScheduleSeeder struct{}

func NewScheduleSeeder() *ScheduleSeeder { return &ScheduleSeeder{} }

func (s *ScheduleSeeder) Name() string { return "ScheduleSeeder" }

func (s *ScheduleSeeder) Run(ctx context.Context, db *bun.DB) error {
	records, err := ReadCSV("schedule.csv")
	if err != nil {
		return err
	}
	var existing []string
	if err := db.NewSelect().Model((*models.Schedule)(nil)).Column("name").Scan(ctx, &existing); err != nil {
		return err
	}
	seen := nameSet(existing)

	now := time.Now()
	rows := make([]models.Schedule, 0, len(records))
	for _, r := range records {
		name := GetString(r, "schedulename", "")
		if name == "" || seen[name] {
			continue
		}
		id := GetUUID(r, "scheduleuuid")
		if id == uuid.Nil {
			id = uuid.New()
		}
		rows = append(rows, models.Schedule{ScheduleUUID: id, Name: name, Status: 1, CreatedAt: &now, UpdatedAt: &now})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err = db.NewInsert().Model(&rows).Exec(ctx)
	return err
}
