package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleConfig(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("simple")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "simple.pac", settings)
}

func TestMultipleProxiesConfig(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("multiple-proxies")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "multiple-proxies.pac", settings)
}

func TestRepeatingProxiesConfig(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("repeating-proxies")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "repeating-proxies.pac", settings)
}

func TestNoProxyConfig(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}

	assert.PanicsWithError(t,
		fmt.Sprintf(NoProxyErrorMessage, "some-domain.com"),
		func() { proxy.ReadConfig("no-proxy") },
	)
}

func TestSimpleIncludeConfig(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("simple-include")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "simple-include.pac", settings)
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
