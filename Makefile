# Makefile — service-order
# Catatan Windows: butuh GNU make (mis. via choco install make / scoop install make).

APP      := service-order
MIGRATE  := go run ./cmd/migrate
BIN_DIR  := bin

.PHONY: help
help: ## Tampilkan daftar perintah
	@echo Targets:
	@echo   make migrate        - jalankan migrasi pending
	@echo   make rollback       - rollback batch terakhir
	@echo   make reset          - rollback semua migrasi
	@echo   make status         - status tiap migrasi
	@echo   make fresh          - drop semua tabel + migrate ulang
	@echo   make seed           - jalankan semua seeder (field_type + template)
	@echo   make seed NAME=X    - jalankan satu seeder (mis. NAME=FieldTypeSeeder)
	@echo   make fresh-seed     - fresh + seed sekaligus
	@echo   make run            - jalankan server lokal (REST :6006 + gRPC :61006)
	@echo   make air            - jalankan server dengan hot-reload (air)
	@echo   make build          - compile binary ke bin/
	@echo   make tidy           - go mod tidy
	@echo   make install-air    - install tool air

# ─── Migrations & Seed ───────────────────────────────────────
.PHONY: migrate rollback reset status fresh seed fresh-seed
migrate: ## Jalankan migrasi pending
	$(MIGRATE) migrate

rollback: ## Rollback batch terakhir
	$(MIGRATE) rollback

reset: ## Rollback semua migrasi
	$(MIGRATE) reset

status: ## Status tiap migrasi
	$(MIGRATE) status

fresh: ## Drop semua tabel + migrate ulang
	$(MIGRATE) fresh

seed: ## Jalankan semua seeder (atau satu: make seed NAME=FieldTypeSeeder)
	$(MIGRATE) seed $(NAME)

fresh-seed: ## fresh + seed sekaligus
	$(MIGRATE) fresh:seed

# ─── Run / Build ─────────────────────────────────────────────
.PHONY: run air build tidy install-air
run: ## Jalankan server lokal (REST + gRPC)
	go run ./cmd

air: ## Hot-reload via air
	air

build: ## Compile binary
	go build -o $(BIN_DIR)/$(APP) ./cmd

tidy: ## go mod tidy
	go mod tidy

install-air: ## Install air
	go install github.com/air-verse/air@latest
