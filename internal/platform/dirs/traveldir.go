package dirs

import (
	"os"
	"path"
	"path/filepath"
)

func TravelDirs(filePath string, fn func(fullPath string) error) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	file, err := os.ReadDir(filePath)
	if err != nil {
		return err
	}

	for _, f := range file {
		fullPath := path.Join(filePath, f.Name())
		if f.IsDir() {
			err := TravelDirs(fullPath, fn)
			if err != nil {
				return err
			}
		} else {
			err := fn(fullPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
