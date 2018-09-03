package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	fs, err := fileList(os.Args[1])
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	var writePkgName func(string)
	writePkgName = func(pkg string) {
		buf.WriteString("package ")
		buf.WriteString(pkg)
		buf.WriteByte('\n')
		writePkgName = func(_ string) {}
	}
	var filename string

	for _, f := range fs {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		pkgName, typs, err := parse(f, data)
		if len(typs) == 0 {
			continue
		}
		writePkgName(pkgName)
		if filename == "" {
			filename = filepath.Join(filepath.Dir(f), "gen_goset.go")
		}
		if err != nil {
			panic(err)
		}
		for _, typ := range typs {
			err := tmpl.Execute(buf, typ)
			if err != nil {
				panic(err)
			}
		}
	}
	code := buf.Bytes()
	if len(code) > 0 && filename != "" {
		err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}
}
