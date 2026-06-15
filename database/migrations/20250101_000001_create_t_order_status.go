package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_status")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_status (
				id          BIGSERIAL    PRIMARY KEY,
				order_status_uuid        UUID         NOT NULL UNIQUE,
				name        VARCHAR(255) NOT NULL,
				status      SMALLINT     NOT NULL DEFAULT 1,
				created_at  TIMESTAMPTZ  DEFAULT NOW(),
				updated_at  TIMESTAMPTZ  DEFAULT NOW()
			);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_status")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_status CASCADE`)
		return err
	})
}
