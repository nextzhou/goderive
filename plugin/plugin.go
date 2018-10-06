//go:generate goderive
package plugin

import (
	"bytes"
	"fmt"
	"go/ast"
	"io"
	"strings"

	"github.com/nextzhou/goderive/utils"
)

const (
	OptionSep   = ";"
	ArgSep      = "="
	ArgValueSep = ","
)

type OptionType string

const (
	OptionTypeNone     = ""
	OptionTypeFlag     = "flag"
	OptionTypeArgKey   = "arg key"
	OptionTypeArgValue = "arg value"
)

func (ot OptionType) IsUnique() bool {
	return ot == OptionTypeFlag || ot == OptionTypeArgKey
}

func (ot OptionType) IsIdent() bool {
	return ot == OptionTypeFlag || ot == OptionTypeArgKey
}

// TODO support base plugin, e.g. plugin to define collection interface for other plugins
// TODO pop operation
// TODO arg with default value as flag
type Options struct {
	// TODO keep order
	Flags          map[Flag]utils.TriBool
	Args           map[string]Arg
	ExistingOption map[string]OptionType
}

func NewOptions() *Options {
	return &Options{
		Flags:          make(map[Flag]utils.TriBool),
		Args:           make(map[string]Arg),
		ExistingOption: make(map[string]OptionType),
	}
}

func (opts *Options) SetFlag(flag string, val utils.TriBool) {
	opts.Flags[Flag(flag)] = val
	opts.ExistingOption[flag] = OptionTypeFlag
}

func (opts *Options) SetArg(arg Arg) {
	opts.Args[arg.Key] = arg
	opts.ExistingOption[arg.Key] = OptionTypeArgKey
}

func (opts Options) ValidateOption(optType OptionType, opt string) error {
	if optType.IsIdent() && !utils.ValidateIdentName(opt) {
		return &utils.InvalidIdentError{Type: string(optType), Ident: opt}
	}
	if t := opts.ExistingOption[opt]; optType.IsUnique() && t != OptionTypeNone {
		return &utils.ConflictingOptionError{Type: string(t), Ident: opt}
	}
	return nil
}

func ParseOptions(optsStr string) (*Options, error) {
	ret := NewOptions()

	optsStr = strings.TrimSpace(optsStr)
	if optsStr == "" {
		return ret, nil
	}

	opts := strings.Split(optsStr, OptionSep)
	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		sepIdx := strings.Index(opt, ArgSep)
		// flag
		if sepIdx == -1 {
			flagVal := true
			if len(opt) > 0 && opt[0] == '!' {
				flagVal = false
				opt = strings.TrimSpace(opt[1:])
			}
			if err := ret.ValidateOption(OptionTypeFlag, opt); err != nil {
				return nil, err
			}
			ret.SetFlag(opt, utils.BoolToTri(flagVal))
			continue
		}
		// arg
		var arg Arg
		key, valsStr := opt[:sepIdx], opt[sepIdx+len(ArgSep):]
		key = strings.TrimSpace(key)

		if err := ret.ValidateOption(OptionTypeArgKey, key); err != nil {
			return nil, err
		}

		arg.Key = key
		vals := strings.Split(valsStr, ArgValueSep)
		for _, val := range vals {
			val = strings.TrimSpace(val)
			if err := ret.ValidateOption(OptionTypeArgValue, val); err != nil {
				return nil, err
			}
			arg.Values = append(arg.Values, Value(val))
		}
		ret.SetArg(arg)
	}
	return ret, nil
}

func (opts *Options) Merge(another *Options) error {
	if another == nil {
		return nil
	}
	for flag, val := range another.Flags {
		if err := opts.ValidateOption(OptionTypeFlag, string(flag)); err != nil {
			return err
		}
		opts.SetFlag(string(flag), val)
	}

	for key, arg := range another.Args {
		if err := opts.ValidateOption(OptionTypeArgKey, key); err != nil {
			return err
		}
		opts.SetArg(arg)
	}
	return nil
}

func (opts *Options) IsEmpty() bool {
	return opts == nil || len(opts.Flags)+len(opts.Args) == 0
}

type Flag string

type Arg struct {
	Key    string
	Values []Value
}

// derive-set:Order=Append
type Value string

func (v *Value) IsNil() bool {
	return v == nil
}

func (v *Value) Str() string {
	if v == nil {
		return ""
	}
	return string(*v)
}

func (arg Arg) GetSingleValue() *Value {
	if len(arg.Values) == 0 {
		return nil
	} else if len(arg.Values) != 1 {
		panic((error)(&utils.ArgNotSingleValueError{ArgKey: arg.Key}))
	}
	val := new(Value)
	*val = arg.Values[0]
	return val
}

func (arg Arg) MustGetSingleValue() Value {
	if len(arg.Values) != 1 {
		panic((error)(&utils.ArgNotSingleValueError{ArgKey: arg.Key}))
	}
	return arg.Values[0]
}

func (opts Options) GetFlag(flag Flag) utils.TriBool {
	return opts.Flags[flag]
}

func (opts Options) WithFlag(flag Flag) bool {
	return opts.Flags[flag].IsTrue()
}

func (opts Options) WithNegativeFlag(flag Flag) bool {
	return opts.Flags[flag].IsFalse()
}

func (opts Options) MustGetValue(key string) Value {
	arg, ok := opts.Args[key]
	if !ok {
		panic((error)(&utils.NotExistedError{Type: OptionTypeArgKey, Ident: key}))
	}
	return arg.MustGetSingleValue()
}
func (opts Options) GetValue(key string) *Value {
	return opts.Args[key].GetSingleValue()
}

