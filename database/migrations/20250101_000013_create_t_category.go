package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_category")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_category (
				id            BIGSERIAL    PRIMARY KEY,
				category_uuid UUID         NOT NULL UNIQUE,
				name          VARCHAR(255) NOT NULL,
				description   TEXT,
				company_uuid  UUID         NOT NULL,
				status        SMALLINT     NOT NULL DEFAULT 1,
				created_at    TIMESTAMPTZ  DEFAULT NOW(),
				updated_at    TIMESTAMPTZ  DEFAULT NOW(),
				created_by    UUID,
				updated_by    UUID
			);
			CREATE INDEX IF NOT EXISTS idx_t_category_company ON t_category(company_uuid);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_category")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_category CASCADE`)
		return err
	})
}
