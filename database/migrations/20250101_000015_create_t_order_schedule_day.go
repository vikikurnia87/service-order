package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_schedule_day")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_schedule_day (
				id                      BIGSERIAL    PRIMARY KEY,
				order_schedule_day_uuid UUID         NOT NULL UNIQUE,
				order_schedule_id       BIGINT       NOT NULL REFERENCES t_order_schedule(id) ON DELETE CASCADE,
				day_id                  BIGINT       NOT NULL REFERENCES t_day(id),
				status                  SMALLINT     NOT NULL DEFAULT 1,
				created_at              TIMESTAMPTZ  DEFAULT NOW(),
				updated_at              TIMESTAMPTZ  DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_t_osd_schedule ON t_order_schedule_day(order_schedule_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_schedule_day")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_schedule_day CASCADE`)
		return err
	})
}
