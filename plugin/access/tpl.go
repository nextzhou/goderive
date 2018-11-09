package access

import (
	"fmt"
	"go/ast"
	"io"
	"text/template"

	"github.com/nextzhou/goderive/utils"

	"github.com/nextzhou/goderive/plugin"
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

func (ta *TemplateArgs) AddField(field *ast.Field, opts *plugin.Options, t *utils.NameWithPkg) error {
	if opts != nil && opts.WithFlag("Ignore") {
		return nil
	}
	var rename *plugin.Value
	if opts != nil {
		rename = opts.GetValue("RenameGet")
	}
	if rename.IsNil() {
		for _, name := range field.Names {
			ta.Fields = append(ta.Fields, Field{Name: name.Name, GetFuncName: genGetName(name.Name), TypeName: t.Name})
		}
	} else {
		if len(field.Names) != 1 {
			return fmt.Errorf(`"RenameGet" field can only have one name`)
		}
		ta.Fields = append(ta.Fields, Field{Name: field.Names[0].Name, GetFuncName: rename.Str(), TypeName: t.Name})
	}
	return nil
}
