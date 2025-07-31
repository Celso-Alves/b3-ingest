package database

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Name            string
	Host            string
	Username        string
	Password        string
	Port            int
	SSL             bool
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

type Database interface {
	Get() *gorm.DB
	Close() error
}
