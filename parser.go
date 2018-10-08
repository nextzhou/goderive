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
	// TODO keep order
	Plugins map[string]*plugin.Options
	Ast     ast.TypeSpec
	Env     plugin.Env
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
		env.Imports.Append(plugin.MakeInportFromAst(i))
	}

	// select types with 'derive' marker
	for _, typ := range d.Types {
		var typeInfo TypeInfo

		cmts := strings.Split(typ.Doc, "\n")
		for _, cmt := range cmts {
			dc, err := MatchDeriveComment(cmt)
			if err != nil {
				return nil, fmt.Errorf("type %s: %v", typ.Name, err)
			}
			if dc == nil {
				continue
			}
			if typeInfo.Name == "" {
				typeInfo.Name = typ.Name
				typeInfo.Plugins = make(map[string]*plugin.Options)
				typeInfo.Ast = *typ.Decl.Specs[0].(*ast.TypeSpec)
				if typeInfo.Ast.Assign.IsValid() {
					assignType := typeInfo.Ast.Name.Obj.Decl.(*ast.TypeSpec).Type
					switch assignType.(type) {
					case *ast.Ident:
						typeInfo.Assigned = assignType.(*ast.Ident).Name
					case *ast.SelectorExpr:
						typeInfo.Assigned = utils.SelectorExprString(assignType.(*ast.SelectorExpr))
					}
				}
			}
			opts, err := plugin.ParseOptions(dc.OptionsStr)
			if err != nil {
				return nil, fmt.Errorf("type %s: %v", typ.Name, err)
			}

			// merge options
			currentOpts := typeInfo.Plugins[dc.Plugin]
			if currentOpts.IsEmpty() {
				currentOpts = opts
			} else {
				err = currentOpts.Merge(opts)
				if err != nil {
					return nil, fmt.Errorf("type %s: %v", typ.Name, err)
				}
			}
			typeInfo.Plugins[dc.Plugin] = currentOpts
		}
		if typeInfo.Name != "" {
			typeInfo.Env = env
			types = append(types, typeInfo)
		}
	}
	return types, nil
}

type DeriveComment struct {
	Plugin     string
	OptionsStr string
}

func MatchDeriveComment(cmt string) (*DeriveComment, error) {
	cmt = strings.TrimSpace(cmt)
	if !strings.HasPrefix(cmt, "derive-") {
		return nil, nil
	}
	cmt = strings.TrimPrefix(cmt, "derive-")
	splitIdx := strings.Index(cmt, ":")
	dc := new(DeriveComment)
	if splitIdx == -1 {
		dc.Plugin = cmt
	} else {
		dc.Plugin = cmt[:splitIdx]
		dc.OptionsStr = strings.TrimSpace(cmt[splitIdx+1:])
	}
	if !utils.ValidateIdentName(dc.Plugin) {
		return nil, fmt.Errorf("invalid plugin name %#v", dc.Plugin)
	}
	return dc, nil
}
