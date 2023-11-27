package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type (
	InputConfig struct {
		Entries []InputConfigEntry `parser:"@@*"`
	}

	InputConfigEntry interface {
		EmitTo(ProxySettings)
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
		EmitTo(ProxySettings)
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

// EmitTo implements InputConfigEntry.
func (entry ProxyConfigEntry) EmitTo(settings ProxySettings) {
	settings.UseProxy(entry)
}

// EmitTo implements InputConfigEntry.
func (entry IncludeFileConfigEntry) EmitTo(settings ProxySettings) {
	panic("unimplemented")
}

// EmitTo implements InputConfigEntry.
func (entry DomainConfigEntry) EmitTo(settings ProxySettings) {
	entry.Match.EmitTo(settings)
}

// EmitTo implements DomainMatch.
func (match SubDomainMatch) EmitTo(settings ProxySettings) {
	settings.AddSubdomainMatch(match)
}

// EmitTo implements DomainMatch.
func (RegexpMatch) EmitTo(settings ProxySettings) {
	panic("unimplemented")
}

// EmitTo implements DomainMatch.
func (KeywordMatch) EmitTo(settings ProxySettings) {
	panic("unimplemented")
}

// EmitTo implements DomainMatch.
func (FullDomainMatch) EmitTo(settings ProxySettings) {
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
