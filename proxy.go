package main

import (
	"fmt"
	"io"
)

type ProxySettings interface {
	AddProxy(ProxyConfigEntry)
	AddSubdomainMatch(SubDomainMatch)
	WriteSettings(*InputConfig, io.Writer)
}

// A Proxy auto-config file generator
// Implements ProxySettings interface
type ProxyPac struct {
	Proxies    []string
	Statements []string
}

// AddProxy implements ProxySettings.
func (proxy *ProxyPac) AddProxy(entry ProxyConfigEntry) {
	proxy.Proxies = append(proxy.Proxies, entry.Type+" "+entry.Address)
}

// AddSubdomainMatch implements ProxySettings.
func (proxy *ProxyPac) AddSubdomainMatch(match SubDomainMatch) {
	statement := fmt.Sprintf("	if (dnsDomainIs(h, '%v')) return p;", match.Value)
	proxy.Statements = append(proxy.Statements, statement)
}

func (proxy *ProxyPac) writeProxies(out io.Writer) {
	if len(proxy.Proxies) == 1 {
		io.WriteString(out, fmt.Sprintf("\tvar p = '%v';\n", proxy.Proxies[0]))
	} else {
		for index, proxy := range proxy.Proxies {
			io.WriteString(out, fmt.Sprintf("\tvar p%v = '%v';\n", index+1, proxy))
		}
	}

	io.WriteString(out, "\n")
}

// WriteSettings implements ProxySettings.
func (proxy *ProxyPac) WriteSettings(config *InputConfig, out io.Writer) {
	io.WriteString(out, "function FindProxyForURL (url, host) {\n")
	io.WriteString(out, "\tvar h = host.toLowerCase();\n")

	for _, entry := range config.Entries {
		entry.EmitTo(proxy)
	}

	proxy.writeProxies(out)

	for _, statement := range proxy.Statements {
		io.WriteString(out, statement)
		io.WriteString(out, "\n")
	}

	io.WriteString(out, "\n\treturn 'DIRECT';\n")
	io.WriteString(out, "}\n")
}
