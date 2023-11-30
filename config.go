package main

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type (
	InputConfig struct {
		Entries []InputConfigEntry `parser:"@@*"`
	}

	InputConfigEntry interface {
		ApplyTo(ProxySettings)
	}

	ProxyConfigEntry struct {
		Type    string `parser:"@Proxy"`
		Address string `parser:"@Value"`
	}

	IncludeFileConfigEntry struct {
		Value string `parser:"'include:' @Value"`
	}

	DomainConfigEntry struct {
		Match      DomainMatch `parser:"@@"`
		Attributes []string    `parser:"@Attribute*"`
	}

	DomainMatch interface {
		ApplyTo(ProxySettings)
	}

	SubDomainMatch struct {
		Value string `parser:"('domain:')? @Value"`
	}

	FullDomainMatch struct {
		Value string `parser:"'full:' @Value"`
	}

	KeywordMatch struct {
		Value string `parser:"'keyword:' @Value"`
	}

	RegexpMatch struct {
		Value string `parser:"'regexp:' @Value"`
	}
)

func NewConfigFrom(path string) *InputConfig {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	} else {
		defer file.Close()
	}

	config, err := parser.Parse(path, file)
	if err != nil {
		panic(err)
	}

	return config
}

// ApplyTo implements InputConfigEntry.
func (entry ProxyConfigEntry) ApplyTo(settings ProxySettings) {
	settings.UseProxy(entry)
}

// ApplyTo implements InputConfigEntry.
func (entry IncludeFileConfigEntry) ApplyTo(settings ProxySettings) {
	settings.ReadConfig(entry.Value)
}

// ApplyTo implements InputConfigEntry.
func (entry DomainConfigEntry) ApplyTo(settings ProxySettings) {
	entry.Match.ApplyTo(settings)
}

// ApplyTo implements DomainMatch.
func (match SubDomainMatch) ApplyTo(settings ProxySettings) {
	settings.AddSubdomainMatch(match)
}

// ApplyTo implements DomainMatch.
func (RegexpMatch) ApplyTo(settings ProxySettings) {
	panic("unimplemented")
}

// ApplyTo implements DomainMatch.
func (KeywordMatch) ApplyTo(settings ProxySettings) {
	panic("unimplemented")
}

// ApplyTo implements DomainMatch.
func (FullDomainMatch) ApplyTo(settings ProxySettings) {
	panic("unimplemented")
}

var (
	lexerDef = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Proxy", Pattern: `(PROXY|SOCKS)\b`},
		{Name: "Key", Pattern: `[a-z]+:`},
		{Name: "Value", Pattern: `[^#@\s]+`},
		{Name: "Attribute", Pattern: `@[^#@\s]+`},
		{Name: "comment", Pattern: `#[^\n]*`},
		{Name: "whitespace", Pattern: `\s+`},
	})

	parser = participle.MustBuild[InputConfig](
		participle.Lexer(lexerDef),
		participle.Union[InputConfigEntry](ProxyConfigEntry{}, IncludeFileConfigEntry{}, DomainConfigEntry{}),
		participle.Union[DomainMatch](SubDomainMatch{}, FullDomainMatch{}, KeywordMatch{}, RegexpMatch{}),
	)
)
