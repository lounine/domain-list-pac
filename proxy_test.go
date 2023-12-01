package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	settings := generateSettings("simple", []string{})
	assertEqualToFile(t, "simple.pac", settings)
}

func TestMultipleProxies(t *testing.T) {
	settings := generateSettings("proxies-multiple", []string{})
	assertEqualToFile(t, "proxies-multiple.pac", settings)
}

func TestRepeatingProxies(t *testing.T) {
	settings := generateSettings("proxies-repeating", []string{})
	assertEqualToFile(t, "proxies-repeating.pac", settings)
}

func TestNoProxy(t *testing.T) {
	assert.PanicsWithError(t,
		fmt.Sprintf(NoProxyErrorMessage, "domain:some-domain.com"),
		func() { generateSettings("proxies-missing", []string{}) },
	)
}

func TestSimpleInclude(t *testing.T) {
	settings := generateSettings("include-simple", []string{})
	assertEqualToFile(t, "include-simple.pac", settings)
}

func TestIncludeFromAnotherDirectory(t *testing.T) {
	settings := generateSettings("include-other-dir", []string{"./test/data"})
	assertEqualToFile(t, "include-simple.pac", settings)
}

func TestIncludesAreScoped(t *testing.T) {
	settings := generateSettings("include-scopes", []string{})
	assertEqualToFile(t, "include-scopes.pac", settings)
}

func TestIncludesAreLocalized(t *testing.T) {
	settings := generateSettings("include-duplicates", []string{"./test/data"})
	assertEqualToFile(t, "include-duplicates.pac", settings)
}

func TestDomainsRulesHandling(t *testing.T) {
	settings := generateSettings("domain-rules", []string{})
	assertEqualToFile(t, "domain-rules.pac", settings)
}

func TestIncludeNotFound(t *testing.T) {
	assert.PanicsWithError(t,
		fmt.Sprintf(FileNotFoundErrorMessage, "missing-file"),
		func() { generateSettings("include-missing", []string{}) },
	)
}

func generateSettings(filename string, locations []string) string {
	proxy := NewProxyPac("./test/config/"+filename, locations)

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
