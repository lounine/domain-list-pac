package main

import (
	"os"
)

func main() {
	config, err := parser.Parse("", os.Stdin)
	if err != nil {
		panic(err)
	}

	proxy := ProxyPac{}
	proxy.ReadConfig(config)
	proxy.WriteSettings(os.Stdout)
}
