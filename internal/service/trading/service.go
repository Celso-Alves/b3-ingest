package trading

import (
	"b3-ingest/internal/infra/repositories/trading"
	"context"
	"time"

	"gorm.io/gorm"
)

type TradingService interface {
	GetQuote(ctx context.Context, ticker string, startDate time.Time) (maxPrice float64, maxVol int64, err error)
}

type tradingService struct {
	repo trading.TradingRepository
	db   *gorm.DB
}

func NewTradingService(repo trading.TradingRepository, db *gorm.DB) TradingService {
	return &tradingService{repo: repo, db: db}
}

func (s *tradingService) GetQuote(ctx context.Context, ticker string, startDate time.Time) (float64, int64, error) {
	stats, err := s.repo.GetQuoteStats(ctx, s.db, ticker, startDate)
	if err != nil {
		return 0, 0, err
	}
	return stats.MaxPrice, stats.MaxDailyVolume, nil
}
