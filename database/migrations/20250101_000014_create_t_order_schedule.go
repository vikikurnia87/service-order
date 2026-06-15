package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_schedule")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_schedule (
				id                  BIGSERIAL    PRIMARY KEY,
				order_schedule_uuid UUID         NOT NULL UNIQUE,
				schedule_id         BIGINT       NOT NULL REFERENCES t_schedule(id),
				date_id             BIGINT       REFERENCES t_date(id),
				repeat_every        SMALLINT     NOT NULL DEFAULT 1,
				company_uuid        UUID         NOT NULL,
				status              SMALLINT     NOT NULL DEFAULT 1,
				created_at          TIMESTAMPTZ  DEFAULT NOW(),
				updated_at          TIMESTAMPTZ  DEFAULT NOW(),
				created_by          UUID,
				updated_by          UUID
			);
			CREATE INDEX IF NOT EXISTS idx_t_order_schedule_company ON t_order_schedule(company_uuid);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_schedule")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_schedule CASCADE`)
		return err
	})
}
