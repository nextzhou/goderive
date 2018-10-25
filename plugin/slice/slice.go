package slice

import (
	"io"

	"github.com/nextzhou/goderive/plugin"
	"github.com/nextzhou/goderive/utils"
)

type Slice struct{}

var _ plugin.Plugin = Slice{}

func (s Slice) Describe() plugin.Description {
	return plugin.Description{
		Identity: "slice",
		Effect:   "slice extension",
		ValidFlags: []plugin.FlagDescription{
			{Key: "Export", Default: utils.TriBoolUndefined, Effect: "force the generated code to be exported/unexported"},
		},
		ValidArgs: []plugin.ArgDescription{
			{Key: "Rename", DefaultValue: nil, ValidValues: nil, AllowEmpty: true, IsMultipleValues: false, Effect: "assign slice type name manually"},
		},
		AllowUnexpectedlyFlag: false,
		AllowUnexpectedlyArg:  false,
	}
}

func (s Slice) GenerateTo(w io.Writer, env plugin.Env, typeInfo plugin.TypeInfo, opt plugin.Options) (plugin.Prerequisites, error) {
	var arg TemplateArgs
	pre := plugin.MakePrerequisites()
	forceExport := opt.GetFlag("Export")
	pre.Imports.Append(plugin.MakeImport("fmt"),
		plugin.MakeImport("encoding/json"),
		plugin.MakeImport("reflect"))
	arg.TypeName = typeInfo.Name

	if typeInfo.Assigned != "" {
		i := env.SelectImportForType(typeInfo.Assigned)
		if i != nil {
			arg.TypeName = typeInfo.Assigned
			pre.Imports.Append(*i)
		}
	}

	if forceExport.UnwrapOr(utils.IsExported(typeInfo.Name)) {
		arg.SliceName = utils.ToExported(typeInfo.Name) + "Slice"
	} else {
		arg.SliceName = utils.ToUnexported(typeInfo.Name) + "Slice"
	}

	if val := opt.GetValue("Rename"); !val.IsNil() {
		if !utils.ValidateIdentName(val.Str()) {
			return pre, &utils.InvalidIdentError{Type: "Rename", Ident: val.Str()}
		}
		arg.SliceName = val.Str()
	}

	if forceExport.UnwrapOr(utils.IsExported(arg.SliceName)) {
		arg.New = "New"
	} else {
		arg.New = "new"
	}

	arg.CapitalizeSliceName = utils.Capitalize(arg.SliceName)
	arg.IsSortable = utils.IsSortableType(typeInfo.Assigned)
	arg.IsComparable = utils.IsComparableType(typeInfo.Ast)
	return pre, arg.GenerateTo(w)
}
