package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_foobar(t *testing.T) {
	expected := true
	actual := true
	assert.Equal(t, expected, actual, "compare result")
}
