package main

import (
	"io"
	"text/template"
)

var setTemplate = `
type {{ .SetName }} map[{{ .TypeName }}]struct{}


func New{{ .CapitalizeSetName }} (capacity int) {{ .SetName }} {
	if capacity > 0 {
		return make(map[{{ .TypeName }}]struct{}, capacity)
	}
	return make(map[{{ .TypeName }}]struct{})
}

func (set {{ .SetName }}) Put(key {{ .TypeName }}) {
	set[key] = struct{}{}
}

func (set {{ .SetName }}) Delete(key {{ .TypeName }}) {
	delete(set, key)
}

func (set {{ .SetName }}) Contains(key {{ .TypeName }}) bool {
	_, ok := set[key]
	return ok
}
`

var tmpl, _ = template.New("set").Parse(setTemplate)

type TemplateArgs struct {
	TypeName          string
	SetName           string
	CapitalizeSetName string
}

func (ta TemplateArgs) GenerateTo(w io.Writer) error {
	return tmpl.Execute(w, ta)
}
