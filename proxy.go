package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ProxySettings interface {
	UseProxy(ProxyConfigEntry)
	AddSubdomainMatch(SubDomainMatch)
	AddKeywordMatch(KeywordMatch)
	AddFullDomainMatch(FullDomainMatch)
	AddRegexpMatch(RegexpMatch)
	ReadConfig(string)
	WriteSettings(io.Writer)
}

// A Proxy auto-config file generator
// Implements ProxySettings interface
type ProxyPac struct {
	Proxies      []ProxyVariable
	Statements   []ProxyStatement
	CurrentScope *FileScope
}

const NoProxyErrorMessage = "no proxy to use for %v"
const CantAccessFileErrorMessage = "error accessing file %v"
const FileNotFoundErrorMessage = "can't find included file: %v"

type ProxyVariable struct {
	Name    string
	Type    string
	Address string
}

func (variable ProxyVariable) String() string {
	return fmt.Sprintf("var %v = '%v %v';", variable.Name, variable.Type, variable.Address)
}

type ProxyStatement struct {
	Check string
	Proxy *ProxyVariable
}

func (statement ProxyStatement) String() string {
	return fmt.Sprintf("%v return %v;", statement.Check, statement.Proxy.Name)
}

type FileScope struct {
	CurrentDirectory    string
	AdditionalLocations []string
	CurrentProxy        *ProxyVariable
	PreviousScope       *FileScope
}

func (scope FileScope) LocalSubScope() *FileScope {
	return &FileScope{
		CurrentDirectory:    scope.CurrentDirectory,
		AdditionalLocations: scope.AdditionalLocations,
		CurrentProxy:        scope.CurrentProxy,
		PreviousScope:       &scope,
	}
}

func (scope FileScope) SubScope(directory string) *FileScope {
	return &FileScope{
		CurrentDirectory: directory,
		CurrentProxy:     scope.CurrentProxy,
		PreviousScope:    &scope,
	}
}

// UseProxy implements ProxySettings.
func (proxy *ProxyPac) UseProxy(entry ProxyConfigEntry) {
	for _, p := range proxy.Proxies {
		if p.Address == entry.Address {
			proxy.CurrentScope.CurrentProxy = &p
			return
		}
	}

	newProxy := ProxyVariable{Type: entry.Type, Address: entry.Address}

	if len(proxy.Proxies) == 0 {
		newProxy.Name = "p"
	} else if len(proxy.Proxies) == 1 {
		proxy.Proxies[0].Name = "p1"
		newProxy.Name = "p2"
	} else {
		newProxy.Name = fmt.Sprintf("p%v", len(proxy.Proxies)+1)
	}

	proxy.Proxies = append(proxy.Proxies, newProxy)
	proxy.CurrentScope.CurrentProxy = &proxy.Proxies[len(proxy.Proxies)-1]
}

// AddSubdomainMatch implements ProxySettings.
func (proxy *ProxyPac) AddSubdomainMatch(match SubDomainMatch) {
	proxy.addMatch(
		fmt.Sprintf(`if (dnsDomainIs(h, '%v'))`, match.Value),
		match,
	)
}

// AddFullDomainMatch implements ProxySettings.
func (proxy *ProxyPac) AddFullDomainMatch(match FullDomainMatch) {
	proxy.addMatch(
		fmt.Sprintf(`if (h == '%v')`, match.Value),
		match,
	)
}

// AddKeywordMatch implements ProxySettings.
func (proxy *ProxyPac) AddKeywordMatch(match KeywordMatch) {
	proxy.addMatch(
		fmt.Sprintf(`if (shExpMatch(h, '*%v*'))`, match.Value),
		match,
	)
}

// AddKeywordMatch implements ProxySettings.
func (proxy *ProxyPac) AddRegexpMatch(match RegexpMatch) {
	proxy.addMatch(
		fmt.Sprintf(`if (/%v/.test(h))`, match.Value),
		match,
	)
}

func (proxy *ProxyPac) addMatch(check string, match DomainMatch) {
	if proxy.CurrentScope.CurrentProxy == nil {
		panic(fmt.Errorf(NoProxyErrorMessage, match))
	}

	statement := ProxyStatement{Check: check, Proxy: proxy.CurrentScope.CurrentProxy}
	proxy.Statements = append(proxy.Statements, statement)
}

func (proxy *ProxyPac) writeProxies(out io.Writer) {
	for _, p := range proxy.Proxies {
		io.WriteString(out, fmt.Sprintf("\t%v\n", p.String()))
	}
}

func (proxy *ProxyPac) writeStatements(out io.Writer) {
	for _, s := range proxy.Statements {
		io.WriteString(out, fmt.Sprintf("\t%v\n", s.String()))
	}
}

func NewProxyPac(configPath string, configLocations []string) *ProxyPac {
	proxy := &ProxyPac{}
	scope := &FileScope{
		CurrentDirectory:    filepath.Dir(configPath),
		AdditionalLocations: configLocations,
	}

	proxy.applyConfigFromPath(configPath, scope)
	return proxy
}

// ReadConfig implements ProxySettings.
func (proxy *ProxyPac) ReadConfig(filename string) {
	if proxy.tryReadConfig(
		proxy.CurrentScope.CurrentDirectory,
		filename,
		proxy.CurrentScope.LocalSubScope()) {
		return
	}

	for _, directory := range proxy.CurrentScope.AdditionalLocations {
		if proxy.tryReadConfig(
			directory,
			filename,
			proxy.CurrentScope.SubScope(directory)) {
			return
		}
	}

	panic(fmt.Errorf(FileNotFoundErrorMessage, filename))
}

func (proxy *ProxyPac) tryReadConfig(directory string, filename string, scope *FileScope) bool {
	path := filepath.Join(directory, filename)
	if _, err := os.Stat(path); err == nil {
		proxy.applyConfigFromPath(path, scope)
		return true
	} else if !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf(CantAccessFileErrorMessage, path))
	}

	return false
}

func (proxy *ProxyPac) applyConfigFromPath(path string, scope *FileScope) {
	config := NewConfigFrom(path)
	proxy.CurrentScope = scope
	for _, e := range config.Entries {
		e.ApplyTo(proxy)
	}
	proxy.CurrentScope = proxy.CurrentScope.PreviousScope
}

// WriteSettings implements ProxySettings.
func (p *ProxyPac) WriteSettings(out io.Writer) {
	io.WriteString(out, "function FindProxyForURL (url, host) {\n")
	io.WriteString(out, "\tvar h = host.toLowerCase();\n")

	p.writeProxies(out)
	io.WriteString(out, "\n")
	p.writeStatements(out)
	io.WriteString(out, "\n")

	io.WriteString(out, "\treturn 'DIRECT';\n")
	io.WriteString(out, "}\n")
}
