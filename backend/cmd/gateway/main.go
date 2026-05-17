package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/macmakobaseballer/support-ops-hub/backend/gateway"
	"github.com/macmakobaseballer/support-ops-hub/backend/internal/config"
	db "github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	sqlDB, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		slog.Error("failed to open db", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("failed to close db", slog.Any("error", err))
		}
	}()

	if err := sqlDB.Ping(); err != nil {
		slog.Error("failed to ping db", slog.Any("error", err))
		os.Exit(1)
	}

	queries := db.New(sqlDB)
	router := gateway.NewRouter(sqlDB, queries, cfg.CORSAllowedOrigins)
	srv := gateway.NewServer(cfg.GatewayPort, router)

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("gateway exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}
