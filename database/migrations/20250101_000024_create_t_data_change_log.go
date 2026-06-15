package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_data_change_log")
		_, err := db.ExecContext(ctx, `
			-- Audit append-only perubahan data lintas model (lihat docs: audit trail).
			CREATE TABLE IF NOT EXISTS t_data_change_log (
				id           BIGSERIAL    PRIMARY KEY,
				company_uuid UUID         NOT NULL,
				service      VARCHAR(100) NOT NULL,
				model        VARCHAR(100) NOT NULL,
				action       VARCHAR(100) NOT NULL,
				record_id    BIGINT,
				before_data  TEXT,
				after_data   TEXT,
				created_by   UUID,
				created_at   TIMESTAMPTZ  DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_t_dcl_company ON t_data_change_log(company_uuid);
			CREATE INDEX IF NOT EXISTS idx_t_dcl_record ON t_data_change_log(model, record_id);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_data_change_log")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_data_change_log CASCADE`)
		return err
	})
}
