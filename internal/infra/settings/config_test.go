package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenValidEnvVarsWhenLoadEnvsThenConfigIsPopulated(t *testing.T) {
	// Arrange
	os.Setenv("DATABASE_NAME", "testdb")
	os.Setenv("DATABASE_PASSWORD", "testpass")
	os.Setenv("DATABASE_USERNAME", "testuser")
	os.Setenv("DATABASE_HOST", "testhost")
	os.Setenv("DATABASE_PORT", "1234")
	os.Setenv("CSV_PATH", "/tmp/csv")
	os.Setenv("APP_DEFAULT_PORT", "9999")
	os.Setenv("APP_NAME", "testapp")
	os.Setenv("INGESTION_CORES", "2")

	// Act
	err := LoadEnvs()
	cfg := LoadConfig()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "testdb", cfg.DatabaseName)
	assert.Equal(t, "testpass", cfg.DatabasePassword)
	assert.Equal(t, "testuser", cfg.DatabaseUsername)
	assert.Equal(t, "testhost", cfg.DatabaseHost)
	assert.Equal(t, 1234, cfg.DatabasePort)
	assert.Equal(t, "/tmp/csv", cfg.CSVPath)
	assert.Equal(t, "9999", cfg.AppPort)
	assert.Equal(t, "testapp", cfg.APPName)
	assert.Equal(t, 2, cfg.IngestionCores)
}

func TestGivenConfigWhenDSNThenReturnsCorrectString(t *testing.T) {
	// Arrange
	cfg := Config{
		DatabaseEnvironment: DatabaseEnvironment{
			DatabaseName:     "db",
			DatabasePassword: "pw",
			DatabaseUsername: "user",
			DatabaseHost:     "host",
			DatabasePort:     5555,
		},
	}

	// Act
	dsn := cfg.DSN()

	// Assert
	assert.Equal(t, "postgres://user:pw@host:5555/db", dsn)
}
