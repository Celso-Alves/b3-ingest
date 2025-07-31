package main

import (
	"b3-ingest/internal/infra/adapter/database"
	"b3-ingest/internal/infra/settings"
	"b3-ingest/internal/logger"
	"b3-ingest/internal/starter"
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		loadFlag     = flag.Bool("load", false, "Load CSV files into the database")
		serveFlag    = flag.Bool("serve", false, "Run HTTP server with trading routes")
		downloadFlag = flag.Bool("download", false, "Download and unzip last 7 workdays' files to bundle/b3files")
	)
	flag.Parse()

	if err := settings.LoadEnvs(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load environment variables: %v\n", err)
		os.Exit(1)
	}
	logger.InitDefaultLogger()
	cfg := settings.LoadConfig()
	log := logger.GetDefaultLogger()

	mode := ""
	if *downloadFlag {
		mode = "download"
	} else if *loadFlag {
		mode = "load"
	} else if *serveFlag {
		mode = "serve"
	}

	starterCfg := starter.StarterConfig{
		Mode:    mode,
		CSVPath: cfg.CSVPath,
		AppPort: cfg.AppPort,
		DSN:     cfg.DSN(),
		DBConfig: database.Config{
			Name:     cfg.DatabaseName,
			Host:     cfg.DatabaseHost,
			Username: cfg.DatabaseUsername,
			Password: cfg.DatabasePassword,
			Port:     cfg.DatabasePort,
			SSL:      cfg.DatabaseSSL,
		},
		Logger: log,
	}
	starter.Start(starterCfg)
}
