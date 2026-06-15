package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_category")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_category (
				id                  BIGSERIAL   PRIMARY KEY,
				order_category_uuid UUID        NOT NULL UNIQUE,
				order_id            BIGINT      NOT NULL,
				category_id         BIGINT      NOT NULL REFERENCES t_category(id) ON DELETE CASCADE,
				status              SMALLINT    NOT NULL DEFAULT 1,
				created_at          TIMESTAMPTZ DEFAULT NOW(),
				updated_at          TIMESTAMPTZ DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_t_order_category_order ON t_order_category(order_id);
			CREATE INDEX IF NOT EXISTS idx_t_order_category_cat ON t_order_category(category_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_category")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_category CASCADE`)
		return err
	})
}
