// cmd/migrate adalah CLI migrasi skema service-vendor.
//
// Usage:
//
//	go run ./cmd/migrate migrate    # jalankan migrasi
//	go run ./cmd/migrate rollback   # rollback 1 batch
//	go run ./cmd/migrate reset      # rollback semua
//	go run ./cmd/migrate status     # status
//	go run ./cmd/migrate fresh      # reset + migrate
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/vikikurnia87/service-order/configs"
	"github.com/vikikurnia87/service-order/database"
	"github.com/vikikurnia87/service-order/database/migrations"
	"github.com/vikikurnia87/service-order/database/seeders"

	sudb "github.com/vikikurnia87/service-utils/database"
)

func main() {
	configs.LoadEnv()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	db := database.InitPostgresDatabase(buildDBConfig(), configs.ServerEnv, logger)
	defer func() { _ = database.Close(db) }()

	ctx := context.Background()
	migrator := migrations.NewMigrator(db)
	seeder := seeders.NewDatabaseSeeder(db)

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "migrate":
		must(migrator.Migrate(ctx))
	case "rollback":
		must(migrator.Rollback(ctx))
	case "reset":
		must(migrator.Reset(ctx))
	case "status":
		must(migrator.Status(ctx))
	case "fresh":
		must(migrator.Fresh(ctx))
	case "seed":
		if len(os.Args) > 2 {
			must(seeder.RunOne(ctx, os.Args[2]))
		} else {
			must(seeder.RunAll(ctx))
		}
	case "fresh:seed":
		must(migrator.Fresh(ctx))
		must(seeder.RunAll(ctx))
	default:
		log.Fatalf("unknown command: %q (migrate|rollback|reset|status|fresh|seed|fresh:seed)", cmd)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func buildDBConfig() sudb.Config {
	return sudb.Config{
		Host:            configs.DatabaseHost,
		Port:            configs.DatabasePort,
		User:            configs.DatabaseUser,
		Password:        configs.DatabasePassword,
		Name:            configs.DatabaseName,
		Timeout:         configs.DatabaseTimeout,
		DialTimeout:     configs.DatabaseDialTimeout,
		ReadTimeout:     configs.DatabaseReadTimeout,
		WriteTimeout:    configs.DatabaseWriteTimeout,
		SSLMode:         configs.DatabaseSSLMode,
		MaxOpenConns:    configs.DatabaseMaxConn,
		MaxIdleConns:    configs.DatabaseMaxIdleConn,
		ConnMaxLifetime: time.Duration(configs.DatabaseMaxLifeConn) * time.Minute,
	}
}
