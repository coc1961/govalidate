package main

import (
	"github.com/coc1961/govalidate/cmd/func-linter/pkg/valerror"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(valerror.Analyzer) }
