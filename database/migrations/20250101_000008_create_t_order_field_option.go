package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_field_option")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_field_option (
				id                      BIGSERIAL    PRIMARY KEY,
				order_field_option_uuid UUID         NOT NULL UNIQUE,
				company_uuid            UUID         NOT NULL,
				name                    VARCHAR(255) NOT NULL,
				order_field_id          BIGINT       NOT NULL REFERENCES t_order_field(id) ON DELETE CASCADE,
				value                   VARCHAR(255),
				status                  SMALLINT     NOT NULL DEFAULT 1,
				created_at              TIMESTAMPTZ  DEFAULT NOW(),
				updated_at              TIMESTAMPTZ  DEFAULT NOW(),
				created_by              UUID,
				updated_by              UUID
			);
			CREATE INDEX IF NOT EXISTS idx_t_order_field_option_field ON t_order_field_option(order_field_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_field_option")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_field_option CASCADE`)
		return err
	})
}
