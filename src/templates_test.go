package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTemplates(t *testing.T) {
	tmpl := initTemplates()
	amount := 7

	assert.Equal(t, len(tmpl.templates), amount, fmt.Sprintf("Should be %d templates", amount))
}
