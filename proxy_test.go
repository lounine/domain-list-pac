package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	proxy := NewProxyPac("./test/config/simple", []string{})
	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "simple.pac", settings)
}

func TestMultipleProxies(t *testing.T) {
	proxy := NewProxyPac("./test/config/proxies-multiple", []string{})
	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "proxies-multiple.pac", settings)
}

func TestRepeatingProxies(t *testing.T) {
	proxy := NewProxyPac("./test/config/proxies-repeating", []string{})
	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "proxies-repeating.pac", settings)
}

func TestNoProxy(t *testing.T) {
	assert.PanicsWithError(t,
		fmt.Sprintf(NoProxyErrorMessage, "some-domain.com"),
		func() { NewProxyPac("./test/config/proxies-missing", []string{}) },
	)
}

func TestSimpleInclude(t *testing.T) {
	proxy := NewProxyPac("./test/config/include-simple", []string{})
	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "include-simple.pac", settings)
}

func TestIncludeFromAnotherDirectory(t *testing.T) {
	proxy := NewProxyPac("./test/config/include-other-dir", []string{"./test/data"})
	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "include-simple.pac", settings)
}

func TestIncludesAreScoped(t *testing.T) {
	proxy := NewProxyPac("./test/config/include-scopes", []string{})
	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "include-scopes.pac", settings)
}

func TestIncludesAreLocalized(t *testing.T) {
	proxy := NewProxyPac("./test/config/include-duplicates", []string{"./test/data"})
	settings := proxy.GenerateSettings()
	assertEqualToFile(t, "include-duplicates.pac", settings)
}

func (proxy ProxyPac) GenerateSettings() string {
	var result bytes.Buffer
	proxy.WriteSettings(&result)

	return result.String()
}

func assertEqualToFile(t *testing.T, filename string, result string) {
	expected, err := os.ReadFile("./test/expect/" + filename)
	assert.NoError(t, err)
	if err == nil {
		assert.Equal(t, string(expected), result)
	}
}
