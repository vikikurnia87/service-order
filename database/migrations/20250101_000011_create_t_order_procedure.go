package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_procedure")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_procedure (
				id                    BIGSERIAL    PRIMARY KEY,
				order_procedure_uuid  UUID         NOT NULL UNIQUE,
				procedure_uuid_origin UUID,                                -- ref service-procedure (no FK)
				name                  VARCHAR(255) NOT NULL,
				description           TEXT,
				company_uuid          UUID         NOT NULL,
				status                SMALLINT     NOT NULL DEFAULT 1,
				created_at            TIMESTAMPTZ  DEFAULT NOW(),
				updated_at            TIMESTAMPTZ  DEFAULT NOW(),
				created_by            UUID,
				updated_by            UUID
			);
			CREATE INDEX IF NOT EXISTS idx_t_order_procedure_company ON t_order_procedure(company_uuid);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_procedure")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_procedure CASCADE`)
		return err
	})
}
