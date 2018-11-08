package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"

	"github.com/nextzhou/goderive/plugin"
	"github.com/nextzhou/goderive/utils"
)

type TypeInfo struct {
	Name     string
	Assigned string
	Plugins  *plugin.Entries
	Ast      ast.Expr
	Env      plugin.Env
}

func ExtractTypes(src []byte) ([]TypeInfo, error) {
	var types []TypeInfo
	// parse source code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	pkg, _ := ast.NewPackage(fset, map[string]*ast.File{"": file}, nil, nil)
	d := doc.New(pkg, file.Name.String(), doc.AllDecls)

	env := plugin.MakeEnv(pkg.Name)

	for _, i := range file.Imports {
		env.Imports.Append(plugin.MakeImportFromAst(i))
	}

	// select types with 'derive' marker
	for _, typ := range d.Types {
		var typeInfo TypeInfo

		cmts := strings.Split(typ.Doc, "\n")
		for _, cmt := range cmts {
			dc, err := utils.MatchDeriveComment(cmt)
			if err != nil {
				return nil, fmt.Errorf("type %s: %v", typ.Name, err)
			}
			if dc == nil {
				continue
			}
			if typeInfo.Name == "" {
				typeInfo.Name = typ.Name
				typeInfo.Plugins = plugin.NewEntries(0)
				spec := typ.Decl.Specs[0].(*ast.TypeSpec)
				typeInfo.Ast = spec.Name.Obj.Decl.(*ast.TypeSpec).Type
				if spec.Assign.IsValid() {
					switch assigned := typeInfo.Ast.(type) {
					case *ast.Ident:
						typeInfo.Assigned = assigned.Name
					case *ast.SelectorExpr:
						typeInfo.Assigned = utils.SelectorExprString(assigned)
					}
				}
			}
			opts, err := plugin.ParseOptions(dc.OptionsStr)
			if err != nil {
				return nil, fmt.Errorf("type %s: %v", typ.Name, err)
			}

			// merge options
			idx := typeInfo.Plugins.FindBy(func(e plugin.Entry) bool { return e.Plugin == dc.Plugin })
			if idx == -1 {
				typeInfo.Plugins.Append(plugin.MakeEntry(dc.Plugin, opts))
			} else {
				err := typeInfo.Plugins.Index(idx).Opts.Merge(opts)
				if err != nil {
					return nil, fmt.Errorf("type %s: %v", typ.Name, err)
				}
			}
		}
		if typeInfo.Name != "" {
			typeInfo.Env = env
			types = append(types, typeInfo)
		}
	}
	return types, nil
}
