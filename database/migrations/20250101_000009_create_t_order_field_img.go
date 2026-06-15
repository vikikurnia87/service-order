package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_field_img")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_field_img (
				id                   BIGSERIAL    PRIMARY KEY,
				order_field_img_uuid UUID         NOT NULL UNIQUE,
				order_field_id       BIGINT       NOT NULL REFERENCES t_order_field(id) ON DELETE CASCADE,
				name                 VARCHAR(255) NOT NULL,
				status               SMALLINT     NOT NULL DEFAULT 1,
				created_at           TIMESTAMPTZ  DEFAULT NOW(),
				updated_at           TIMESTAMPTZ  DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_t_order_field_img_field ON t_order_field_img(order_field_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_field_img")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_field_img CASCADE`)
		return err
	})
}
