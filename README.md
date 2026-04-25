# Boilerplate Clean Architecture + DDD (Go)

Dokumentasi singkat struktur folder pada boilerplate ini.

## Struktur Folder dan Deskripsi

```text
.
├── cli/                           # Entry point untuk command-line interface
├── cmd/                           # Entry point aplikasi yang bisa dieksekusi
│   └── server/                    # Bootstrap server (main app runner)
├── config/                        # Definisi dan loader konfigurasi aplikasi
├── deploy/                        # Kebutuhan deployment
│   └── docker/                    # Dockerfile dan docker-compose untuk local/deploy
├── internal/                      # Kode inti aplikasi (private to this module)
│   ├── app/                       # Wiring aplikasi, routing, dan factory
│   │   ├── factory/               # Inisialisasi dependency/komponen aplikasi
│   │   └── routes/                # Registrasi route endpoint
│   ├── delivery/                  # Layer delivery (adapter masuk: HTTP/GRPC)
│   │   ├── grpc/                  # Handler GRPC
│   │   └── http/                  # Handler HTTP
│   ├── domain/                    # Domain model (entity, contract/repository interface)
│   │   └── deposit/               # Bounded context/aggregate deposit
│   ├── dto/                       # Data Transfer Object antar layer
│   ├── infrastructure/            # Adapter keluar (DB, cache, broker, external service)
│   │   ├── cache/                 # Implementasi akses cache
│   │   ├── db/                    # Koneksi dan utilitas database
│   │   ├── external/              # Integrasi API/service eksternal
│   │   ├── kafka/                 # Producer/consumer Kafka
│   │   ├── queue/                 # Message queue implementation
│   │   └── repository/            # Implementasi repository untuk domain
│   ├── middleware/                # Middleware aplikasi (contoh: JWT)
│   └── usecase/                   # Business use case / application service
├── migrations/                    # File migrasi database
├── pkg/                           # Shared package reusable lintas modul
│   ├── graceful/                  # Graceful shutdown helper
│   ├── helper/                    # Utility/helper umum
│   └── response/                  # Standarisasi response payload
├── test/                          # Folder test (integration/e2e/supporting tests)
├── go.mod                         # Modul dan dependency Go
├── Makefile                       # Kumpulan perintah automation (build/run/test)
└── README.md                      # Dokumentasi proyek
```

## Ringkasan Layer

- `domain`: inti bisnis, bebas dari framework/teknologi eksternal.
- `usecase`: orkestrasi alur bisnis berdasarkan domain.
- `delivery`: menerima request dari luar (HTTP/GRPC) dan memanggil use case.
- `infrastructure`: implementasi teknis untuk DB, cache, broker, dan service eksternal.
- `app/cmd/cli`: inisialisasi aplikasi dan entry point proses.

Struktur ini membantu menjaga separation of concerns, memudahkan testing, dan membuat kode lebih mudah dikembangkan.