package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_field")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_field (
				id               BIGSERIAL    PRIMARY KEY,
				order_field_uuid UUID         NOT NULL UNIQUE,
				company_uuid     UUID         NOT NULL,
				name             VARCHAR(255) NOT NULL,
				field_type_id    BIGINT,                                  -- ref service-procedure (snapshot, no FK)
				isrequired       BOOLEAN      NOT NULL DEFAULT FALSE,
				isheading        BOOLEAN      NOT NULL DEFAULT FALSE,
				status           SMALLINT     NOT NULL DEFAULT 1,
				created_at       TIMESTAMPTZ  DEFAULT NOW(),
				updated_at       TIMESTAMPTZ  DEFAULT NOW(),
				created_by       UUID,
				updated_by       UUID
			);
			CREATE INDEX IF NOT EXISTS idx_t_order_field_company ON t_order_field(company_uuid);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_field")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_field CASCADE`)
		return err
	})
}
