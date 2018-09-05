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

func (set {{ .SetName }}) ContainsAny(keys []{{ .TypeName }}) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set {{ .SetName }}) ContainsAll(keys []{{ .TypeName }}) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
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
