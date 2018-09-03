package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func fileList(path string) ([]string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var files []string
	if stat.IsDir() {
		dirInfo, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, entry := range dirInfo {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
	} else {
		if strings.HasSuffix(path, ".go") {
			files = []string{path}
		}
	}
	return files, nil
}
