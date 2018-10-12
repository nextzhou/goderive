package access

import (
	"go/ast"
	"io"

	"github.com/nextzhou/goderive/plugin"
	"github.com/nextzhou/goderive/utils"
)

type Access struct{}

var _ plugin.Plugin = Access{}

func (a Access) Describe() plugin.Description {
	return plugin.Description{
		Identity:   "access",
		Effect:     "access fields for struct type",
		ValidFlags: []plugin.FlagDescription{},
		ValidArgs: []plugin.ArgDescription{
			{Key: "Receiver", DefaultValue: nil, ValidValues: nil, AllowEmpty: false, IsMultipleValues: false, Effect: "receiver of methods"},
		},
		AllowUnexpectedlyFlag: false,
		AllowUnexpectedlyArg:  false,
	}
}

func (a Access) GenerateTo(w io.Writer, env plugin.Env, typeInfo plugin.TypeInfo, opt plugin.Options) (plugin.Prerequisites, error) {
	var args TemplateArgs
	pre := plugin.MakePrerequisites()

	s, ok := typeInfo.Ast.(*ast.StructType)
	if !ok {
		return pre, &utils.OnlySupportError{Supported: "Struct", Got: utils.ExprTypeStr(typeInfo.Ast)}
	}

	args.TypeName = typeInfo.Name
	r := opt.GetValue("Receiver").Str()
	if !utils.ValidateIdentName(r) {
		return pre, &utils.InvalidIdentError{Type: "Receiver", Ident: r}
	}
	args.Receiver = r

	for _, field := range s.Fields.List {
		t := utils.TypeNameWithPkg(field.Type)
		if t == nil {
			continue
		}
		t.Pkgs.ForEach(func(s string) {
			found := env.Imports.FindBy(func(i plugin.Import) bool {
				pkgName, ok := i.PkgName()
				return ok && pkgName == s
			})
			if found != nil {
				pre.Imports.Append(*found)
			}
		})
		for _, name := range field.Names {
			args.Fields = append(args.Fields, Field{Name: name.Name, GetFuncName: genGetName(name.Name), TypeName: t.Name})
		}
	}
	return pre, args.GenerateTo(w)
}

func genGetName(field string) string {
	if utils.IsExported(field) {
		return "Get" + field
	}
	return "get" + utils.Capitalize(field)
}
