package trading

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type QuoteStats struct {
	MaxPrice       float64
	MaxDailyVolume int64
}

type TradingRepository interface {
	GetQuoteStats(ctx context.Context, db *gorm.DB, ticker string, startDate time.Time) (QuoteStats, error)
}

type tradingRepository struct{}

func NewTradingRepository() TradingRepository {
	return &tradingRepository{}
}

func (r *tradingRepository) GetQuoteStats(ctx context.Context, db *gorm.DB, ticker string, startDate time.Time) (QuoteStats, error) {
	var stats QuoteStats

	query := `
		WITH daily_volumes AS (
			SELECT data_negocio, SUM(quantidade_negociada) AS soma
			FROM tradings
			WHERE codigo_instrumento = ? AND data_negocio >= ?
			GROUP BY data_negocio
		)
		SELECT 
			COALESCE(MAX(t.preco_negocio), 0) AS max_price,
			COALESCE(MAX(dv.soma), 0) AS max_daily_volume
		FROM tradings t
		LEFT JOIN daily_volumes dv ON t.data_negocio = dv.data_negocio
		WHERE t.codigo_instrumento = ? AND t.data_negocio >= ?
	`
	row := db.WithContext(ctx).Raw(query, ticker, startDate, ticker, startDate).Row()
	if err := row.Scan(&stats.MaxPrice, &stats.MaxDailyVolume); err != nil {
		return QuoteStats{}, err
	}
	return stats, nil
}
