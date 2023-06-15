package logging

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoLogger(t *testing.T) {
	level := "info"
	log := CreateLogger(level)
	defer log.Sync()

	assert.NotEmpty(t, log, "Logger should not be nil")
	assert.Equal(t, log.Level().String(), level, fmt.Sprintf("Level should be %s", level))
}

func TestWarnLogger(t *testing.T) {
	level := "warn"
	log := CreateLogger(level)
	defer log.Sync()

	assert.NotEmpty(t, log, "Logger should not be nil")
	assert.Equal(t, log.Level().String(), level, fmt.Sprintf("Level should be %s", level))
}

func TestDebugLogger(t *testing.T) {
	level := "debug"
	log := CreateLogger(level)
	defer log.Sync()

	assert.NotEmpty(t, log, "Logger should not be nil")
	assert.Equal(t, log.Level().String(), level, fmt.Sprintf("Level should be %s", level))
}

func TestInvalidLogger(t *testing.T) {
	level := "invalid"
	log := CreateLogger(level)
	defer log.Sync()

	assert.NotEmpty(t, log, "Logger should not be nil")
	assert.Equal(t, log.Level().String(), "info", "Level should be info")
}
