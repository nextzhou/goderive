package access

import (
	"go/ast"
	"io"
	"strings"

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
			{Key: "Receiver", DefaultValue: nil, ValidValues: nil, AllowEmpty: true, IsMultipleValues: false, Effect: "receiver of methods"},
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
	args.Receiver = strings.ToLower(typeInfo.Name[:1])
	if r := opt.GetValue("Receiver"); !r.IsNil() {
		args.Receiver = r.Str()
	}
	if !utils.ValidateIdentName(args.Receiver) {
		return pre, &utils.InvalidIdentError{Type: "Receiver", Ident: args.Receiver}
	}

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
		var cmtList []*ast.Comment
		if field.Doc != nil {
			cmtList = field.Doc.List
		}

		// field options parse
		var fieldOpts = plugin.NewOptions()
		for _, cmt := range cmtList {
			dc, err := utils.MatchPluginComment(cmt.Text)
			if err != nil {
				return pre, err
			}
			opts, err := plugin.ParseOptions(dc.OptionsStr)
			if err != nil {
				return pre, err
			}
			if err = fieldOpts.Merge(opts); err != nil {
				return pre, err
			}
		}

		args.AddField(field, fieldOpts, t)
	}
	return pre, args.GenerateTo(w)
}

func genGetName(field string) string {
	if utils.IsExported(field) {
		return "Get" + field
	}
	return "get" + utils.Capitalize(field)
}