func (opts Options) GetValues(key string) ([]Value, error) {
	arg, ok := opts.Args[key]
	if !ok {
		return nil, &utils.NotExistedError{Type: OptionTypeArgKey, Ident: key}
	}
	if len(arg.Values) == 0 {
		return nil, &utils.ArgEmptyValueError{ArgKey: arg.Key}
	}
	return arg.Values, nil
}

func (opts Options) GetValuesOrEmpty(key string) []Value {
	return opts.Args[key].Values
}

// derive-set: Order=Append
type Plugin interface {
	Describe() Description
	GenerateTo(w io.Writer, typeInfo TypeInfo, opt Options) (Prerequisites, error)
}

type Prerequisites struct {
	Imports []string
	// TODO depending plugins
}

type Description struct {
	Identity              string
	Effect                string
	ValidFlags            []FlagDescription
	ValidArgs             []ArgDescription
	AllowUnexpectedlyFlag bool
	AllowUnexpectedlyArg  bool
}

type FlagDescription struct {
	Key     string
	Default utils.TriBool
	Effect  string
}

type ArgDescription struct {
	Key string
	//ValueDescription string
	DefaultValue     *Value
	ValidValues      *ValueSet
	AllowEmpty       bool
	IsMultipleValues bool
	Effect           string
}

func (desc Description) ToHelpString() string {
	help := bytes.NewBufferString(fmt.Sprintf("Plugin: %s\n\n", desc.Identity))

	help.WriteString(desc.Effect + "\n\n")

	// flags
	if len(desc.ValidFlags) > 0 || desc.AllowUnexpectedlyFlag {
		w := utils.NewTableWriter(help)
		help.WriteString("Flags:\n")
		for _, flag := range desc.ValidFlags {
			if flag.Default.IsTrue() {
				w.Append([]string{flag.Key, flag.Effect + "(default true)"})
			} else {
				w.Append([]string{flag.Key, flag.Effect})
			}
		}
		if desc.AllowUnexpectedlyArg {
			w.Append([]string{"other flags"})
		}
		w.Render()
		help.WriteByte('\n')
	}

	// args
	if len(desc.ValidArgs) > 0 || desc.AllowUnexpectedlyArg {
		w := utils.NewTableWriter(help)
		help.WriteString("Args:\n")
		for _, arg := range desc.ValidArgs {
			var valNum string
			if arg.IsMultipleValues {
				if arg.AllowEmpty {
					valNum = "multiple values"
				} else {
					valNum = "one or more values"
				}
			} else {
				valNum = "single value"
			}
			effect := arg.Effect
			if val := arg.DefaultValue; val != nil {
				effect += "(default: " + string(*val) + ")"
			}
			var validValues string
			if !arg.ValidValues.IsEmpty() {
				validValues = arg.ValidValues.String()
			}
			w.Append([]string{arg.Key, valNum, validValues, effect})
		}
		if desc.AllowUnexpectedlyArg {
			w.Append([]string{"other args"})
		}
		w.Render()
	}
	return help.String()
}

func (desc Description) Validate(opts *Options) error {
	if err := desc.validateFlags(opts); err != nil {
		return err
	}
	if err := desc.validateArgs(opts); err != nil {
		return err
	}
	return nil
}

func (desc Description) validateFlags(opts *Options) error {
	uncheckedFlags := make(map[string]bool)
	for flag := range opts.Flags {
		uncheckedFlags[string(flag)] = true
	}
	for _, flag := range desc.ValidFlags {
		delete(uncheckedFlags, flag.Key)
		// set default value
		if _, ok := opts.Flags[Flag(flag.Key)]; !ok {
			opts.Flags[Flag(flag.Key)] = flag.Default
		}
	}

	if len(uncheckedFlags) > 0 && !desc.AllowUnexpectedlyFlag {
		flags := make([]string, 0, len(uncheckedFlags))
		for flag := range uncheckedFlags {
			flags = append(flags, flag)
		}
		return &utils.UnexpectedError{Type: string(OptionTypeFlag), Idents: flags}
	}
	return nil
}

func (desc Description) validateArgs(opts *Options) error {
	uncheckedArgs := make(map[string]bool)
	for key := range opts.Args {
		uncheckedArgs[key] = true
	}

	for _, validArg := range desc.ValidArgs {
		delete(uncheckedArgs, validArg.Key)
		if arg, ok := opts.Args[validArg.Key]; ok {
			if validArg.IsMultipleValues {
				if !validArg.AllowEmpty && len(arg.Values) == 0 {
					return &utils.ArgEmptyValueError{ArgKey: validArg.Key}
				}
			} else {
				if len(arg.Values) != 1 {
					return &utils.ArgNotSingleValueError{ArgKey: validArg.Key}
				}
			}
			if !validArg.ValidValues.IsEmpty() {
				for _, value := range arg.Values {
					if !validArg.ValidValues.Contains(value) {
						return &utils.UnsupportedError{Type: OptionTypeArgValue, Idents: []string{value.Str()}}
					}
				}

			}
		} else if validArg.DefaultValue != nil {
			// set default value
			opts.Args[validArg.Key] = Arg{Key: validArg.Key, Values: []Value{*validArg.DefaultValue}}
		}
	}

	if len(uncheckedArgs) > 0 && !desc.AllowUnexpectedlyArg {
		args := make([]string, 0, len(uncheckedArgs))
		for arg := range uncheckedArgs {
			args = append(args, arg)
		}
		return &utils.UnexpectedError{Type: string(OptionTypeArgKey), Idents: args}
	}
	return nil
}

type TypeInfo struct {
	Name     string
	Assigned string
	Ast      ast.TypeSpec
}
