package trading

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockTradingService struct{}

func (m *mockTradingService) GetQuote(ctx context.Context, ticker string, startDate time.Time) (float64, int64, error) {
	if ticker == "FAIL" {
		return 0, 0, assert.AnError
	}
	return 123.45, 6789, nil
}

func TestGetQuoteHandlerGivenValidTickerAndDateWhenRequestIsMadeThenReturnsSuccess(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/quote", GetQuoteHandler(&mockTradingService{}))
	req, _ := http.NewRequest("GET", "/quote?ticker=TEST&data_inicio=2025-07-01", nil)

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp QuoteResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "TEST", resp.Ticker)
	assert.Equal(t, 123.45, resp.MaxRangeValue)
	assert.Equal(t, int64(6789), resp.MaxDailyVolume)
}

func TestGetQuoteHandlerGivenMissingTickerWhenRequestIsMadeThenReturnsBadRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/quote", GetQuoteHandler(&mockTradingService{}))
	req, _ := http.NewRequest("GET", "/quote", nil)

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetQuoteHandlerGivenInvalidDateWhenRequestIsMadeThenReturnsBadRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/quote", GetQuoteHandler(&mockTradingService{}))
	req, _ := http.NewRequest("GET", "/quote?ticker=TEST&data_inicio=invalid-date", nil)

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetQuoteHandlerGivenServiceErrorWhenRequestIsMadeThenReturnsInternalServerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/quote", GetQuoteHandler(&mockTradingService{}))
	req, _ := http.NewRequest("GET", "/quote?ticker=FAIL", nil)

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
