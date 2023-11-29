package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("simple")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "simple.pac", settings)
}

func TestMultipleProxies(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("proxies-multiple")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "proxies-multiple.pac", settings)
}

func TestRepeatingProxies(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("proxies-repeating")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "proxies-repeating.pac", settings)
}

func TestNoProxy(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}

	assert.PanicsWithError(t,
		fmt.Sprintf(NoProxyErrorMessage, "some-domain.com"),
		func() { proxy.ReadConfig("proxies-missing") },
	)
}

func TestSimpleInclude(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config"}}
	proxy.ReadConfig("include-simple")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "include-simple.pac", settings)
}

func TestIncludeFromAnotherDirectory(t *testing.T) {
	proxy := ProxyPac{ConfigLocations: []string{"./test/config", "./test/data"}}
	proxy.ReadConfig("include-other-dir")

	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "include-simple.pac", settings)
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
