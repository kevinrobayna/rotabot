package main

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestName(t *testing.T) {
	assert.Equal(t, hello(), "Hello world from rotabot running unknown built on unknown")
}
