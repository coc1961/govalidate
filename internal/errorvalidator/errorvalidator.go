package errorvalidator

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/coc1961/govalidate/internal/platform/dirs"
)

type ErrorStatus struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

func ValidateErrors(filePath string, skipPackages ...string) ([]ErrorStatus, error) {
	ret := []ErrorStatus{}
	dirs.TravelDirs(filePath, func(fullPath string) error {
		if strings.HasSuffix(fullPath, ".go") && !strings.HasSuffix(fullPath, "_test.go") {
			msg, err := processFile(fullPath, skipPackages...)
			if err != nil {
				return err
			}
			if len(msg) > 0 {
				ret = append(ret, ErrorStatus{
					File:  fullPath,
					Error: msg,
				})
			}
		}
		return nil
	})
	return ret, nil
}

func processFile(filePath string, skipPackages ...string) (string, error) {
	for _, s := range skipPackages {
		tmp := "/" + s + "/"
		s = strings.ReplaceAll(tmp, "//", "/")
		if strings.Contains(filePath, s) {
			return "", nil
		}
	}
	_, fileName := path.Split(filePath)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	arr := strings.Split(string(b), "\n")
	comment := false
	function := false
	llave := 0
	by := bytes.NewBufferString("")
	for i, s := range arr {
		s := strings.TrimSpace(s)
		if strings.HasPrefix(s, "/*") {
			comment = true
			continue
		}
		if strings.HasSuffix(s, "*/") {
			comment = true
			continue
		}
		if comment {
			continue
		}
		if strings.Index(s, "//") == 0 {
			continue
		}

		if strings.HasPrefix(s, "func ") {
			function = true
		}

		if function {
			llave = strings.Count(s, "{")
			tmp, err := processLine(s, fileName, i+1)
			if err != nil {
				return "", err
			}
			by.WriteString(tmp)
			llave -= strings.Count(s, "}")
			if llave == 0 {
				function = false
			}
		}

	}

	return by.String(), nil
}

func has(s, x string) bool {
	com := false
	for i := range s {
		if s[i] == '\'' || s[i] == '"' || s[i] == '`' {
			com = !com
		}
		if com {
			continue
		}

		if string(s[i:]) == "//" {
			return false
		}
		if string(s[i:]) == "/*" {
			return false
		}
		if strings.Index(string(s[i:]), x) == 0 {
			return true
		}
	}

	return false
}

func processLine(line, filename string, lineNumber int) (string, error) {
	e := ""
	if has(line, `fmt.Errorf(`) {
		e += "File: " + filename + " contain fmt.Errorf(" + fmt.Sprint(lineNumber) + ")\n"
	}
	if has(line, `errors.New(`) {
		e += "File:" + filename + " contain errors.New(" + fmt.Sprint(lineNumber) + ")\n"
	}
	return e, nil
}
