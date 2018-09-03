package main

import "text/template"

var setTemplate = `
type {{ .Name }}Set map[{{ .Name }}]struct{}


func New{{ .Name }}Set(capacity int) {{ .Name }}Set {
	if capacity > 0 {
		return make(map[{{ .Name }}]struct{}, capacity)
	}
	return make(map[{{ .Name }}]struct{})
}

func (set {{ .Name }}Set) Put(key {{ .Name }}) {
	set[key] = struct{}{}
}

func (set {{ .Name }}Set) Delete(key {{ .Name }}) {
	delete(set, key)
}

func (set {{ .Name }}Set) Contains(key {{ .Name }}) bool {
	_, ok := set[key]
	return ok
}
`

var tmpl, _ = template.New("set").Parse(setTemplate)
