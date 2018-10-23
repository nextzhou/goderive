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

func (s *{{ .SliceName }}) IsEmpty() bool {
	return s.Len() == 0
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

func (s *{{ .SliceName }}) Clear() {
	s.elements = s.elements[:0]
}

func (s *{{ .SliceName }}) Equal(another *{{ .SliceName }}) bool {
	if s.Len() != another.Len() {
		return false
	}
	for idx, item := range s.elements {
		if item != another.elements[idx] {
			return false
		}
	}
	return false
}

func (s *{{ .SliceName }}) Insert(idx int, items ...{{ .TypeName }}) {
	if idx < 0 {
		idx += s.Len()
	}
	if l := len(s.elements) + len(items); l > cap(s.elements) {
		// reallocate
		result := make([]{{ .TypeName }}, l)
		copy(result, s.elements[:idx])
		copy(result[idx:], items)
		copy(result[idx+len(items):], s.elements[idx:])
		s.elements = result
		return
	}

	l := s.Len()
	s.elements = append(s.elements, items...)
	copy(s.elements[idx+len(items):], s.elements[idx:l])
	copy(s.elements[idx:], items)
}

func (s *{{ .SliceName }}) Remove(idx int) {
	if idx < 0 {
		idx += s.Len()
	}
	s.elements = append(s.elements[:idx], s.elements[idx+1:]...)
}

func (s *{{ .SliceName }}) RemoveRange(from, to int) {
	if from < 0 {
		from += s.Len()
	}
	if to < 0 {
		to += s.Len()
	}
	s.elements = append(s.elements[:from], s.elements[to+1:]...)
}

func (s *{{ .SliceName }}) RemoveFrom(idx int) {
	if idx < 0 {
		idx += s.Len()
	}
	s.elements = s.elements[:idx]
}

func (s *{{ .SliceName }}) RemoveTo(idx int) {
	if idx < 0 {
		idx += s.Len()
	}
	s.elements = s.elements[idx + 1:]
}

func (s *{{ .SliceName }}) Concat(another *{{ .SliceName }}) *{{ .SliceName }} {
	result := s.Clone()
	if another.IsEmpty() {
		return result
	}
	result.Append(another.elements...)
	return result
}

func (s *{{ .SliceName }}) InPlaceConcat(another *{{ .SliceName }}) {
	if another.IsEmpty() {
		return
	}
	s.Append(another.elements...)
}

func (s *{{ .SliceName }}) ForEach(f func({{ .TypeName }})) {
	if s.IsEmpty() {
		return
	}
	for _, item := range s.elements {
		f(item)
	}
}

func (s *{{ .SliceName }}) Filter(f func({{ .TypeName }}) bool) *{{ .SliceName }} {
	result := {{ .New }}{{ .CapitalizeSliceName }}(0)
	for _, item := range s.elements {
		if f(item) {
			result.Append(item)
		}
	}
	return result
}

func (s *{{ .SliceName }}) Index(idx int) *{{ .TypeName }} {
	if idx < 0 {
		idx += s.Len()
	}
	return &s.elements[idx]
}

func (s *{{ .SliceName }}) IndexRange(from, to int) []{{ .TypeName }} {
	if from < 0 {
		from += s.Len()
	}
	if to < 0 {
		to += s.Len()
	}
	return s.elements[from:to]
}

func (s *{{ .SliceName }}) IndexFrom(idx int) []{{ .TypeName }} {
	if idx < 0 {
		idx += s.Len()
	}
	return s.elements[idx:]
}

func (s *{{ .SliceName }}) IndexTo(idx int) []{{ .TypeName }} {
	if idx < 0 {
		idx += s.Len()
	}
	return s.elements[:idx]
}

func (s *{{ .SliceName }}) Find(item {{ .TypeName }}) int {
	if s.IsEmpty() {
		return -1
	}
	for idx, n := range s.elements {
		if n == item {
			return idx
		}
	}
	return -1
}

func (s *{{ .SliceName }}) FindLast(item {{ .TypeName }}) int {
	for idx := s.Len() - 1; idx >= 0; idx-- {
		if s.elements[idx] == item {
			return idx
		}
	}
	return -1
}

func (s *{{ .SliceName }}) FindBy(f func({{ .TypeName }}) bool) int {
	if s.IsEmpty() {
		return -1
	}
	for idx, n := range s.elements {
		if f(n) {
			return idx
		}
	}
	return -1
}

func (s *{{ .SliceName }}) FindLastBy(f func({{ .TypeName }}) bool) int {
	for idx := s.Len() - 1; idx >= 0; idx-- {
		if f(s.elements[idx]) {
			return idx
		}
	}
	return -1
}

func (s *{{ .SliceName }}) Count(item {{ .TypeName }}) uint {
	count := uint(0)
	s.ForEach(func(n {{ .TypeName }}) {
		if n == item {
			count++
		}
	})
	return count
}

func (s *{{ .SliceName }}) CountBy(f func({{ .TypeName }}) bool) uint {
	count := uint(0)
	s.ForEach(func(item {{ .TypeName }}) {
		if f(item) {
			count++
		}
	})
	return count
}

func (s *{{ .SliceName }}) GroupByBool(f func({{ .TypeName }}) bool) (trueGroup, falseGroup *{{ .SliceName }}) {
	trueGroup, falseGroup = {{ .New }}{{ .CapitalizeSliceName }}(0), {{ .New }}{{ .CapitalizeSliceName }}(0)
	s.ForEach(func(item {{ .TypeName }}) {
		if f(item) {
			trueGroup.Append(item)
		} else {
			falseGroup.Append(item)
		}
	})
	return trueGroup, falseGroup
}

func (s {{ .SliceName }}) GroupByStr(f func({{ .TypeName }}) string) map[string]*{{ .SliceName }} {
	groups := make(map[string]*{{ .SliceName }})
	s.ForEach(func(item {{ .TypeName }}) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = {{ .New }}{{ .CapitalizeSliceName }}(0)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (s {{ .SliceName }}) GroupByInt(f func({{ .TypeName }}) int) map[int]*{{ .SliceName }} {
	groups := make(map[int]*{{ .SliceName }})
	s.ForEach(func(item {{ .TypeName }}) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = {{ .New }}{{ .CapitalizeSliceName }}(0)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (s *{{ .SliceName }}) GroupBy(f func({{ .TypeName }}) interface{}) map[interface{}]*{{ .SliceName }} {
	groups := make(map[interface{}]*{{ .SliceName }})
	s.ForEach(func(item {{ .TypeName }}) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = {{ .New }}{{ .CapitalizeSliceName }}(0)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (s *{{ .SliceName }}) String() string {
	return fmt.Sprint(s.elements)
}

func (s {{ .SliceName }}) MarshalJSON() ([]byte, error) {
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
