package set

import (
	"go/ast"
	"io"

	"github.com/nextzhou/goderive/plugin"
	"github.com/nextzhou/goderive/utils"
)

type Set struct{}

var _ plugin.Plugin = Set{}

func (set Set) Describe() plugin.Description {
	return plugin.Description{
		Identity: "set",
		Effect:   "set collection",
		ValidFlags: []plugin.FlagDescription{
			{Key: "ByPoint", IsDefault: false, Effect: "store elements by pointer"},
			// TODO ThreadSafe
			//{Key: "ThreadSafe", IsDefault: false, Effect: "thread-safe implementation"},
			// TODO Order
			//{Key: "Order", IsDefault: false, Effect: "keep order"},
		},
		ValidArgs: []plugin.ArgDescription{
			{Key: "Rename", DefaultValue: nil, ValidValues: nil, AllowEmpty: true, IsMultipleValues: false, Effect: "assign set type name manually"},
		},
		AllowUnexpectedlyFlag: false,
		AllowUnexpectedlyArg:  false,
	}
}

func (set Set) GenerateTo(w io.Writer, typeName string, typeInfo ast.TypeSpec, opt plugin.Options) (plugin.Prerequisites, error) {
	var arg TemplateArgs
	var pre plugin.Prerequisites
	pre.Imports = []string{"fmt", "encoding/json"}
	arg.RawTypeName = typeName
	arg.SetName = typeName + "Set"
	if val := opt.GetValue("Rename"); !val.IsNil() {
		if !utils.ValidateIdentName(val.Str()) {
			return pre, &utils.InvalidIdentError{Type: "rename", Ident: val.Str()}
		}
		arg.SetName = val.Str()
	}
	arg.CapitalizeSetName = utils.Capitalize(arg.SetName)
	if opt.WithFlag("ByPoint") {
		arg.TypeName = "*" + arg.RawTypeName
	} else {
		arg.TypeName = arg.RawTypeName
	}

	return pre, arg.GenerateTo(w)
}
