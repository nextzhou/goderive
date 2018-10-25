package utils

import (
	"go/ast"
	"regexp"
	"strings"
)

var identPattern = regexp.MustCompile(`^[_a-zA-Z]\w*$`)

func ValidateIdentName(name string) bool {
	return identPattern.MatchString(name)
}

func IsSortableType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "string", "uintptr", "rune", "byte":
		return true
	default:
		return false
	}
}

func IsComparableSimpleType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "complex64", "complex128",
		"string", "uintptr", "rune", "byte", "bool":
		return true
	default:
		return false
	}
}

// check if operator == is defined on this type
func IsComparableType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return IsComparableSimpleType(e.Name)
	case *ast.BasicLit:
		return IsComparableSimpleType(e.Value)
	case *ast.StarExpr, *ast.ChanType:
		// pointer and channel types are aways comparable
		return true
	case *ast.ArrayType:
		// array type is comparable, but not slice
		return e.Len != nil
	default:
		return false
	}
}

func IsBaseType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "complex64", "complex128",
		"string", "uintptr", "rune", "byte":
		return true
	default:
		return false
	}
}

func IsExported(ident string) bool {
	if len(ident) == 0 {
		return false
	}
	// selector expr. e.g. time.Time should be exported
	if idx := strings.IndexByte(ident, '.'); idx > 0 {
		return IsExported(ident[idx+1:])
	}
	initial := ident[0]
	return 'A' <= initial && initial <= 'Z'
}

func IsSelectorExpr(expr string) bool {
	return strings.IndexByte(expr, '.') > 0
}

func IsLocalPath(path string) bool {
	return !strings.ContainsRune(strings.SplitN(path, "/", 2)[0], '.')
}
