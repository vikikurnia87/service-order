# service-order — Plan & Guide (Migration, Seeder, Arsip)

> **Status:** rencana implementasi (belum di-scaffold). Dokumen ini = acuan tetap saat membangun service-order.
> **Sumber:** project lama `../../order` (Node.js + Sequelize, 25 tabel) — diadaptasi ke konvensi AMANOS (Go + bun).
> **Prinsip arsip:** mencontoh `service-myneuops-api` (data lama dipindah keluar tabel panas), tetapi karena order **relasional & kompleks**, memakai **PostgreSQL partitioning** untuk order dan **tabel `_year` manual** hanya untuk log/append.

---

## 1. Lingkup & posisi

service-order = domain **work order / inspeksi terjadwal**. Order **menyalin (snapshot)** sebuah procedure dari **service-procedure** lalu diisi, dijadwalkan berulang, ditugaskan ke tim/vendor, dan diarsip saat periodenya habis.

**Fokus tahap awal (dokumen ini):** **migration (struktur DB)** + **seeder**, serta **strategi arsip** untuk data kedaluwarsa.

---

## 2. Konvensi AMANOS (pemetaan dari kode lama)

| Kode lama (Sequelize) | service-order (Go + bun) |
|---|---|
| Tabel PascalCase plural `OrderPriorities` | snake_case `t_order_priority` |
| dual `id` INT + `<entity>uuid`, **FK by-uuid** | `id` BIGINT PK + `<entity>_uuid` UUID; **FK internal by-id** |
| FK lintas domain (1 DB) | **referensi UUID tanpa FK** (asset/location/procedure/vendor/user) |
| `isactive` (INTEGER 0/1) | `status` SMALLINT default 1 |
| flag INTEGER (`isorder`,`issystem`) | BOOLEAN |
| `useruuid` (pembuat) | `created_by` / `updated_by` UUID + `company_uuid` |
| `HOrderLogs`/`HCategoryLogs` (audit) | **`t_data_change_log`** (1 tabel, pola service-utils/vendor) |
| Migrasi 1 file per tabel | **dipertahankan** (anti-pattern menggabung) |
| Seeder | **CSV** (pola `csv_reader.go`) |

Referensi lintas-service (simpan **UUID saja**, enrich via gRPC/event — lihat `../../docs/saas-architecture.md`):

| Tabel | Kolom | Milik service |
|---|---|---|
| `t_order` | `asset_uuid` | service-assets |
| `t_order` | `location_uuid` | (location) |
| `t_order`, `t_order_procedure` | `procedure_uuid_origin` | service-procedure |
| `t_order_vendor` | `vendor_uuid` | service-vendor |
| `t_order_assign` | `team_uuid`, `user_uuid` | service-user |

---

## 3. Migration plan (±23 tabel, 1 file/tabel)

### 3.1 Master / referensi (di-seed)
| Tabel | Kolom inti |
|---|---|
| `t_order_status` | id, order_status_uuid, name, status, ts |
| `t_order_priority` | id, order_priority_uuid, name, status, ts |
| `t_schedule` | id, schedule_uuid, name, status, ts |
| `t_day` | id, day_uuid, name, status, ts |
| `t_date` | id, date_uuid, name, value, status, ts |
| `t_master_date` | id, master_date_uuid, date(DATE), day_no, day, week, month, year, status, ts |

### 3.2 Snapshot form (order menyalin procedure terisi)
| Tabel | Kolom inti | FK internal |
|---|---|---|
| `t_order_field` | id, field_uuid, company_uuid, name, field_type_id?, isrequired, isheading, status, audit | — |
| `t_order_field_option` | id, field_option_uuid, company_uuid, name, field_id, value, status, audit | → t_order_field |
| `t_order_field_img` | id, …_uuid, field_id, name, status, audit | → t_order_field |
| `t_order_field_file` | id, …_uuid, field_id, name, status, audit | → t_order_field |
| `t_order_procedure` | id, order_procedure_uuid, **procedure_uuid_origin**, name, desc, company_uuid, status, audit | — |
| `t_order_procedure_field` | id, …_uuid, order_procedure_id, field_id, fieldorder, **value** (isian), company_uuid, status, audit | → t_order_procedure, t_order_field |

> **Catatan:** `field_type_id` boleh referensi **lokal** (jika service-order menyimpan salinan jenis field) atau hanya menyimpan label. Karena snapshot, idealnya order menyimpan cukup info agar tak bergantung service-procedure setelah dibuat.

### 3.3 Penjadwalan (recurring)
| Tabel | Kolom inti | FK internal |
|---|---|---|
| `t_order_schedule` | id, order_schedule_uuid, schedule_id, date_id, repeat_every, status, audit | → t_schedule, t_date |
| `t_order_schedule_day` | id, …_uuid, order_schedule_id, day_id, status, audit | → t_order_schedule, t_day |
| `t_transaction_schedule` | id, transaction_schedule_uuid, order_id, start_date, end_date, master_date_id, schedule_id, status, ts | → t_order, t_master_date, t_schedule |

