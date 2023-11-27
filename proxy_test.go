package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleConfig(t *testing.T) {
	config := readTestConfig(t, "simple")

	proxyPac := ProxyPac{}
	proxyPac.ReadConfig(config)

	contents := buildProxyPacContents(t, &proxyPac)
	assertEqualToFile(t, "simple.pac", contents)
}

func TestMultipleProxiesConfig(t *testing.T) {
	config := readTestConfig(t, "multiple-proxies")

	proxyPac := ProxyPac{}
	proxyPac.ReadConfig(config)

	contents := buildProxyPacContents(t, &proxyPac)
	assertEqualToFile(t, "multiple-proxies.pac", contents)
}

func readTestConfig(t *testing.T, filename string) *InputConfig {
	input, err := os.Open("./test/config/" + filename)
	assert.NoError(t, err)
	if err == nil {
		defer input.Close()
	}

	config, err := parser.Parse("", input)
	assert.NoError(t, err)

	return config
}

func buildProxyPacContents(t *testing.T, proxyPac *ProxyPac) string {
	var result bytes.Buffer
	proxyPac.WriteSettings(&result)

	return result.String()
}

func assertEqualToFile(t *testing.T, filename string, result string) {
	expected, err := os.ReadFile("./test/expect/" + filename)
	assert.NoError(t, err)
	assert.Equal(t, string(expected), result)
}
