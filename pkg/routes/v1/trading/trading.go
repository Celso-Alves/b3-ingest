package trading

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type QuoteResponse struct {
	Ticker         string  `json:"ticker"`
	MaxRangeValue  float64 `json:"max_range_value"`
	MaxDailyVolume int64   `json:"max_daily_volume"`
}

type TradingService interface {
	GetQuote(ctx context.Context, ticker string, startDate time.Time) (maxPrice float64, maxVol int64, err error)
}

func GetQuoteHandler(svc TradingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticker := c.Query("ticker")
		if ticker == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticker is required"})
			return
		}

		dataInicio := c.Query("data_inicio")
		var startDate time.Time
		var err error
		fmt.Printf("Received request for ticker: %s, startDate: %s\n", ticker, dataInicio)

		if dataInicio != "" {
			startDate, err = time.Parse("2006-01-02", dataInicio)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "data_inicio inválida, use formato YYYY-MM-DD"})
				return
			}
		} else {
			startDate = time.Now().AddDate(0, 0, -7) // 7 dias atrás
		}

		maxPrice, maxVol, err := svc.GetQuote(c.Request.Context(), ticker, startDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, QuoteResponse{
			Ticker:         ticker,
			MaxRangeValue:  maxPrice,
			MaxDailyVolume: maxVol,
		})
	}
}
