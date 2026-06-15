# service-order

Microservice **work order / inspeksi terjadwal** (AMANOS). Order **menyalin (snapshot)**
procedure dari service-procedure, dijadwalkan berulang, ditugaskan ke tim/vendor, lalu
**diarsip** saat periodenya habis.

> Rencana lengkap (migration, seeder, arsip partitioning): **[docs/service-order-plan.md](docs/service-order-plan.md)**

## Struktur
Pola standar AMANOS (sama vendor/procedure): configs · container · app · auth · clients ·
middlewares · database(migrations,seeders) · models · structs · repositories · services ·
mappers · handlers · routes. Auth via gRPC ke service-user; ELK/APM; DI; status 400/422/404.

## Menjalankan
    cp .env.example .env        # sesuaikan DB, USER_GRPC_ADDR, MASTERDATE_*
    make migrate                # jalankan migrasi pending (24 tabel)
    make seed                   # seed referensi (status/priority/schedule/day/date/master_date)
    make run                    # HTTP :6006

### Migrasi & seeder (detail)
    make migrate                # jalankan migrasi pending
    make status                 # cek status tiap migrasi
    make rollback               # rollback batch terakhir
    make reset                  # rollback semua
    make fresh                  # drop semua tabel + migrate ulang
    make seed                   # semua seeder referensi (idempotent)
    make seed NAME=DaySeeder    # satu seeder saja
    make fresh-seed             # fresh + seed sekaligus

Seeder referensi (urut): `OrderStatusSeeder` · `OrderPrioritySeeder` · `ScheduleSeeder` ·
`DaySeeder` · `DateSeeder` · `MasterDateSeeder`. Lima pertama dari CSV
(`database/seeders/file/`, idempotent via cek nama); `MasterDateSeeder` meng-generate
kalender programatik dari `MASTERDATE_START_YEAR`..`MASTERDATE_END_YEAR` (default tahun
berjalan s/d +2), idempotent via `ON CONFLICT (full_date) DO NOTHING`.

> Tanpa GNU make, pakai langsung: `go run ./cmd/migrate <migrate|rollback|reset|status|fresh|seed [NAME]|fresh:seed>`.

## API (semua butuh Bearer token + company aktif)
| Method | Path | Keterangan |
|---|---|---|
| GET | `/api/v1/order` | list order (paginated `?page&limit`, `?search` nama/nomor) |
| GET | `/api/v1/order/:id` | detail order + anak (kategori/vendor/assignee/komentar) |
| POST | `/api/v1/order` | buat order |

POST body (inti + anak; `priority_uuid`/`status_uuid` wajib, sisanya opsional):
```json
{
  "name": "Inspeksi AC Lt.3",
  "description": "rutin",
  "order_number": "WO-001",
  "priority_uuid": "<order_priority_uuid>",
  "status_uuid": "<order_status_uuid>",
  "asset_uuid": "<uuid>", "location_uuid": "<uuid>",
  "start_date": "2026-06-20T08:00:00Z", "due_date": "2026-06-20T17:00:00Z",
  "period_date": "2026-06-20",
  "categories": ["<category_uuid>"],
  "vendors": ["<vendor_uuid>"],
  "assignees": [{ "team_uuid": "<uuid>", "user_uuid": "<uuid>" }],
  "comments": ["catatan awal"]
}
```
> `period_date` (partition key) default dari `start_date`/hari ini bila kosong. Snapshot
> prosedur (`procedure_uuid`) & penjadwalan berulang menyusul di Fase 3/4.

## Roadmap (detail di docs/service-order-plan.md)
| Fase | Status | Isi |
|---|---|---|
| 0 Scaffold | done | kerangka service (config/app/container/middleware), tanpa minio |
| 1 Master + seeder | done | order_status/priority/schedule/day/date/master_date + seeder CSV |
| 2 Migrasi sisa tabel | done | snapshot form + order inti (partisi) + penjadwalan + audit |
| 3 CRUD order | wip | GET (list/detail) + POST create (inti + kategori/vendor/assignee/komentar) done; PUT/DELETE + snapshot prosedur via gRPC service-procedure menyusul |
| 4 Mesin jadwal | todo | generator okurensi |
| 5 Arsip | todo | partitioning + worker archiver |
| 6 Media | todo | order_file/img -> service-media (satu pintu) |

Port: HTTP :6006, gRPC :61006.
