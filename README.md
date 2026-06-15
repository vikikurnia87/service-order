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
    cp .env.example .env        # sesuaikan DB, JWT, USER_GRPC_ADDR
    make migrate                # jalankan migrasi (Fase 1+)
    make seed                   # seed referensi (Fase 1+)
    make run                    # HTTP :6006

## Roadmap (detail di docs/service-order-plan.md)
| Fase | Status | Isi |
|---|---|---|
| 0 Scaffold | done | kerangka service (config/app/container/middleware), tanpa minio |
| 1 Master + seeder | todo | order_status/priority/schedule/day/date/master_date + seeder CSV |
| 2 Migrasi sisa tabel | todo | snapshot form + order inti (partisi) + penjadwalan |
| 3 CRUD order | todo | repo/service/handler + snapshot via gRPC service-procedure |
| 4 Mesin jadwal | todo | generator okurensi |
| 5 Arsip | todo | partitioning + worker archiver |
| 6 Media | todo | order_file/img -> service-media (satu pintu) |

Port: HTTP :6006, gRPC :61006.
