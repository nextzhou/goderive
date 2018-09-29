package set

import (
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
		},
		ValidArgs: []plugin.ArgDescription{
			{Key: "Rename", DefaultValue: nil, ValidValues: nil, AllowEmpty: true, IsMultipleValues: false, Effect: "assign set type name manually"},
			{Key: "Order", DefaultValue: &UnstableOrder, ValidValues: plugin.NewValueSetFromSlice([]plugin.Value{UnstableOrder, AppendOrder}),
				AllowEmpty: true, IsMultipleValues: false, Effect: "keep order"},
		},
		AllowUnexpectedlyFlag: false,
		AllowUnexpectedlyArg:  false,
	}
}

func (set Set) GenerateTo(w io.Writer, typeInfo plugin.TypeInfo, opt plugin.Options) (plugin.Prerequisites, error) {
	var arg TemplateArgs
	var pre plugin.Prerequisites
	pre.Imports = []string{"fmt", "encoding/json"}
	arg.RawTypeName = typeInfo.Name
	arg.SetName = typeInfo.Name + "Set"
	if val := opt.GetValue("Rename"); !val.IsNil() {
		if !utils.ValidateIdentName(val.Str()) {
			return pre, &utils.InvalidIdentError{Type: "rename", Ident: val.Str()}
		}
		arg.SetName = val.Str()
	}
	arg.Order = string(opt.MustGetValue("Order"))
	arg.CapitalizeSetName = utils.Capitalize(arg.SetName)
	if opt.WithFlag("ByPoint") {
		arg.TypeName = "*" + arg.RawTypeName
	} else {
		arg.TypeName = arg.RawTypeName
	}

	return pre, arg.GenerateTo(w)
}

// TODO more sort type
var (
	UnstableOrder = plugin.Value("Unstable")
	AppendOrder   = plugin.Value("Append")
)
