package importvalidator

import (
	"go/parser"
	"go/token"
	"path"
	"strings"

	"github.com/coc1961/govalidate/internal/platform/dirs"
)

var fs = token.NewFileSet()

type ImportsStatus struct {
	File           string   `json:"file"`
	InvalidImports []string `json:"invalid_imports"`
}

func ValidateImports(filePath string, skipPackages ...string) ([]ImportsStatus, error) {
	ret := []ImportsStatus{}

	dirs.TravelDirs(filePath, func(fullPath string) error {
		if strings.HasSuffix(fullPath, ".go") {
			arr, err := processFile(fullPath, skipPackages...)
			if err != nil {
				return err
			}
			if len(arr) > 0 {
				ret = append(ret, ImportsStatus{
					File:           fullPath,
					InvalidImports: arr,
				})
			}
		}
		return nil
	})
	return ret, nil
}

func processFile(filePath string, skipPackages ...string) ([]string, error) {
	arr := []string{}
	for _, s := range skipPackages {
		tmp := "/" + s + "/"
		s = strings.ReplaceAll(tmp, "//", "/")
		if strings.Contains(filePath, s) {
			return arr, nil
		}
	}

	dir, _ := path.Split(filePath)
	basePackage, _ := path.Split(strings.TrimRight(dir, "/"))

	f, err := parser.ParseFile(fs, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, dec := range f.Imports {
		pref, _ := path.Split(strings.ReplaceAll(dec.Path.Value, `"`, ""))
		if pref == "" {
			continue
		}
		if strings.HasSuffix(basePackage, pref) {
			arr = append(arr, strings.ReplaceAll(dec.Path.Value, `"`, ""))
		}
	}

	return arr, nil
}
