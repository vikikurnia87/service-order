package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_transaction_schedule")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_transaction_schedule (
				id                        BIGSERIAL   PRIMARY KEY,
				transaction_schedule_uuid UUID        NOT NULL UNIQUE,
				order_id                  BIGINT      NOT NULL,
				start_date                TIMESTAMPTZ,
				end_date                  TIMESTAMPTZ,
				master_date_id            BIGINT      REFERENCES t_master_date(id),
				schedule_id               BIGINT      REFERENCES t_schedule(id),
				status                    SMALLINT    NOT NULL DEFAULT 1,
				created_at                TIMESTAMPTZ DEFAULT NOW(),
				updated_at                TIMESTAMPTZ DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_t_txsched_order ON t_transaction_schedule(order_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_transaction_schedule")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_transaction_schedule CASCADE`)
		return err
	})
}
