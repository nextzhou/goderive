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
		Identity:   "set",
		Effect:     "set collection",
		ValidFlags: []plugin.FlagDescription{
		// TODO ThreadSafe
		// {Key: "ThreadSafe", IsDefault: false, Effect: "thread-safe implementation"},
		},
		ValidArgs: []plugin.ArgDescription{
			{Key: "Rename", DefaultValue: nil, ValidValues: nil, AllowEmpty: true, IsMultipleValues: false, Effect: "assign set type name manually"},
			{Key: "Order", DefaultValue: &UnstableOrder, ValidValues: plugin.NewValueSetFromSlice([]plugin.Value{UnstableOrder, AppendOrder, KeyOrder}),
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
	arg.TypeName = typeInfo.Name
	arg.SetName = typeInfo.Name + "Set"
	if val := opt.GetValue("Rename"); !val.IsNil() {
		if !utils.ValidateIdentName(val.Str()) {
			return pre, &utils.InvalidIdentError{Type: "rename", Ident: val.Str()}
		}
		arg.SetName = val.Str()
	}
	arg.Order = string(opt.MustGetValue("Order"))
	if arg.Order == KeyOrder.Str() {
		pre.Imports = append(pre.Imports, "sort")
	}
	arg.CapitalizeSetName = utils.Capitalize(arg.SetName)
	arg.IsComparable = utils.IsComparableType(typeInfo.Assigned)

	return pre, arg.GenerateTo(w)
}

var (
	UnstableOrder = plugin.Value("Unstable")
	AppendOrder   = plugin.Value("Append")
	KeyOrder      = plugin.Value("Key")
)
