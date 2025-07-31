package starter

import (
	"b3-ingest/internal/infra/adapter/database"
	"b3-ingest/internal/infra/adapter/database/provider/postgres"
	"b3-ingest/internal/infra/repositories/trading"
	"b3-ingest/internal/logger"
	"b3-ingest/internal/service/ingestion"
	tradingServicePkg "b3-ingest/internal/service/trading"
	tradingRoute "b3-ingest/pkg/routes/v1/trading"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

type StarterConfig struct {
	Mode     string
	CSVPath  string
	AppPort  string
	DSN      string
	DBConfig database.Config
	Logger   *logger.Logger
}

func Start(cfg StarterConfig) {
	switch cfg.Mode {
	case "download":
		startDownload(cfg)
	case "load":
		startIngestion(cfg)
	case "serve":
		startServer(cfg)
	default:
		fmt.Println("Usage:")
		fmt.Println("  b3-ingest -load   # Load CSV files into the database")
		fmt.Println("  b3-ingest -serve  # Run HTTP server with trading routes")
		os.Exit(1)
	}
}

func startDownload(cfg StarterConfig) {
	cfg.Logger.Info("Downloading and extracting last 7 workdays' files...")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := ingestion.DownloadAndUnzipLast7Workdays(ctx, cfg.CSVPath, func(msg string, args ...interface{}) { cfg.Logger.Info(msg, args...) })
	if err != nil {
		cfg.Logger.Error("Download/unzip failed: %v", err)
		os.Exit(1)
	}
	cfg.Logger.Info("Download and extraction complete.")
}

func startIngestion(cfg StarterConfig) {
	db, err := postgres.NewPostgres(cfg.DBConfig)
	if err != nil {
		cfg.Logger.Error("Error connecting to database: %v", err)
		os.Exit(1)
	}
	cfg.Logger.Info("Starting CSV ingestion mode...")
	start := time.Now()
	ingestionService := ingestion.NewService(db, cfg.DSN, cfg.Logger)
	err = ingestionService.IngestFromCSV(cfg.CSVPath)
	if err != nil {
		cfg.Logger.Error("Error loading CSV data: %v", err)
		os.Exit(1)
	}
	elapsed := time.Since(start)
	cfg.Logger.Info("LoadFromCSV completed: %f", elapsed.Seconds())
}

func startServer(cfg StarterConfig) {
	db, err := postgres.NewPostgres(cfg.DBConfig)
	if err != nil {
		cfg.Logger.Error("Error connecting to database: %v", err)
		os.Exit(1)
	}
	cfg.Logger.Info("Starting HTTP server mode...")
	repo := trading.NewTradingRepository()
	service := tradingServicePkg.NewTradingService(repo, db)
	r := gin.Default()
	r.GET("/quote", tradingRoute.GetQuoteHandler(service))
	port := cfg.AppPort
	if port == "" {
		port = "8000"
	}
	addr := fmt.Sprintf(":%s", port)
	cfg.Logger.Info("HTTP server running, addr: %s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.Logger.Error("Failed to start HTTP server: %v", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	cfg.Logger.Info("Shutdown signal received, shutting down HTTP server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		cfg.Logger.Error("HTTP server forced to shutdown: %v", err)
	}
	cfg.Logger.Info("HTTP server stopped.")
}
