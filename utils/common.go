//go:generate goderive
package utils

import (
	"go/ast"
	"io"
	"reflect"
	"strings"
	"unicode"

	"github.com/olekukonko/tablewriter"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	if 'a' <= s[0] && s[0] <= 'z' {
		b := []byte(s)
		b[0] = b[0] - (byte('a') - 'A')
		return string(b)
	}
	return s
}

type TableWriter struct {
	table *tablewriter.Table
}

func NewTableWriter(w io.Writer) *TableWriter {
	t := tablewriter.NewWriter(w)
	t.SetColWidth(60)
	t.SetBorders(tablewriter.Border{Left: false, Right: false, Top: false, Bottom: false})
	t.SetHeaderLine(false)
	t.SetColumnSeparator("\t")
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	return &TableWriter{table: t}
}

func (tw *TableWriter) Append(row []string) {
	tw.table.Append(row)
}

func (tw *TableWriter) Render() {
	tw.table.Render()
}

// derive-set: Order=Append
type Str = string

// derive-set: Rename=StrOrderSet;Order=Key
type Str2 = string

const HeaderComment = "// Code generated by https://github.com/nextzhou/goderive. DO NOT EDIT.\n\n"

func ToExported(ident string) string {
	if len(ident) == 0 {
		return ident
	}
	return string(unicode.ToUpper(rune(ident[0]))) + ident[1:]
}

func ToUnexported(ident string) string {
	if len(ident) == 0 {
		return ident
	}
	return string(unicode.ToLower(rune(ident[0]))) + ident[1:]
}

// time.Time => (time, Time)
func SplitSelectorExpr(expr string) (string, string) {
	dotIdx := strings.IndexByte(expr, '.')
	if dotIdx > 0 {
		return expr[:dotIdx], expr[dotIdx+1:]
	}
	return "", expr
}

func SelectorExprString(se *ast.SelectorExpr) string {
	return se.X.(*ast.Ident).Name + "." + se.Sel.Name
}

// a/b/c.v5 => c
func PkgNameFromPath(path string) string {
	path = strings.Trim(path, `"`)
	terms := strings.Split(path, "/")
	name := terms[len(terms)-1]
	dotIdx := strings.IndexByte(name, '.')
	if dotIdx == -1 {
		return name
	}
	return name[:dotIdx]
}

type NameWithPkg struct {
	Name string
	Pkgs *StrSet
}

func NewNameWithPkg(name string) *NameWithPkg {
	return &NameWithPkg{
		Name: name,
		Pkgs: NewStrSet(0),
	}
}

func TypeNameWithPkg(expr ast.Expr) *NameWithPkg {
	ret := NewNameWithPkg("")
	switch e := expr.(type) {
	case *ast.Ident:
		ret.Name = e.Name
	case *ast.StarExpr:
		s := TypeNameWithPkg(e.X)
		if s == nil {
			return nil
		}
		ret.Name = "*" + s.Name
		ret.Pkgs.InPlaceUnion(s.Pkgs)
	case *ast.SelectorExpr:
		s := TypeNameWithPkg(e.X)
		if s == nil {
			return nil
		}
		ret.Pkgs.Append(s.Name)
		ret.Pkgs.InPlaceUnion(s.Pkgs)
		ret.Name = s.Name + "." + e.Sel.Name
	case *ast.SliceExpr:
		s := TypeNameWithPkg(e.X)
		if s == nil {
			return nil
		}
		ret.Pkgs.InPlaceUnion(s.Pkgs)
		ret.Name = "[]" + s.Name
	case *ast.MapType:
		k, v := TypeNameWithPkg(e.Key), TypeNameWithPkg(e.Value)
		if k == nil || v == nil {
			return nil
		}
		ret.Pkgs.InPlaceUnion(k.Pkgs)
		ret.Pkgs.InPlaceUnion(v.Pkgs)
		ret.Name = "map[" + k.Name + "]" + v.Name
	case *ast.ArrayType:
		elem := TypeNameWithPkg(e.Elt)
		if elem == nil {
			return nil
		}
		ret.Pkgs.InPlaceUnion(elem.Pkgs)
		if e.Len == nil {
			ret.Name = "[]" + elem.Name
			break
		}
		l := TypeNameWithPkg(e.Len)
		if l == nil {
			return nil
		}
		ret.Pkgs.InPlaceUnion(l.Pkgs)
		ret.Name = "[" + l.Name + "]" + elem.Name
	case *ast.BasicLit:
		ret.Name = e.Value
	case *ast.ChanType:
		elem := TypeNameWithPkg(e.Value)
		if elem == nil {
			return nil
		}
		ret.Pkgs.InPlaceUnion(elem.Pkgs)
		if !e.Arrow.IsValid() {
			ret.Name = "chan " + elem.Name
			break
		}
		if e.Arrow == e.Begin {
			ret.Name = "<-chan " + elem.Name
			break
		}
		ret.Name = "chan<- " + elem.Name
	case *ast.FuncType:
		return FuncTypeNameWithPkg(*e)
	default:
		return nil
	}
	return ret
}

func FuncTypeNameWithPkg(e ast.FuncType) *NameWithPkg {
	ret := NewNameWithPkg("")
	ret.Name = "func("
	for idx, param := range e.Params.List {
		if idx != 0 {
			ret.Name += ", "
		}
		for idx, name := range param.Names {
			if idx != 0 {
				ret.Name += ", "
			}
			ret.Name += name.Name
		}
		p := TypeNameWithPkg(param.Type)
		if p == nil {
			return nil
		}
		ret.Pkgs.InPlaceUnion(p.Pkgs)
		if idx != 0 || len(param.Names) > 0 {
			ret.Name += " "
		}
		ret.Name += p.Name
	}
	ret.Name += ")"
	switch len(e.Results.List) {
	case 0:
		break
	case 1:
		ret.Name += " "
		item := e.Results.List[0]
		r := TypeNameWithPkg(item.Type)
		if r == nil {
			return nil
		}
		ret.Pkgs.InPlaceUnion(r.Pkgs)
		if len(item.Names) > 0 {
			ret.Name += "("
			for idx, name := range item.Names {
				if idx != 0 {
					ret.Name += ", "
				}
				ret.Name += name.Name
			}
			ret.Name += " " + r.Name + ")"
		} else {
			ret.Name += r.Name
		}
	default:
		ret.Name += " ("
		for idx, result := range e.Results.List {
			if idx != 0 {
				ret.Name += ", "
			}
			for idx, name := range result.Names {
				if idx != 0 {
					ret.Name += ", "
				}
				ret.Name += name.Name
			}
			r := TypeNameWithPkg(result.Type)
			if r == nil {
				return nil
			}
			ret.Pkgs.InPlaceUnion(r.Pkgs)
			if idx != 0 || len(result.Names) > 0 {
				ret.Name += " "
			}
			ret.Name += r.Name
		}
		ret.Name += ")"
	}
	return ret
}

func ExprTypeStr(expr ast.Expr) string {
	s := reflect.TypeOf(expr).String()
	s = strings.TrimPrefix(s, "*ast.")
	s = strings.TrimSuffix(s, "Expr")
	s = strings.TrimSuffix(s, "Type")
	return s
}
