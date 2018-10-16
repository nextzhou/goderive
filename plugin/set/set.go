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
			{Key: "Export", Default: utils.TriBoolUndefined, Effect: "force the generated code to be exported/unexported"},
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

func (set Set) GenerateTo(w io.Writer, env plugin.Env, typeInfo plugin.TypeInfo, opt plugin.Options) (plugin.Prerequisites, error) {
	var arg TemplateArgs
	pre := plugin.MakePrerequisites()
	forceExport := opt.GetFlag("Export")
	pre.Imports.Append(plugin.MakeImport("fmt"),
		plugin.MakeImport("encoding/json"),
		plugin.MakeImport("reflect"))
	arg.TypeName = typeInfo.Name

	// use assigned type as type name
	if typeInfo.Assigned != "" {
		i := env.SelectImportForType(typeInfo.Assigned)
		if i != nil {
			arg.TypeName = typeInfo.Assigned
			pre.Imports.Append(*i)
		}
	}

	if forceExport.UnwrapOr(utils.IsExported(typeInfo.Name)) {
		arg.SetName = utils.ToExported(typeInfo.Name) + "Set"
	} else {
		arg.SetName = utils.ToUnexported(typeInfo.Name) + "Set"
	}

	if val := opt.GetValue("Rename"); !val.IsNil() {
		if !utils.ValidateIdentName(val.Str()) {
			return pre, &utils.InvalidIdentError{Type: "Rename", Ident: val.Str()}
		}
		arg.SetName = val.Str()
	}

	if forceExport.UnwrapOr(utils.IsExported(arg.SetName)) {
		arg.New = "New"
	} else {
		arg.New = "new"
	}

	arg.Order = string(opt.MustGetValue("Order"))
	if arg.Order == KeyOrder.Str() {
		pre.Imports.Append(plugin.MakeImport("sort"))
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
