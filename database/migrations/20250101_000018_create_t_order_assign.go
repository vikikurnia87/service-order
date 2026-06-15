package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_assign")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_assign (
				id                BIGSERIAL   PRIMARY KEY,
				order_assign_uuid UUID        NOT NULL UNIQUE,
				order_id          BIGINT      NOT NULL,
				team_uuid         UUID,
				user_uuid         UUID,
				status            SMALLINT    NOT NULL DEFAULT 1,
				created_at        TIMESTAMPTZ DEFAULT NOW(),
				updated_at        TIMESTAMPTZ DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_t_order_assign_order ON t_order_assign(order_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_assign")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_assign CASCADE`)
		return err
	})
}
