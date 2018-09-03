package main

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

type SetType struct {
	Name string
	opts []setTypeOption
}

type setTypeOption struct {
	key string
	val string
}

func parse(path string, src []byte) (string, []SetType, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	pkg, _ := ast.NewPackage(fset, map[string]*ast.File{path: file}, nil, nil)

	d := doc.New(pkg, file.Name.String(), doc.AllDecls)
	var typs []SetType
	for _, typ := range d.Types {
		var opts []setTypeOption
		marked := false
		cmts := strings.Split(typ.Doc, "\n")
		for _, cmt := range cmts {
			cmt = strings.TrimSpace(cmt)
			if cmt == "goset" {
				marked = true
			} else if strings.HasPrefix(cmt, "goset:") {
				marked = true
				cmt = strings.TrimPrefix(cmt, "goset:")
				optStrs := strings.Split(cmt, ",")
				for _, optStr := range optStrs {
					optTerms := strings.SplitN(optStr, "=", 2)
					var opt setTypeOption
					opt.key = optTerms[0]
					if len(optTerms) > 1 {
						opt.val = optTerms[1]
					}
					opts = append(opts, opt)
				}
			}
		}
		if marked {
			typs = append(typs, SetType{Name: typ.Name, opts: opts})
		}
	}
	return pkg.Name, typs, nil
}
