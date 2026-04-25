# Boilerplate Clean Architecture + DDD (Go)

Boilerplate opinionated untuk membangun sistem Go yang **maintainable, testable, dan scalable** menggunakan prinsip **Clean Architecture** dan **Domain-Driven Design (DDD)**.

Cocok digunakan sebagai titik awal proyek baru atau sebagai referensi belajar arsitektur Go yang terstruktur.

## Tech Stack

| Komponen | Library |
|---|---|
| HTTP Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| ORM | [GORM](https://gorm.io) + PostgreSQL |
| Cache | [go-redis v9](https://github.com/redis/go-redis) |
| Queue (async job) | [Asynq](https://github.com/hibiken/asynq) via Redis |
| JWT Auth | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Config | [Viper](https://github.com/spf13/viper) (`.config.yml`) |
| Hot-reload (dev) | [Air](https://github.com/air-verse/air) |
| Process manager | [Overseer](https://github.com/jpillora/overseer) (zero-downtime restart) |
| Migration | [golang-migrate](https://github.com/golang-migrate/migrate) |

---

## Struktur Folder

```text
.
├── cli/                           # Entry point CLI (perintah non-HTTP seperti worker, seed)
├── cmd/
│   └── server/
│       └── main.go                # Entry point utama: init config → DB → Redis → Asynq → App
├── config/                        # Struct dan loader konfigurasi (via Viper)
│   ├── app.go                     # Konfigurasi aplikasi (port, env, dll)
│   ├── database.go                # Konfigurasi PostgreSQL + connection pooling
│   ├── redis.go                   # Konfigurasi Redis (cache)
│   ├── queue.go                   # Konfigurasi Asynq queue (Redis-backed)
│   └── config.go                  # Aggregator semua config loader
├── deploy/
│   └── docker/
│       ├── Dockerfile             # Multi-stage: development (Air) + production (static binary)
│       └── docker-compose.yml     # PostgreSQL + Redis + Asynqmon + App
├── internal/                      # Kode inti — tidak boleh diimport dari luar module
│   ├── app/
│   │   ├── app.go                 # Rakit Fiber + Asynq worker, daftarkan routes
│   │   ├── factory/
│   │   │   └── factory.go         # Dependency injection: wire Infrastructure → UseCase → Handler
│   │   └── routes/
│   │       ├── routes.go          # Root router: public dan protected (JWT) group
│   │       ├── healtz.go          # Route GET /api/healthz
│   │       └── deposit.go         # Route /api/deposits (POST, GET /:id, GET /user/:id)
│   ├── delivery/
│   │   └── http/
│   │       └── deposit.go         # HTTP handler: parse request → call use case → write response
│   ├── domain/
│   │   └── deposit/
│   │       ├── entity.go          # Aggregate root, business rules, factory function
│   │       └── repository.go      # Interface kontrak persistence (tidak tahu GORM/SQL)
│   ├── dto/
│   │   └── deposit.go             # Request/Response struct antar layer
│   ├── infrastructure/
│   │   ├── cache/
│   │   │   ├── redis.go           # Buat Redis client, ping test koneksi
│   │   │   └── deposit.go         # Implementasi DepositCache: Get/Set/Delete via Redis
│   │   ├── db/
│   │   │   └── db.go              # Koneksi GORM + connection pool config
│   │   ├── queue/
│   │   │   ├── client.go          # Buat Asynq client (produce) dan server (consume)
│   │   │   └── deposit.go         # Enqueue task + processor worker deposit:created
│   │   └── repository/
│   │       └── deposit.go         # Implementasi Repository: GORM, mapping model↔domain
│   ├── middleware/
│   │   └── jwt.go                 # JWT middleware + GenerateToken helper
│   └── usecase/
│       └── deposit.go             # Application service: orkestrasi domain + cache + queue
├── migrations/
│   ├── 000001_create_deposits_table.up.sql    # Buat tabel deposits, index, trigger
│   └── 000001_create_deposits_table.down.sql  # Rollback: drop tabel deposits
├── pkg/
│   ├── graceful/
│   │   └── graceful.go            # Setup signal handler untuk graceful shutdown
│   └── response/
│       └── response.go            # Standar JSON response: success, error, with-meta
├── .air.toml                      # Konfigurasi Air hot-reload (poll mode untuk Docker volume)
├── .config.yml                    # Konfigurasi aplikasi (lihat .config.exampe.yml)
├── .config.exampe.yml             # Template konfigurasi
├── .env                           # Environment variables untuk docker-compose
├── go.mod                         # Module dependencies
└── Makefile                       # Automation: docker, migrate, build, test
```

---

## Arsitektur

Proyek ini mengikuti **Clean Architecture** — dependency hanya boleh mengarah ke dalam (ke domain), tidak boleh ke luar.

```
 ┌─────────────────────────────────────────────────────┐
 │  Delivery (HTTP Handler)                            │
 │  • Parse request, panggil use case, tulis response  │
 └────────────────────┬────────────────────────────────┘
                      │ memanggil
 ┌────────────────────▼────────────────────────────────┐
 │  Use Case (Application Service)                     │
 │  • Orkestrasi domain, cache-aside, enqueue job      │
 │  • Hanya tahu interface (Repository/Cache/Queue)    │
 └──────┬─────────────┬──────────────┬─────────────────┘
        │             │              │
    domain        cache port     queue port
        │             │              │
 ┌──────▼──┐   ┌──────▼──────┐  ┌───▼────────────────┐
 │ Domain  │   │ Redis (impl)│  │ Asynq (impl)       │
 │ Entity  │   │ cache/      │  │ queue/             │
 │ + Rules │   │ deposit.go  │  │ deposit.go         │
 └──────┬──┘   └─────────────┘  └────────────────────┘
        │
 ┌──────▼──────────────┐
 │ Repository (impl)   │
 │ GORM + PostgreSQL   │
 │ repository/deposit  │
 └─────────────────────┘
```

**Aturan utama:** `domain` tidak boleh mengimport package apapun dari `infrastructure`, `delivery`, atau `usecase`. Ketergantungan hanya satu arah — ke dalam.

---

## Alur Request (Contoh: `POST /api/deposits`)

```
Client
  └─▶ Fiber Router (routes.go)
        └─▶ JWT Middleware (middleware/jwt.go)
              └─▶ DepositHandler.CreateDeposit (delivery/http/deposit.go)
                    │  parse body → dto.CreateDepositRequest
                    └─▶ depositUseCase.CreateDeposit (usecase/deposit.go)
                          │  deposit.NewDeposit()  ← validasi di domain
                          │  depositRepo.Create()  ← simpan ke PostgreSQL
                          └─▶ depositQueue.EnqueueDepositCreated()  ← kirim job ke Redis
                                └─▶ [async] ProcessDepositCreated worker
```

**Alur `GET /api/deposits/:id` — cache-aside pattern:**

```
Handler.GetDeposit
  └─▶ useCase.GetDeposit
        ├─▶ cache.Get(id)
        │     ├─ HIT  → return dari Redis (tidak ke DB)
        │     └─ MISS → repo.FindByID() → cache.Set() → return
```

---

## Cara Penggunaan

### Prasyarat

- Docker & Docker Compose
- Go 1.24+
- `make`

### 1. Setup konfigurasi

Salin template config dan sesuaikan:

```bash
cp .config.exampe.yml .config.yml
```

Nilai penting di `.config.yml` untuk development dengan Docker:

```yaml
database:
  host: postgres   # nama service di docker-compose
  name: boilerplate
  user: postgres
  password: postgres

cache:
  redis:
    host: redis    # nama service di docker-compose
    password: password

queue:
  redis:
    host: redis
    password: password
    database: 1    # pisah dari cache (database: 0)
```

### 2. Jalankan development (hot-reload)

```bash
make dev
```

Air akan otomatis rebuild dan restart binary setiap kali ada perubahan file `.go` atau `.config.yml`. Log bisa dipantau dengan:

```bash
make logs-app
```

### 3. Jalankan migration

```bash
# Install migrate CLI (sekali saja)
make migrate-install

# Jalankan semua migration
make migrate-up
```

### 4. Akses services

| Service | URL |
|---|---|
| App API | http://localhost:8000 |
| Asynqmon (queue monitor) | http://localhost:8080 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

---

## Makefile Reference

### Docker

| Perintah | Fungsi |
|---|---|
| `make dev` | Jalankan semua service (foreground) |
| `make dev-d` | Jalankan semua service (background) |
| `make build` | Rebuild image app |
| `make up` / `make down` | Start / stop semua container |
| `make down-v` | Stop + hapus volume (data DB hilang) |
| `make stop` / `make restart` | Stop / restart tanpa destroy |
| `make logs` / `make logs-app` | Tail log realtime |
| `make shell` | Masuk ke shell container app |

### Migration

| Perintah | Fungsi |
|---|---|
| `make migrate-install` | Install `migrate` CLI |
| `make migrate-create name=xxx` | Buat file migration baru |
| `make migrate-up` | Jalankan semua migration pending |
| `make migrate-down` | Rollback 1 migration terakhir |
| `make migrate-version` | Cek versi migration aktif |
| `make migrate-force version=N` | Fix dirty state |

Override DB connection: `make migrate-up DB_HOST=myhost DB_PASS=secret`

---

## Menambah Domain Baru

Ikuti pola yang sama dengan domain `deposit`:

1. **Domain** — buat folder `internal/domain/<nama>/`, isi `entity.go` dan `repository.go`
2. **DTO** — tambah file di `internal/dto/<nama>.go`
3. **Use Case** — buat `internal/usecase/<nama>.go`, definisikan interface cache/queue jika perlu
4. **Infrastructure** — implementasi repository di `internal/infrastructure/repository/<nama>.go`
5. **Handler** — buat `internal/delivery/http/<nama>.go`
6. **Routes** — buat `internal/app/routes/<nama>.go`, daftarkan di `routes.go`
7. **Factory** — wire semua layer di `internal/app/factory/factory.go`
8. **Migration** — `make migrate-create name=create_<nama>_table`

---

## Prinsip yang Diterapkan

- **Dependency Inversion** — use case hanya bergantung pada interface, bukan implementasi konkret
- **Single Responsibility** — setiap file/struct punya satu tanggung jawab
- **Domain isolation** — domain tidak mengimport framework, ORM, atau library eksternal apapun
- **Cache-aside pattern** — cek cache → miss → DB → set cache
- **Graceful shutdown** — semua resource (DB, Redis, Asynq) ditutup dengan urutan yang benar saat menerima sinyal OS