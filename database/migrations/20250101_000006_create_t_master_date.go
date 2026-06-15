package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_master_date")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_master_date (
				id               BIGSERIAL    PRIMARY KEY,
				master_date_uuid UUID         NOT NULL UNIQUE,
				full_date        DATE         NOT NULL UNIQUE,
				day_no           SMALLINT     NOT NULL,
				day              VARCHAR(20)  NOT NULL,
				day_of_month     SMALLINT     NOT NULL,
				week             SMALLINT     NOT NULL,
				month            SMALLINT     NOT NULL,
				year             SMALLINT     NOT NULL,
				status           SMALLINT     NOT NULL DEFAULT 1,
				created_at       TIMESTAMPTZ  DEFAULT NOW(),
				updated_at       TIMESTAMPTZ  DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_t_master_date_ym ON t_master_date(year, month);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_master_date")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_master_date CASCADE`)
		return err
	})
}
