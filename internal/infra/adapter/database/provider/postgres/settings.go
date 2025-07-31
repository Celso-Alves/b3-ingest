package postgres

import (
	"b3-ingest/internal/infra/adapter/database"
	"b3-ingest/internal/infra/adapter/database/orm/models"
	"fmt"

	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func New(config database.Config) (database.Database, error) {
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=America/Sao_Paulo",
		config.Host, config.Username, config.Password, config.Name, config.Port, getSSLMode(config.SSL))

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connection, PreferSimpleProtocol: true}), &gorm.Config{
		Logger:                 logger.Discard,
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &Postgres{
		db: db,
	}, nil
}

func NewPostgres(config database.Config) (*gorm.DB, error) {
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=America/Sao_Paulo",
		config.Host, config.Username, config.Password, config.Name, config.Port, getSSLMode(config.SSL))
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connection, PreferSimpleProtocol: true}), &gorm.Config{
		Logger:                 logger.Discard,
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Trading{})
	return db, nil
}

func (p *Postgres) Get() *gorm.DB {
	return p.db
}

func (p *Postgres) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

func getSSLMode(sslmode bool) string {
	if sslmode {
		return "require"
	}
	return "disable"
}
