package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  creating table: t_order")
		_, err := db.ExecContext(ctx, `
			-- t_order DIPARTISI RANGE per period_date (arsip per tahun, lihat docs).
			-- PK & UNIQUE wajib memuat partition key. Anak yang merujuk order pakai
			-- order_id BIGINT TANPA FK (FK ke tabel terpartisi rumit).
			CREATE TABLE IF NOT EXISTS t_order (
				id                    BIGINT       GENERATED ALWAYS AS IDENTITY,
				order_uuid            UUID         NOT NULL,
				company_uuid          UUID         NOT NULL,
				order_number          VARCHAR(100),
				name                  VARCHAR(255) NOT NULL,
				description           TEXT,
				order_procedure_id    BIGINT       REFERENCES t_order_procedure(id),
				order_schedule_id     BIGINT       REFERENCES t_order_schedule(id),
				order_priority_id     BIGINT       NOT NULL REFERENCES t_order_priority(id),
				order_status_id       BIGINT       NOT NULL REFERENCES t_order_status(id),
				asset_uuid            UUID,
				location_uuid         UUID,
				procedure_uuid_origin UUID,
				start_date            TIMESTAMPTZ,
				due_date              TIMESTAMPTZ,
				complete_date         TIMESTAMPTZ,
				period_date           DATE         NOT NULL,
				is_order              BOOLEAN      NOT NULL DEFAULT TRUE,
				is_system             BOOLEAN      NOT NULL DEFAULT FALSE,
				status                SMALLINT     NOT NULL DEFAULT 1,
				created_at            TIMESTAMPTZ  DEFAULT NOW(),
				updated_at            TIMESTAMPTZ  DEFAULT NOW(),
				created_by            UUID,
				updated_by            UUID,
				PRIMARY KEY (id, period_date),
				UNIQUE (order_uuid, period_date)
			) PARTITION BY RANGE (period_date);
			CREATE INDEX IF NOT EXISTS idx_t_order_company ON t_order(company_uuid);
			CREATE INDEX IF NOT EXISTS idx_t_order_status ON t_order(order_status_id);
			CREATE TABLE IF NOT EXISTS t_order_default PARTITION OF t_order DEFAULT;
			DO $$
			DECLARE y int;
			BEGIN
				FOR y IN (EXTRACT(YEAR FROM now())::int - 1)..(EXTRACT(YEAR FROM now())::int + 2) LOOP
					EXECUTE format('CREATE TABLE IF NOT EXISTS t_order_y%s PARTITION OF t_order FOR VALUES FROM (%L) TO (%L)', y, make_date(y,1,1), make_date(y+1,1,1));
				END LOOP;
			END $$;`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("  dropping table: t_order")
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS t_order CASCADE`)
		return err
	})
}
