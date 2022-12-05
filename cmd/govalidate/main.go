package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coc1961/govalidate/internal/importvalidator"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

type manyStrings []string

func (i *manyStrings) String() string {
	return "my string representation"
}

func (i *manyStrings) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var excludePackage, src manyStrings
	flag.Var(&src, "p", "project root path")
	verbose := flag.Bool("v", false, "verbose, print import erros")
	flag.Var(&excludePackage, "e", "exclude packages (optional)")

	flag.Parse()

	if len(src) == 0 {
		flag.CommandLine.Usage()
		return
	}

	ret := []importvalidator.ImportsStatus{}
	for _, filePath := range src {
		tmp, err := importvalidator.ValidateImports(filePath, excludePackage...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ret = append(ret, tmp...)
	}

	if *verbose {
		for _, e := range ret {
			fmt.Println(Blue, "File: ", Gray, e.File)
			for _, i := range e.InvalidImports {
				fmt.Println(Yellow, "\tInvalid Import: ", Gray, i)
			}
			fmt.Println(Reset)
		}
	}
	if len(ret) > 0 {
		os.Exit(1)
	}
}
