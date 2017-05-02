package util

import (
	"os"
	"path/filepath"
)

func StringToFile(s string, filename string) error {
	// ensure the directory of the file exists
	dirpath := filepath.Dir(filename)
	err := os.MkdirAll(dirpath, 755)
	if err != nil {
		return err
	}

	// now open the file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// write the file contents
	f.WriteString(s)
	return nil
}
