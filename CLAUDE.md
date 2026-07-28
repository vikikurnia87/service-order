# service-order

Work order / inspeksi terjadwal. Detail: [README.md](README.md) dan
[docs/service-order-plan.md](docs/service-order-plan.md).

- REST `:6006` · gRPC `:6106` ⏳ · DB `service_order`
- Module: `github.com/vikikurnia87/service-order`

> **Service ini masih dalam pembangunan.** gRPC server belum ada — port `6106`
> baru dialokasikan di `.env`, belum dibaca kode. Saat membangunnya, ikuti pola
> `service-vendor` (`grpcutil.NewServer` di `app/app.go`).

## Ketergantungan

| Tujuan | Alamat | Status |
|---|---|---|
| service-user | `127.0.0.1:6101` | ✅ jalan — `clients/userclient.go`, untuk `ValidateToken` |
| service-procedure | `127.0.0.1:6105` | ⏳ **direncanakan** — `PROCEDURE_GRPC_ADDR` sudah di `.env`, klien belum dibuat |
| service-vendor | `127.0.0.1:6104` | ⏳ **direncanakan** — tabel `order_vendor` sudah ada, klien belum |

Snapshot procedure adalah rencana yang belum diimplementasi. Jangan menulis kode
yang mengasumsikannya sudah tersedia — `service-procedure` pun belum serve gRPC.

## Tabel yang dimiliki

`order`, `order_assign`, `order_category`, `order_comment`, `order_priority`,
`order_status`, `order_vendor`, `category`, `schedule`, `date`, `day`,
`master_date`

## Hal yang mudah salah

- **Order menyalin (snapshot) procedure**, tidak mereferensikannya. Perubahan
  procedure setelah order dibuat **tidak** boleh mengubah order yang berjalan.
  Ini disengaja — jangan "diperbaiki" jadi relasi hidup.
- **Order dijadwalkan berulang** lalu **diarsip** saat periodenya habis. Perhatikan
  rencana partitioning arsip di `docs/service-order-plan.md`.
- Order bisa ditugaskan ke tim internal maupun vendor (`order_vendor`).
- Seeder referensi wajib jalan sebelum dipakai: status, priority, schedule, day,
  date, master_date.
