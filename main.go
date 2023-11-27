package main

import (
	"os"
)

func main() {
	config, err := parser.Parse("", os.Stdin)
	if err != nil {
		panic(err)
	}

	proxyPac := ProxyPac{}
	proxyPac.ReadConfig(config)
	proxyPac.WriteSettings(os.Stdout)
}
