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

## Roadmap (detail di docs/service-order-plan.md)
| Fase | Status | Isi |
|---|---|---|
| 0 Scaffold | done | kerangka service (config/app/container/middleware), tanpa minio |
| 1 Master + seeder | done | order_status/priority/schedule/day/date/master_date + seeder CSV |
| 2 Migrasi sisa tabel | done | snapshot form + order inti (partisi) + penjadwalan + audit |
| 3 CRUD order | todo | repo/service/handler + snapshot via gRPC service-procedure |
| 4 Mesin jadwal | todo | generator okurensi |
| 5 Arsip | todo | partitioning + worker archiver |
| 6 Media | todo | order_file/img -> service-media (satu pintu) |

Port: HTTP :6006, gRPC :61006.
