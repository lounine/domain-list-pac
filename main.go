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
	proxyPac.WriteSettings(config, os.Stdout)

	//config.EmitProxy(os.Stdout)

	// repr.Println(config, repr.Indent("  "), repr.OmitEmpty(true))

}
