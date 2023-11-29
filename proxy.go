package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type ProxySettings interface {
	UseProxy(ProxyConfigEntry)
	AddSubdomainMatch(SubDomainMatch)
	ReadConfig(string)
	WriteSettings(io.Writer)
}

// A Proxy auto-config file generator
// Implements ProxySettings interface
type ProxyPac struct {
	ConfigLocations []string
	Proxies         []ProxyVariable
	Statements      []ProxyStatement
	CurrentScope    *FileScope
}

const NoProxyErrorMessage = "no proxy to use for %v"
const NoConfigLocationsErrorMessage = "no config locations set up"
const CantAccessFileErrorMessage = "error accessing file %v"
const FileNotFoundErrorMessage = "can't find included file %v"

type ProxyVariable struct {
	Name    string
	Type    string
	Address string
}

func (variable ProxyVariable) String() string {
	return fmt.Sprintf("var %v = '%v %v';", variable.Name, variable.Type, variable.Address)
}

type ProxyStatement struct {
	Statement string
	Proxy     *ProxyVariable
}

func (statement ProxyStatement) String() string {
	return fmt.Sprintf("%v return %v;", statement.Statement, statement.Proxy.Name)
}

type FileScope struct {
	CurrentProxy  *ProxyVariable
	PreviousScope *FileScope
}

func (scope FileScope) NewSubScope() *FileScope {
	return &FileScope{PreviousScope: &scope, CurrentProxy: scope.CurrentProxy}
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
	if proxy.CurrentScope.CurrentProxy == nil {
		panic(fmt.Errorf(NoProxyErrorMessage, match.Value))
	}

	statement := ProxyStatement{
		Statement: fmt.Sprintf("if (dnsDomainIs(h, '%v'))", match.Value),
		Proxy:     proxy.CurrentScope.CurrentProxy,
	}

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

// ReadConfig implements ProxySettings.
func (proxy *ProxyPac) ReadConfig(filename string) {
	config := proxy.getConfigByName(filename)

	if proxy.CurrentScope == nil {
		proxy.CurrentScope = &FileScope{}
	} else {
		proxy.CurrentScope = proxy.CurrentScope.NewSubScope()
	}

	for _, e := range config.Entries {
		e.EmitTo(proxy)
	}

	proxy.CurrentScope = proxy.CurrentScope.PreviousScope
}

func (proxy *ProxyPac) getConfigByName(filename string) *InputConfig {
	if len(proxy.ConfigLocations) == 0 {
		panic(errors.New(NoConfigLocationsErrorMessage))
	}

	for _, dir := range proxy.ConfigLocations {
		path := dir + "/" + filename
		if _, err := os.Stat(path); err == nil {
			return proxy.getConfigByPath(path)
		} else if !errors.Is(err, os.ErrNotExist) {
			panic(fmt.Errorf(NoProxyErrorMessage, path))
		}
	}

	panic(fmt.Errorf(FileNotFoundErrorMessage, filename))
}

func (proxy *ProxyPac) getConfigByPath(path string) *InputConfig {
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
