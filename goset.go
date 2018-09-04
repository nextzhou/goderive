package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	args := ParseArgs()

	var files []string
	for _, path := range args.InputPaths {
		fs, err := ListGoFiles(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		files = append(files, fs...)
	}

	groupFileInfoByPath := make(map[string]*FileInfo)
	for _, file := range files {
		path, err := filepath.Abs(filepath.Dir(file))
		if err != nil {
			panic(err)
		}
		fi, ok := groupFileInfoByPath[path]
		if !ok {
			groupFileInfoByPath[path] = new(FileInfo)
		}

		src, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %#v : %s\n", file, err.Error())
			os.Exit(1)
		}
		fileInfo, err := ExtractTypes(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:%s\n", file, err.Error())
			os.Exit(1)
		}
		if fileInfo.PkgName == "" || len(fileInfo.Types) == 0 {
			continue
		}
		if fi.PkgName == "" {
			fi.PkgName = fileInfo.PkgName
		}
		fi.Types = append(fi.Types, fileInfo.Types...)
		groupFileInfoByPath[path] = fi
	}

	for path, fileInfo := range groupFileInfoByPath {
		filename := filepath.Join(path, args.OutputFilename)
		if fileInfo.PkgName == "" {
			os.Remove(filename)
			continue
		}
		buf := bytes.NewBuffer(nil)
		buf.WriteString(fmt.Sprintf("package %s\n\n", fileInfo.PkgName))
		for _, typ := range fileInfo.Types {
			err := typ.ToTemplateArgs().GenerateTo(buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to generate code of %#v: %s\n", typ.Name, err.Error())
				os.Exit(1)
			}
		}
		err := ioutil.WriteFile(filename, buf.Bytes(), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write %#v error : %s\n", filename, err.Error())
			os.Exit(1)
		}
	}
}

type Args struct {
	OutputFilename string
	InputPaths     []string
}

func ParseArgs() Args {
	var args Args
	flag.StringVar(&args.OutputFilename, "o", "gen_goset.go", "output filename")
	flag.Parse()
	args.InputPaths = flag.Args()
	if len(args.InputPaths) == 0 {
		args.InputPaths = []string{"."}
	}
	return args
}
