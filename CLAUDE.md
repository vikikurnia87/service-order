# service-order

Work order / inspeksi terjadwal. Detail: [README.md](README.md) dan
[docs/service-order-plan.md](docs/service-order-plan.md).

- REST `:6006` · gRPC `:61006` · DB `service_order`
- Module: `github.com/vikikurnia87/service-order`

> **Catatan:** `GRPC_PORT=61006` menyimpang dari pola fleet (`6104`, `6105` →
> seharusnya `6106`). Kemungkinan salah ketik. Jangan diubah diam-diam — ada
> kemungkinan sudah dipakai di konfigurasi lain.

## Ketergantungan

| Tujuan | Alamat | Untuk |
|---|---|---|
| service-user | `127.0.0.1:6101` | `ValidateToken` |
| service-procedure | `127.0.0.1:6105` | ambil procedure saat pembuatan order |

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
