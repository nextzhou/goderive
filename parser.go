package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

type FileInfo struct {
	PkgName string
	Types   []TypeInfo
}

type TypeInfo struct {
	Name string
	Opts Options
}

func (ti TypeInfo) ToTemplateArgs() TemplateArgs {
	var ta TemplateArgs
	ta.TypeName = ti.Name
	if ti.Opts.Rename == "" {
		ta.SetName = ti.Name + "Set"
	} else {
		ta.SetName = ti.Opts.Rename
	}
	ta.CapitalizeSetName = capitalize(ta.SetName)
	if ti.Opts.ByPointer {
		ta.TypeName = "*" + ta.TypeName
	}
	return ta
}

type Options struct {
	Rename    string // Rename Set struct
	ByPointer bool   // Record elements by pointer
}

func (opt *Options) Parse(desc string) error {
	optStrs := strings.Split(desc, ",")
	for _, optStr := range optStrs {
		terms := strings.SplitN(optStr, "=", 2)
		switch strings.TrimSpace(terms[0]) {
		case OptionRename:
			if opt.Rename != "" {
				return fmt.Errorf("duplicate '%s' option", OptionRename)
			}
			if len(terms) == 2 {
				name := strings.TrimSpace(terms[1])
				if name != "" {
					opt.Rename = name
					break
				}
			}
			return fmt.Errorf("wrong '%s' format, expect '%s=<name>'", OptionRename, OptionRename)
		case OptionByPointer:
			if opt.ByPointer {
				return fmt.Errorf("duplicate '%s' option", OptionByPointer)
			}
			if len(terms) == 1 {
				opt.ByPointer = true
			} else {
				return fmt.Errorf("wrong '%s' format, expect '%s'", OptionByPointer, OptionByPointer)
			}
		}
	}
	return nil
}

func ExtractTypes(src []byte) (FileInfo, error) {
	var fileInfo FileInfo

	// parse source code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return fileInfo, err
	}
	pkg, _ := ast.NewPackage(fset, map[string]*ast.File{"": file}, nil, nil)
	d := doc.New(pkg, file.Name.String(), doc.AllDecls)

	// select types with 'goset' marker
	for _, typ := range d.Types {
		var typeInfo TypeInfo

		cmts := strings.Split(typ.Doc, "\n")
		for _, cmt := range cmts {
			matched, optsStr := MatchGosetComment(cmt)
			if matched {
				typeInfo.Name = typ.Name
				err = typeInfo.Opts.Parse(optsStr)
				if err != nil {
					return fileInfo, err
				}
			}
		}
		if typeInfo.Name != "" {
			fileInfo.PkgName = pkg.Name
			fileInfo.Types = append(fileInfo.Types, typeInfo)
		}
	}
	return fileInfo, nil
}

func MatchGosetComment(cmt string) (matched bool, opts string) {
	cmt = strings.TrimSpace(cmt)
	if cmt == "goset" {
		return true, ""
	} else if strings.HasPrefix(cmt, "goset:") {
		cmt = strings.TrimPrefix(cmt, "goset:")
		return true, cmt
	}
	return false, ""
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	if 'a' <= s[0] && s[0] <= 'z' {
		b := []byte(s)
		b[0] = b[0] - (byte('a') - 'A')
		return string(b)
	}
	return s
}
