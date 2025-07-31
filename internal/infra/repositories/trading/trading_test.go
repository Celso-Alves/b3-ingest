package trading

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Trading struct {
	ID                  uint `gorm:"primaryKey"`
	DataNegocio         time.Time
	CodigoInstrumento   string
	PrecoNegocio        float64
	QuantidadeNegociada int64
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&Trading{})
	return db
}

func TestGivenValidDataWhenGetQuoteStatsThenReturnsStats(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	db.Create(&Trading{
		DataNegocio:         time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC),
		CodigoInstrumento:   "WDOQ25",
		PrecoNegocio:        100.0,
		QuantidadeNegociada: 500,
	})
	db.Create(&Trading{
		DataNegocio:         time.Date(2025, 7, 30, 0, 0, 0, 0, time.UTC),
		CodigoInstrumento:   "WDOQ25",
		PrecoNegocio:        200.0,
		QuantidadeNegociada: 1000,
	})
	repo := NewTradingRepository()

	// Act
	stats, err := repo.GetQuoteStats(context.Background(), db, "WDOQ25", time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200.0, stats.MaxPrice)
	assert.Equal(t, int64(1000), stats.MaxDailyVolume)
}

func TestGivenNoDataWhenGetQuoteStatsThenReturnsZeroStats(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewTradingRepository()

	// Act
	stats, err := repo.GetQuoteStats(context.Background(), db, "FOO", time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0.0, stats.MaxPrice)
	assert.Equal(t, int64(0), stats.MaxDailyVolume)
}

func TestGivenDBErrorWhenGetQuoteStatsThenReturnsError(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	db.Migrator().DropTable(&Trading{}) // force error
	repo := NewTradingRepository()

	// Act
	_, err := repo.GetQuoteStats(context.Background(), db, "WDOQ25", time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC))

	// Assert
	assert.Error(t, err)
}