### 3.4 Order inti + lampiran
| Tabel | Kolom inti | FK internal | Lintas-service |
|---|---|---|---|
| `t_order` | id, order_uuid, order_number, name, desc, order_procedure_id?, start_date, due_date, complete_date, order_schedule_id?, order_priority_id, order_status_id, is_order, is_system, **period_date**, company_uuid, status, audit | → priority, status | asset_uuid, location_uuid, procedure_uuid_origin, user(created_by) |
| `t_category` | id, category_uuid, name, desc, company_uuid, status, audit | — | — |
| `t_order_category` | id, …_uuid, order_id, category_id, status, audit | → t_order, t_category | — |
| `t_order_assign` | id, order_assign_uuid, order_id, status, audit | → t_order | team_uuid, user_uuid |
| `t_order_vendor` | id, order_vendor_uuid, order_id, status, audit | → t_order | vendor_uuid |
| `t_order_comment` | id, order_comment_uuid, order_id, name(text), status, audit | → t_order | user (created_by) |
| `t_order_file` | id, order_file_uuid, order_id, name, status, audit | → t_order | (file → service-media) |
| `t_order_img` | id, order_img_uuid, order_id, name, status, audit | → t_order | (img → service-media) |

> **Audit:** `HOrderLogs`/`HCategoryLogs` lama **dilebur** ke `t_data_change_log` (pola AMANOS).
> **Media:** `t_order_file`/`t_order_img` hanya simpan metadata; upload aktual ke **service-media** (pola "satu pintu" seperti vendor) — *menyusul, bukan tahap migrasi*.
> **`period_date`** di `t_order` = kunci partisi arsip (lihat §5).

### 3.5 Urutan migrasi (dependensi FK)
```
order_status, order_priority, schedule, day, date, master_date
→ order_field → order_field_option / order_field_img / order_field_file
→ order_procedure → order_procedure_field
→ order_schedule → order_schedule_day
→ category
→ order (partisi, §5)
→ order_category, order_assign, order_vendor, order_comment, order_file, order_img
→ transaction_schedule
→ data_change_log
```

---

## 4. Seeder plan (CSV — pola `csv_reader.go`)

| Seeder | Sumber | Isi |
|---|---|---|
| `OrderStatusSeeder` | `order_status.csv` | Open, On Hold, In Progress, Done |
| `OrderPrioritySeeder` | `order_priority.csv` | None, Low, Medium, High |
| `ScheduleSeeder` | `schedule.csv` | None, Daily, Weekly, Monthly, Yearly |
| `DaySeeder` | `day.csv` | sunday … saturday |
| `DateSeeder` | `date.csv` | 1st…31st (value 01…31) |
| `MasterDateSeeder` | **generator** | kalender per tahun (programatik) |

- File CSV asli ada di `../../order/src/seeders/csvs/` (acuan nilai + bila ingin uuid deterministik).
- Tiap entitas yang punya `<entity>_uuid` → seeder set `uuid.New()` (atau pakai uuid dari CSV agar deterministik lintas rebuild).
- **MasterDate**: JANGAN CSV raksasa. Buat **generator** Go: loop `time.Time` dari `MASTERDATE_START_YEAR` s/d `END_YEAR`, isi `day_no` (1=Sun…7=Sat atau ISO), `day`, `week` (ISO week), `month`, `year`. Idempotent per tahun.
- Urutan runner: master/referensi dulu, lalu master_date.

---

## 5. ⭐ Arsip data kedaluwarsa — PostgreSQL Partitioning

### 5.1 Kenapa partitioning (bukan `_year` manual ala neuops) untuk order
neuops memakai tabel fisik per-tahun `report_checklist_detail_<year>` (nama dibangun `fmt.Sprintf("..%d", year)`, akses via `ModelTableExpr`) — **cocok untuk data append-only/leaf**. Order = **agregat relasional** dengan kebutuhan query lintas-tahun → `_year` manual memaksa nama tabel dinamis + UNION. **Partitioning** mencapai tujuan yang sama (data lama keluar dari tabel panas) tapi **query tetap transparan** (Postgres me-route).

### 5.2 Struktur partisi (RANGE per tahun)
```sql
CREATE TABLE t_order (
    id              BIGINT GENERATED ALWAYS AS IDENTITY,
    order_uuid      UUID  NOT NULL,
    company_uuid    UUID  NOT NULL,
    period_date     DATE  NOT NULL,                 -- partition key (mis. due/periode)
    order_status_id BIGINT NOT NULL,
    ...
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, period_date),                  -- WAJIB sertakan partition key
    UNIQUE (order_uuid, period_date)                -- uuid unik komposit (lihat gotcha)
) PARTITION BY RANGE (period_date);

CREATE TABLE t_order_y2025 PARTITION OF t_order FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
CREATE TABLE t_order_y2026 PARTITION OF t_order FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
CREATE TABLE t_order_default PARTITION OF t_order DEFAULT;        -- penampung, jaga tetap kecil

CREATE INDEX idx_t_order_company ON t_order (company_uuid);       -- index di induk → turun ke partisi
```

