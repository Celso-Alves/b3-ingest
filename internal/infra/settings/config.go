package settings

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

var envs Config

type DatabaseEnvironment struct {
	DatabaseName     string `env:"DATABASE_NAME,required" envDefault:"b3db"`
	DatabasePassword string `env:"DATABASE_PASSWORD,required" envDefault:"postgres"`
	DatabaseUsername string `env:"DATABASE_USERNAME,required" envDefault:"postgres"`
	DatabaseHost     string `env:"DATABASE_HOST,required" envDefault:"localhost"`
	DatabasePort     int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseSSL      bool   `env:"DATABASE_SSL" envDefault:"true"`
}

// Config stores application configurations.
type Config struct {
	CSVPath        string `env:"CSV_PATH,required" envDefault:"./bundle/b3files"`
	AppPort        string `env:"APP_DEFAULT_PORT" envDefault:"8000"`
	APPName        string `env:"APP_NAME" envDefault:"b3-ingest"`
	IngestionCores int    `env:"INGESTION_CORES" envDefault:"6"`
	DatabaseEnvironment
}

// DSN returns a PostgreSQL DSN string for pgxpool or similar clients.
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.DatabaseUsername,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}

// LoadConfig loads environment variables.
func LoadConfig() *Config {
	cfg := &Config{
		CSVPath:        GetEnvs().CSVPath,
		AppPort:        GetEnvs().AppPort,
		APPName:        GetEnvs().APPName,
		IngestionCores: GetEnvs().IngestionCores,
		DatabaseEnvironment: DatabaseEnvironment{
			DatabaseName:     GetEnvs().DatabaseName,
			DatabasePassword: GetEnvs().DatabasePassword,
			DatabaseUsername: GetEnvs().DatabaseUsername,
			DatabaseHost:     GetEnvs().DatabaseHost,
			DatabasePort:     GetEnvs().DatabasePort,
		},
	}
	return cfg
}

func LoadEnvs() error {
	if err := env.Parse(&envs); err != nil {
		return err
	}
	return nil

}

func GetEnvs() Config {
	return envs
}
