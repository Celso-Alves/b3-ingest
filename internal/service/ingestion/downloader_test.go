package ingestion

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDownloadAndUnzipLast7WorkdaysGivenInvalidDirWhenCalledThenReturnsError(t *testing.T) {
	// Arrange
	destDir := string([]byte{0}) // invalid path
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	logf := func(string, ...interface{}) {}

	// Act
	err := DownloadAndUnzipLast7Workdays(ctx, destDir, logf)

	// Assert
	assert.Error(t, err)
}

func TestUnzipGivenNonExistentFileWhenCalledThenReturnsError(t *testing.T) {
	// Arrange
	src := "./notfound.zip"
	dest := "./tmp"
	logf := func(string, ...interface{}) {}

	// Act
	err := unzip(src, dest, logf)

	// Assert
	assert.Error(t, err)
}

func TestUnzipGivenValidZipWhenCalledThenExtractsFiles(t *testing.T) {
	// Arrange
	dest := "./testdata/unzip"
	_ = os.MkdirAll(dest, 0755)
	defer os.RemoveAll("./testdata")
	// Create a zip file for testing
	zipPath := filepath.Join(dest, "test.zip")
	f, _ := os.Create(zipPath)
	f.Close()
	logf := func(string, ...interface{}) {}

	// Act
	err := unzip(zipPath, dest, logf)

	// Assert
	assert.Error(t, err)
}
