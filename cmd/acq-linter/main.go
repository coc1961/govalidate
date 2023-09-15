package main

import (
	"github.com/coc1961/govalidate/cmd/func-linter/pkg/valerror"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(valerror.Analyzer)
}
