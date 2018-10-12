package utils

import "fmt"

type InvalidIdentError struct {
	Type  string
	Ident string
}

func (e *InvalidIdentError) Error() string {
	return fmt.Sprintf("invalid %s %#v", e.Type, e.Ident)
}

type ConflictingOptionError struct {
	Type  string
	Ident string
}

func (e *ConflictingOptionError) Error() string {
	return fmt.Sprintf("already existed %s %#v", e.Type, e.Ident)
}

type UnexpectedError struct {
	Type   string
	Idents []string
}

func (e *UnexpectedError) Error() string {
	if len(e.Idents) == 1 {
		return fmt.Sprintf("unexpected %s %#v", e.Type, e.Idents[0])
	}
	return fmt.Sprintf("unexpected %s %#v", e.Type, e.Idents)
}

type OnlySupportError struct {
	Supported string
	Got       string
}

func (e *OnlySupportError) Error() string {
	return fmt.Sprintf("support only %s, not %s", e.Supported, e.Got)
}

type UnsupportedError struct {
	Type   string
	Idents []string
}

func (e *UnsupportedError) Error() string {
	if len(e.Idents) == 1 {
		return fmt.Sprintf("unsupported %s %#v", e.Type, e.Idents[0])
	}
	return fmt.Sprintf("unsupported %s %#v", e.Type, e.Idents)
}

type ArgNotSingleValueError struct {
	ArgKey string
}

func (e *ArgNotSingleValueError) Error() string {
	return fmt.Sprintf("option %#v should accept a single value", e.ArgKey)
}

type ArgEmptyValueError struct {
	ArgKey string
}

func (e *ArgEmptyValueError) Error() string {
	return fmt.Sprintf("value of option %#v shouldn't be empty", e.ArgKey)
}

type NotExistedError struct {
	Type  string
	Ident string
}

func (e *NotExistedError) Error() string {
	return fmt.Sprintf("not existed %s %#v", e.Type, e.Ident)
}
