package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_schedule")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_schedule (
				id          BIGSERIAL    PRIMARY KEY,
				schedule_uuid        UUID         NOT NULL UNIQUE,
				name        VARCHAR(255) NOT NULL,
				status      SMALLINT     NOT NULL DEFAULT 1,
				created_at  TIMESTAMPTZ  DEFAULT NOW(),
				updated_at  TIMESTAMPTZ  DEFAULT NOW()
			);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_schedule")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_schedule CASCADE`)
		return err
	})
}
