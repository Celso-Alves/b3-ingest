package starter

import (
	"b3-ingest/internal/logger"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeLogger struct{}

func (f *fakeLogger) Info(format string, args ...interface{})  {}
func (f *fakeLogger) Error(format string, args ...interface{}) {}

func TestStartGivenUnknownModeWhenCalledThenExitsWithUsage(t *testing.T) {
	// Arrange
	cfg := StarterConfig{
		Mode:   "unknown",
		Logger: logger.NewLogger(io.Discard, "", 0, logger.INFO),
	}
	// Act & Assert
	if os.Getenv("BE_CRASHER") == "1" {
		Start(cfg)
		return
	}
	cmd := os.Args[0]
	os.Setenv("BE_CRASHER", "1")
	p, err := os.StartProcess(cmd, []string{cmd, "-test.run=TestStartGivenUnknownModeWhenCalledThenExitsWithUsage"}, &os.ProcAttr{Env: os.Environ(), Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	assert.NoError(t, err)
	state, err := p.Wait()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, state.ExitCode())
	os.Unsetenv("BE_CRASHER")
}
