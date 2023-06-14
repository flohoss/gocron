package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLogger(t *testing.T) {
	level := "info"
	log := createLogger(level)

	assert.NotEmpty(t, log, "Logger should not be nil")
	assert.Equal(t, log.Level().String(), level, fmt.Sprintf("Level should be %s", level))

	level = "warn"
	log = createLogger(level)

	assert.NotEmpty(t, log, "Logger should not be nil")
	assert.Equal(t, log.Level().String(), level, fmt.Sprintf("Level should be %s", level))
	log.Sync()

	level = "invalid"
	log = createLogger(level)

	assert.NotEmpty(t, log, "Logger should not be nil")
	assert.Equal(t, log.Level().String(), "info", "Level should be info")
	log.Sync()
}
