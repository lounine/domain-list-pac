package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleConfig(t *testing.T) {
	proxy := ProxyPac{}
	proxy.ReadConfig(readTestConfig(t, "simple"))

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "simple.pac", settings)
}

func TestMultipleProxiesConfig(t *testing.T) {
	proxy := ProxyPac{}
	proxy.ReadConfig(readTestConfig(t, "multiple-proxies"))

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "multiple-proxies.pac", settings)
}

func TestRepeatingProxiesConfig(t *testing.T) {
	proxy := ProxyPac{}
	proxy.ReadConfig(readTestConfig(t, "repeating-proxies"))

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "repeating-proxies.pac", settings)
}

func TestNoProxyConfig(t *testing.T) {
	proxy := ProxyPac{}

	assert.PanicsWithError(t,
		fmt.Sprintf(NoProxyErrorMessage, "some-domain.com"),
		func() { proxy.ReadConfig(readTestConfig(t, "no-proxy")) },
	)
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

func (proxy ProxyPac) GenerateSettings() string {
	var result bytes.Buffer
	proxy.WriteSettings(&result)

	return result.String()
}

func assertEqualToFile(t *testing.T, filename string, result string) {
	expected, err := os.ReadFile("./test/expect/" + filename)
	assert.NoError(t, err)
	assert.Equal(t, string(expected), result)
}