### 5.3 Arsip = DETACH partisi
```sql
ALTER TABLE t_order DETACH PARTITION t_order_y2023 CONCURRENTLY;  -- PG14+: tanpa lock berat
```
`t_order_y2023` jadi tabel mandiri → `pg_dump` ke storage murah / pindah tablespace / `DROP` (retensi). Tabel panas mengecil, query cepat — tanpa ubah kode aplikasi.

### 5.4 Gotchas WAJIB
- **PK/unique harus memuat partition key** → PK `(id, period_date)`; `order_uuid` unik jadi `(order_uuid, period_date)`. Keunikan murni `order_uuid` tak dijamin DB → andalkan `uuid.New()` (+ unique komposit). `id` tetap unik global (sequence bersama).
- **FK ke tabel terpartisi** perlu **PG 12+**; lebih aman: anak bervolume tinggi **ikut dipartisi** key tahun yang sama (denormalisasi `period_date`).
- **Partisi tidak dibuat otomatis** → pre-create via scheduler tahunan atau ekstensi **`pg_partman`**; sediakan `DEFAULT` kecil agar INSERT di luar rentang tak error.
- **bun tidak perlu tahu partisi** — model `t_order` tetap satu; insert/select ke induk. DDL partisi ditulis raw SQL di migrasi.
- **Jangan DETACH** partisi tahun yang masih punya order `Open` → archiver cek "semua Done" dulu atau pindahkan straggler.
- Syarat: **PostgreSQL ≥ 13** (idealnya **14+** untuk `DETACH … CONCURRENTLY`).

### 5.5 Anak (agregat) yang ikut dipartisi
- **Partisi**: `t_order` + `t_order_procedure_field` (paling banyak baris — nilai isian) dengan `period_date` yang sama → DETACH satu tahun melepas agregat bersama.
- **Normal** (volume rendah): `t_order_assign`, `t_order_comment`, `t_order_category`, dst. (boleh dipartisi nanti bila perlu detach total).

### 5.6 `_year` manual hanya untuk log/append
Tabel **murni append/log** (mis. `t_data_change_log`) boleh pakai pola neuops `t_data_change_log_<year>` (nama dinamis via int year) — lebih simpel, tak butuh query relasional.

---

## 6. Mesin penjadwalan (catatan, fase lanjut)

Dua job (scheduler/worker, terpisah):
1. **Generator** — dari `t_order_schedule` (+ `t_order_schedule_day`) × `t_master_date` × `t_schedule` (Daily/Weekly/Monthly/Yearly) → buat okurensi `t_transaction_schedule` + order mendatang.
2. **Archiver** — cari order **Done + periode berakhir + tahun lampau**, pastikan partisi tahun siap, lalu `DETACH`/arsip (§5.3).

"Periode habis" = `t_transaction_schedule.end_date` < now, atau `t_order.due_date`/`period_date` berlalu.

---

## 7. Roadmap bertahap

| Fase | Isi | Status |
|---|---|---|
| 0 — Scaffold | infra (config/app/container/middleware) pola procedure; media menyusul | ⬜ |
| **1 — Master + seeder** | `t_order_status/priority/schedule/day/date/master_date` + 6 seeder (CSV + MasterDate generator) | ⬜ ← mulai sini |
| **2 — Migrasi sisa tabel** | snapshot form + order inti (partisi) + penjadwalan (struktur) | ⬜ |
| 3 — CRUD order | repo/service/handler; snapshot procedure via gRPC ke service-procedure | ⬜ |
| 4 — Mesin jadwal | generator okurensi | ⬜ |
| 5 — Arsip | partitioning + worker archiver (+ `_year` utk log) | ⬜ |
| 6 — Media | `t_order_file/img` → service-media (pola satu pintu vendor) | ⬜ |

---

## 8. Keputusan & catatan terbuka
1. **Arsip**: partitioning untuk `t_order`(+`t_order_procedure_field`); `_year` manual untuk `t_data_change_log`. → *default dokumen ini*.
2. **MasterDate**: generator programatik (bukan CSV besar).
3. **Audit**: lebur ke `t_data_change_log`.
4. **Partition key**: `period_date` (selaras makna "periode habis"), bukan `created_at`.
5. **Versi Postgres**: pastikan ≥ 13/14 sebelum mengandalkan DETACH CONCURRENTLY.
6. **`field_type` di order**: simpan label/salinan agar snapshot mandiri (tak bergantung service-procedure pasca-create).

---

## 9. Referensi
- Kode lama: `../../order/src/migrations`, `../../order/src/seeders` (Sequelize).
- Pola `_year` neuops: `../../service-myneuops-api/repositories/report_checklist_detail_repository.go`.
- Konvensi service AMANOS: `../../service-procedure` (struktur), `../../service-vendor` (CRUD+media+audit).
- Arsitektur SaaS & master data: `../../docs/saas-architecture.md`.
- Seeder CSV helper: `../../service-procedure/database/seeders/csv_reader.go`.
