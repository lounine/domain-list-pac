package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqualToFile(t *testing.T, file string, result bytes.Buffer) {
	expected, err := os.ReadFile("./test/expect/" + file)
	assert.NoError(t, err)
	assert.Equal(t, string(expected), result.String())
}

func TestSimpleConfig(t *testing.T) {
	input, err := os.Open("./test/config/simple")
	assert.NoError(t, err)

	config, err := parser.Parse("", input)
	assert.NoError(t, err)

	var result bytes.Buffer
	proxyPac := ProxyPac{}
	proxyPac.WriteSettings(config, &result)

	AssertEqualToFile(t, "simple.pac", result)
}

