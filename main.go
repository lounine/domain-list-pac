package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type arrayFlags []string

var (
	input  string
	output string
	data   arrayFlags
)

func (array *arrayFlags) String() string {
	return strings.Join(*array, ", ")
}

func (array *arrayFlags) Set(value string) error {
	*array = append(*array, value)
	return nil
}

func main() {
	flag.Var(&data, "d", "Additional directories to lookup for included configs, may be repeated")
	flag.StringVar(&output, "o", "./proxy.pac", "Save generated file as")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage:")
		fmt.Fprintln(flag.CommandLine.Output(), "domain-list-pac | go run . <input-file>")
		fmt.Fprintln(flag.CommandLine.Output(), "domain-list-pac | go run . [-d <data-dir>]... [-o <output-file>] <input-file>")
		fmt.Fprintln(flag.CommandLine.Output(), "\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	input = flag.Args()[0]

	proxy := NewProxyPac(input, data)
	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			panic(err)
		}
	}()

	proxy.WriteSettings(out)
	fmt.Fprintf(os.Stdout, "Saved config to: %v\n", output)
}
