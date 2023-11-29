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
	Proxies         []*ProxyVariable
	Statements      []ProxyStatement
	CurrentProxy    *ProxyVariable
}

const NoProxyErrorMessage = "no proxy to use for %v"
const NoConfigLocationsErrorMessage = "no config locations set up"

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

// UseProxy implements ProxySettings.
func (proxy *ProxyPac) UseProxy(entry ProxyConfigEntry) {
	for _, p := range proxy.Proxies {
		if p.Address == entry.Address {
			proxy.CurrentProxy = p
			return
		}
	}

	proxy.CurrentProxy = &ProxyVariable{Type: entry.Type, Address: entry.Address}

	if len(proxy.Proxies) == 0 {
		proxy.CurrentProxy.Name = "p"
	} else if len(proxy.Proxies) == 1 {
		proxy.Proxies[0].Name = "p1"
		proxy.CurrentProxy.Name = "p2"
	} else {
		proxy.CurrentProxy.Name = fmt.Sprintf("p%v", len(proxy.Proxies)+1)
	}

	proxy.Proxies = append(proxy.Proxies, proxy.CurrentProxy)
}

// AddSubdomainMatch implements ProxySettings.
func (proxy *ProxyPac) AddSubdomainMatch(match SubDomainMatch) {
	if proxy.CurrentProxy == nil {
		panic(fmt.Errorf(NoProxyErrorMessage, match.Value))
	}

	statement := ProxyStatement{
		Statement: fmt.Sprintf("if (dnsDomainIs(h, '%v'))", match.Value),
		Proxy:     proxy.CurrentProxy,
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
	config := proxy.getConfig(filename)
	for _, e := range config.Entries {
		e.EmitTo(proxy)
	}
}

func (proxy *ProxyPac) getConfig(filename string) *InputConfig {
	if len(proxy.ConfigLocations) == 0 {
		panic(errors.New(NoConfigLocationsErrorMessage))
	}

	path := proxy.ConfigLocations[0] + "/" + filename
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
