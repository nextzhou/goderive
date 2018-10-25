package access

import (
	"io"
	"text/template"
)

var accessTemplate = `
{{ range $idx, $field := .Fields }}
{{- if not (eq $idx 0) }}

{{ end -}}
func ({{ $.Receiver }} *{{ $.TypeName }}) {{ $field.GetFuncName }}() {{ $field.TypeName }} {
	if {{ $.Receiver }} == nil {
		var defaultVal {{ $field.TypeName }}
		return defaultVal
	}
	return {{ $.Receiver }}.{{ $field.Name }}
}
{{- end }}
`

var tpl, err = template.New("access").Parse(accessTemplate)

type TemplateArgs struct {
	TypeName string
	Receiver string
	Fields   []Field
}

type Field struct {
	Name        string
	GetFuncName string
	TypeName    string
}

func (ta TemplateArgs) GenerateTo(w io.Writer) error {
	if err != nil {
		return err
	}
	return tpl.Execute(w, ta)
}
