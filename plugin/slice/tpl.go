package slice

import (
	"io"
	"text/template"
)

var sliceTemplate = `
type {{ .SliceName }} struct {
	elements []{{ .TypeName }}
}

func {{ .New }}{{ .CapitalizeSliceName }}(capacity int) *{{ .SliceName }} {
	return &{{ .SliceName }}{
		elements: make([]{{ .TypeName }}, 0, capacity),
	}
}

func {{ .New }}{{ .CapitalizeSliceName }}FromSlice(slice []{{ .TypeName }}) *{{ .SliceName }} {
	return &{{ .SliceName }}{
		elements: slice,
	}
}

func (s *{{ .SliceName }}) Len() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

func (s *{{ .SliceName }}) Append(items ...{{ .TypeName }}) {
	s.elements = append(s.elements, items...)
}

func (s *{{ .SliceName }}) Clone() *{{ .SliceName }} {
	cloned := &{{ .SliceName }}{
		elements: make([]{{ .TypeName }}, s.Len()),
	}
	copy(cloned.elements, s.elements)
	return cloned
}

func (s *{{ .SliceName }}) ToSlice() []{{ .TypeName }} {
	slice := make([]{{ .TypeName }}, s.Len())
	copy(slice, s.elements)
	return slice
}

func (s *{{ .SliceName }}) ToSliceRef() []{{ .TypeName }} {
	return s.elements
}

func (s *{{ .SliceName }}) String() string {
	return fmt.Sprint(s.elements)
}

func (s *{{ .SliceName }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.elements)
}

func (s *{{ .SliceName }}) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &s.elements)
}
`

var tpl, _ = template.New("slice").Parse(sliceTemplate)

type TemplateArgs struct {
	TypeName            string
	SliceName           string
	CapitalizeSliceName string
	IsComparable        bool
	New                 string
}

func (ta TemplateArgs) GenerateTo(w io.Writer) error {
	return tpl.Execute(w, ta)
}
