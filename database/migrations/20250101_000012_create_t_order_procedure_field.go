package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order_procedure_field")
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS t_order_procedure_field (
				id                         BIGSERIAL    PRIMARY KEY,
				order_procedure_field_uuid UUID         NOT NULL UNIQUE,
				company_uuid               UUID         NOT NULL,
				order_procedure_id         BIGINT       NOT NULL REFERENCES t_order_procedure(id) ON DELETE CASCADE,
				order_field_id             BIGINT       NOT NULL REFERENCES t_order_field(id) ON DELETE CASCADE,
				fieldorder                 SMALLINT     NOT NULL DEFAULT 0,
				value                      TEXT,                            -- jawaban/isian inspeksi
				status                     SMALLINT     NOT NULL DEFAULT 1,
				created_at                 TIMESTAMPTZ  DEFAULT NOW(),
				updated_at                 TIMESTAMPTZ  DEFAULT NOW(),
				created_by                 UUID,
				updated_by                 UUID
			);
			CREATE INDEX IF NOT EXISTS idx_t_opf_procedure ON t_order_procedure_field(order_procedure_id);
			CREATE INDEX IF NOT EXISTS idx_t_opf_field ON t_order_procedure_field(order_field_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order_procedure_field")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order_procedure_field CASCADE`)
		return err
	})
}
