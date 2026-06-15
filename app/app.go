// Package app adalah orchestrator lifecycle service-order.
package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/vikikurnia87/service-order/clients"
	"github.com/vikikurnia87/service-order/configs"
	"github.com/vikikurnia87/service-order/container"
	"github.com/vikikurnia87/service-order/database"
	"github.com/vikikurnia87/service-order/routes"
	"github.com/vikikurnia87/service-order/utils"

	"github.com/uptrace/bun"
	sudb "github.com/vikikurnia87/service-utils/database"
	"github.com/vikikurnia87/service-utils/monitoring"
)

const gracefulShutdownTimeout = 30 * time.Second

type App struct {
	Logger *slog.Logger

	ctx        context.Context
	db         *bun.DB
	server     *http.Server
	userClient *clients.UserClient
}

func New(ctx context.Context) *App {
	configs.LoadEnv()
	monitoring.InitAPM(configs.ElasticAPMServerURL)
	utils.SetValidatorLang(configs.ServerLang)

	logger := buildLogger()
	db := database.InitPostgresDatabase(buildDBConfig(), configs.ServerEnv, logger)
	database.InitRedisDatabase(ctx, logger)

	// gRPC client ke service-user (auth/validate token) — dipakai auth middleware.
	userClient, err := clients.NewUserClient(configs.UserGrpcAddr)
	if err != nil {
		logger.ErrorContext(ctx, "failed to dial service-user gRPC", slog.String("error", err.Error()))
	}

	c := container.New(container.Deps{DB: db, Logger: logger})

	return &App{
		Logger:     logger,
		ctx:        ctx,
		db:         db,
		server:     buildHTTPServer(routes.SetupRouter(c, userClient)),
		userClient: userClient,
	}
}

func (a *App) Start() {
	go func() {
		logServerStart(a.ctx, a.Logger)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.ErrorContext(a.ctx, "HTTP server error", slog.String("error", err.Error()))
		}
	}()
}

func (a *App) Shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.Logger.ErrorContext(shutdownCtx, "HTTP server shutdown error", slog.String("error", err.Error()))
	} else {
		a.Logger.InfoContext(shutdownCtx, "✅ HTTP server stopped")
	}
}

func (a *App) Cleanup() {
	ctx := context.Background()
	if a.db != nil {
		if err := database.Close(a.db); err != nil {
			a.Logger.Error("Failed to close database", slog.String("error", err.Error()))
		}
	}
	database.RedisClose(ctx, a.Logger)
	if a.userClient != nil {
		_ = a.userClient.Close()
	}
}

// ─── Build helpers ───────────────────────────────────────────

func buildLogger() *slog.Logger {
	var w io.Writer = os.Stdout
	if path := os.Getenv("LOG_FILE"); path != "" {
		if f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
			w = io.MultiWriter(os.Stdout, f)
		}
	}
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	return logger.With("service.name", configs.ServerName, "app", configs.ServerName, "env", configs.ServerEnv)
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

func buildHTTPServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", configs.ServerPort),
		Handler:      handler,
		ReadTimeout:  configs.ServerReadTimeout,
		WriteTimeout: configs.ServerWriteTimeout,
		IdleTimeout:  configs.ServerIdleTimeout,
	}
}

func logServerStart(ctx context.Context, logger *slog.Logger) {
	logger.InfoContext(ctx, "🚀 Server running", "address", fmt.Sprintf(":%s", configs.ServerPort))
}
