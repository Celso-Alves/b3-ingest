package ingestion

import (
	"b3-ingest/internal/logger"
	"io"

	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIngestFromCSVGivenNonExistentDirWhenCalledThenReturnsError(t *testing.T) {
	// Arrange
	testLogger := logger.NewLogger(io.Discard, "", 0, logger.INFO)
	s := &Service{DB: &gorm.DB{}, DSN: "", Log: testLogger}
	dir := "./nonexistent_dir"

	// Act
	err := s.IngestFromCSV(dir)

	// Assert
	assert.Error(t, err)
}
