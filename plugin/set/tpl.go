package set

import (
	"io"
	"text/template"
)

var setTemplate = `
type {{ .SetName }} struct {
	{{ if eq .Order "Key" -}}
	cmp             func(i, j {{ .TypeName }}) bool
	elements        map[{{ .TypeName }}]uint32
	elementSequence []{{ .TypeName }}
	{{- else if eq .Order "Append" -}}
	elements        map[{{ .TypeName }}]uint32
	elementSequence []{{ .TypeName }}
	{{- else -}}
	elements map[{{ .TypeName }}]struct{}
	{{- end }}
}

{{ if eq .Order "Key" -}}
func {{ .New }}{{ .CapitalizeSetName }}(capacity int, cmp func(i, j {{ .TypeName }}) bool) *{{ .SetName }} {
{{- else -}}
func {{ .New }}{{ .CapitalizeSetName }}(capacity int) *{{ .SetName }} {
{{- end }}
	set := new({{ .SetName }})
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	if capacity > 0 {
		set.elements = make(map[{{ .TypeName }}]uint32, capacity)
		set.elementSequence = make([]{{ .TypeName }}, 0, capacity)
	} else {
		set.elements = make(map[{{ .TypeName }}]uint32)
	}
	{{- else -}}
	if capacity > 0 {
		set.elements = make(map[{{ .TypeName }}]struct{}, capacity)
	} else {
		set.elements = make(map[{{ .TypeName }}]struct{})
	}
	{{- end }}
	{{- if eq .Order "Key" }}
	set.cmp = cmp
	{{- end }}
	return set
}

{{ if eq .Order "Key" -}}
func {{ .New }}{{ .CapitalizeSetName }}FromSlice(items []{{ .TypeName }}, cmp func(i, j {{ .TypeName }}) bool) *{{ .SetName }} {
	set := {{ .New }}{{ .CapitalizeSetName }}(len(items), cmp)
{{- else -}}
func {{ .New }}{{ .CapitalizeSetName }}FromSlice(items []{{ .TypeName }}) *{{ .SetName }} {
	set := {{ .New }}{{ .CapitalizeSetName }}(len(items))
{{- end }}
	for _, item := range items {
		set.Append(item)
	}
	return set
}
{{- if and (eq .Order "Key") .IsComparable }}

func {{ .New }}Ascending{{ .CapitalizeSetName }}(capacity int) *{{ .SetName }} {
	return {{ .New }}{{ .CapitalizeSetName }}(capacity, func(i, j {{ .TypeName }}) bool { return i < j })
}

func {{ .New }}Descending{{ .CapitalizeSetName }}(capacity int) *{{ .SetName }} {
	return {{ .New }}{{ .CapitalizeSetName }}(capacity, func(i, j {{ .TypeName }}) bool { return i > j })
}

func {{ .New }}Ascending{{ .CapitalizeSetName }}FromSlice(items []{{ .TypeName }}) *{{ .SetName }} {
	return {{ .New }}{{ .CapitalizeSetName }}FromSlice(items, func(i, j {{ .TypeName }}) bool { return i < j })
}

func {{ .New }}Descending{{ .CapitalizeSetName }}FromSlice(items []{{ .TypeName }}) *{{ .SetName }} {
	return {{ .New }}{{ .CapitalizeSetName }}FromSlice(items, func(i, j {{ .TypeName }}) bool { return i > j })
}
{{- end }}

func (set *{{ .SetName }}) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *{{ .SetName }}) IsEmpty() bool {
	return set.Len() == 0
}

func (set *{{ .SetName }}) ToSlice() []{{ .TypeName }} {
	if set == nil {
		return nil
	}
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	s := make([]{{ .TypeName }}, set.Len())
	copy(s, set.elementSequence)
	{{- else -}}
	s := make([]{{ .TypeName }}, 0, set.Len())
	set.ForEach(func(item {{.TypeName}}) {
		s = append(s, item)
	})
	{{- end }}
	return s
}
{{- if or (eq .Order "Append") (eq .Order "Key") }}

// NOTICE: efficient but unsafe
func (set *{{ .SetName }}) ToSliceRef() []{{ .TypeName }} {
	return set.elementSequence
}
{{- end }}

func (set *{{ .SetName }}) Append(keys ...{{ .TypeName }}) {
	for _, key := range keys {
		{{ if eq .Order "Append" -}}
		if _, ok := set.elements[key]; !ok {
			set.elements[key] = uint32(len(set.elementSequence))
			set.elementSequence = append(set.elementSequence, key)
		}
		{{- else if eq .Order "Key" -}}
		if _, ok := set.elements[key]; !ok {
			idx := sort.Search(len(set.elementSequence), func(i int) bool {
				return set.cmp(key, set.elementSequence[i])
			})
			l := len(set.elementSequence)
			set.elementSequence = append(set.elementSequence, key)
			for i := l; i > idx; i-- {
				set.elements[set.elementSequence[i]] = uint32(i + 1)
				set.elementSequence[i] = set.elementSequence[i-1]
			}
			set.elements[set.elementSequence[idx]] = uint32(idx + 1)
			set.elementSequence[idx] = key
			set.elements[key] = uint32(idx)
		}
		{{- else -}}
		set.elements[key] = struct{}{}
		{{- end }}
	}
}

func (set *{{ .SetName }}) Clear() {
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	set.elements = make(map[{{ .TypeName }}]uint32)
	set.elementSequence = set.elementSequence[:0]
	{{- else -}}
	set.elements = make(map[{{ .TypeName }}]struct{})
	{{- end }}
}

func (set *{{ .SetName }}) Clone() *{{ .SetName }} {
	{{ if eq .Order "Key" -}}
	cloned := {{ .New }}{{ .CapitalizeSetName }}(set.Len(), set.cmp)
	{{- else -}}
	cloned := {{ .New }}{{ .CapitalizeSetName }}(set.Len())
	{{- end }}
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
	}
	{{- else -}}
	for item := range set.elements {
		cloned.elements[item] = struct{}{}
	}
	{{- end }}
	return cloned
}

func (set *{{ .SetName }}) Difference(another *{{ .SetName }}) *{{ .SetName }} {
	{{ if eq .Order "Key" -}}
	difference := {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp)
	{{- else -}}
	difference := {{ .New }}{{ .CapitalizeSetName }}(0)
	{{- end }}
	set.ForEach(func(item {{ .TypeName }}) {
		if !another.Contains(item) {
			difference.Append(item)
		}
	})
	return difference
}

func (set *{{ .SetName }}) Equal(another *{{ .SetName }}) bool {
	if set.Len() != another.Len() {
		return false
	}
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	return set.ContainsAll(another.elementSequence...)
	{{- else -}}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
	{{- end }}
}

{{ if eq .Order "Append" -}}
// TODO keep order
{{ end -}}
func (set *{{ .SetName }}) Intersect(another *{{ .SetName }}) *{{ .SetName }} {
	{{ if eq .Order "Key" -}}
	intersection := {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp)
	{{- else -}}
	intersection := {{ .New }}{{ .CapitalizeSetName }}(0)
	{{- end }}
	if set.Len() < another.Len() {
		for item := range set.elements {
			if another.Contains(item) {
				intersection.Append(item)
			}
		}
	} else {
		for item := range another.elements {
			if set.Contains(item) {
				intersection.Append(item)
			}
		}
	}
	return intersection
}

func (set *{{ .SetName }}) Union(another *{{ .SetName }}) *{{ .SetName }} {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *{{ .SetName }}) InPlaceUnion(another *{{ .SetName }}) {
	another.ForEach(func(item {{ .TypeName }}) {
		set.Append(item)
	})
}

func (set *{{ .SetName }}) IsProperSubsetOf(another *{{ .SetName }}) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *{{ .SetName }}) IsProperSupersetOf(another *{{ .SetName }}) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *{{ .SetName }}) IsSubsetOf(another *{{ .SetName }}) bool {
	if set.Len() > another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *{{ .SetName }}) IsSupersetOf(another *{{ .SetName }}) bool {
	return another.IsSubsetOf(set)
}

func (set *{{ .SetName }}) ForEach(f func({{ .TypeName }})) {
	if set.IsEmpty() {
		return
	}
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	for _, item := range set.elementSequence {
		f(item)
	}
	{{- else -}}
	for item := range set.elements {
		f(item)
	}
	{{- end }}
}
{{- if or (eq .Order "Append") (eq .Order "Key") }}

func (set *{{ .SetName }}) ForEachWithIndex(f func(int, {{ .TypeName }})) {
	if set.IsEmpty() {
		return
	}
	for idx, item := range set.elementSequence {
		f(idx, item)
	}
}
{{- end }}

func (set *{{ .SetName }}) Filter(f func({{ .TypeName }}) bool) *{{ .SetName }} {
	{{ if eq .Order "Key" -}}
	result := {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp)
	{{- else -}}
	result := {{ .New }}{{ .CapitalizeSetName }}(0)
	{{- end }}
	set.ForEach(func(item {{ .TypeName }}) {
		if f(item) {
			result.Append(item)
		}
	})
	return result
}

func (set *{{ .SetName }}) Remove(key {{ .TypeName }}) {
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	if idx, ok := set.elements[key]; ok {
		l := set.Len()
		delete(set.elements, key)
		for ; idx < uint32(l-1); idx++ {
			item := set.elementSequence[idx+1]
			set.elementSequence[idx] = item
			set.elements[item] = idx
		}
		set.elementSequence = set.elementSequence[:l-1]
	}
	{{- else -}}
	delete(set.elements, key)
	{{- end }}
}

func (set *{{ .SetName }}) Contains(key {{ .TypeName }}) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *{{ .SetName }}) ContainsAny(keys ...{{ .TypeName }}) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *{{ .SetName }}) ContainsAll(keys ...{{ .TypeName }}) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}
{{- if or (eq .Order "Append") (eq .Order "Key") }}

func (set *{{ .SetName }}) DoUntil(f func({{ .TypeName }}) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *{{ .SetName }}) DoWhile(f func({{ .TypeName }}) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}
{{- end }}

func (set *{{ .SetName }}) DoUntilError(f func({{ .TypeName }}) error) error {
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	for _, item := range set.elementSequence {
	{{- else -}}
	for item := range set.elements {
	{{- end }}
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

func (set *{{ .SetName }}) All(f func({{ .TypeName }}) bool) bool {
	for item := range set.elements {
		if !f(item) {
			return false
		}
	}
	return true
}

func (set *{{ .SetName }}) Any(f func({{ .TypeName }}) bool) bool {
	for item := range set.elements {
		if f(item) {
			return true
		}
	}
	return false
}

func (set *{{ .SetName }}) FindBy(f func({{ .TypeName }}) bool) *{{ .TypeName }} {
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	for _, item := range set.elementSequence {
	{{- else -}}
	for item := range set.elements {
	{{- end }}
		if f(item) {
			return &item
		}
	}
	return nil
}
{{- if or (eq .Order "Append") (eq .Order "Key") }}

func (set *{{ .SetName }}) FindLastBy(f func({{ .TypeName }}) bool) *{{ .TypeName }} {
	for i := set.Len() - 1; i >= 0; i-- {
		if item := set.elementSequence[i]; f(item) {
			return &item
		}
	}
	return nil
}
{{- end }}

func (set *{{ .SetName }}) CountBy(f func({{ .TypeName }}) bool) int {
	count := 0
	set.ForEach(func(item {{ .TypeName }}) {
		if f(item) {
			count++
		}
	})
	return count
}

func (set *{{ .SetName }}) GroupByBool(f func({{ .TypeName }}) bool) (trueGroup *{{ .SetName }}, falseGroup *{{ .SetName }}) {
	{{ if eq .Order "Key" -}}
	trueGroup, falseGroup = {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp), {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp)
	{{- else -}}
	trueGroup, falseGroup = {{ .New }}{{ .CapitalizeSetName }}(0), {{ .New }}{{ .CapitalizeSetName }}(0)
	{{- end }}
	set.ForEach(func(item {{ .TypeName }}) {
		if f(item) {
			trueGroup.Append(item)
		} else {
			falseGroup.Append(item)
		}
	})
	return trueGroup, falseGroup
}

func (set *{{ .SetName }}) GroupByStr(f func({{ .TypeName }}) string) map[string]*{{ .SetName }} {
	groups := make(map[string]*{{ .SetName }})
	set.ForEach(func(item {{ .TypeName }}) {
		key := f(item)
		group := groups[key]
		if group == nil {
			{{ if eq .Order "Key" -}}
			group = {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp)
			{{- else -}}
			group = {{ .New }}{{ .CapitalizeSetName }}(0)
			{{- end }}
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (set *{{ .SetName }}) GroupByInt(f func({{ .TypeName }}) int) map[int]*{{ .SetName }} {
	groups := make(map[int]*{{ .SetName }})
	set.ForEach(func(item {{ .TypeName }}) {
		key := f(item)
		group := groups[key]
		if group == nil {
			{{ if eq .Order "Key" -}}
			group = {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp)
			{{- else -}}
			group = {{ .New }}{{ .CapitalizeSetName }}(0)
			{{- end }}
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (set *{{ .SetName }}) GroupBy(f func({{ .TypeName }}) interface{}) map[interface{}]*{{ .SetName }} {
	groups := make(map[interface{}]*{{ .SetName }})
	set.ForEach(func(item {{ .TypeName }}) {
		key := f(item)
		group := groups[key]
		if group == nil {
			{{ if eq .Order "Key" -}}
			group = {{ .New }}{{ .CapitalizeSetName }}(0, set.cmp)
			{{- else -}}
			group = {{ .New }}{{ .CapitalizeSetName }}(0)
			{{- end }}
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

// f: func({{ .TypeName }}) T
// return: []T
func (set *{{ .SetName }}) Map(f interface{}) interface{} {
	expected := "f should be func({{ .TypeName }})T"
	ft := reflect.TypeOf(f)
	fVal := reflect.ValueOf(f)
	if ft.Kind() != reflect.Func {
		panic(expected)
	}
	if ft.NumIn() != 1 {
		panic(expected)
	}
	elemType := reflect.TypeOf(new({{ .TypeName }})).Elem()
	if ft.In(0) != elemType {
		panic(expected)
	}
	if ft.NumOut() != 1 {
		panic(expected)
	}
	outType := ft.Out(0)
	result := reflect.MakeSlice(reflect.SliceOf(outType), 0, set.Len())
	set.ForEach(func(item {{ .TypeName }}) {
		result = reflect.Append(result, fVal.Call([]reflect.Value{reflect.ValueOf(item)})[0])
	})
	return result.Interface()
}

// f: func({{ .TypeName }}) *T
//    func({{ .TypeName }}) (T, bool)
//    func({{ .TypeName }}) (T, error)
// return: []T
func (set *{{ .SetName }}) FilterMap(f interface{}) interface{} {
	expected := "f should be func({{ .TypeName }}) *T / func({{ .TypeName }}) (T, bool) / func({{ .TypeName }}) (T, error)"
	ft := reflect.TypeOf(f)
	fVal := reflect.ValueOf(f)
	if ft.Kind() != reflect.Func {
		panic(expected)
	}
	if ft.NumIn() != 1 {
		panic(expected)
	}
	in := ft.In(0)
	if in != reflect.TypeOf(new({{ .TypeName }})).Elem() {
		panic(expected)
	}
	var outType reflect.Type
	var filter func([]reflect.Value) *reflect.Value
	if ft.NumOut() == 1 {
		// func({{ .TypeName }}) *T
		outType = ft.Out(0)
		if outType.Kind() != reflect.Ptr {
			panic(expected)
		}
		outType = outType.Elem()
		filter = func(values []reflect.Value) *reflect.Value {
			if values[0].IsNil() {
				return nil
			}
			val := values[0].Elem()
			return &val
		}
	} else if ft.NumOut() == 2 {
		outType = ft.Out(0)
		checker := ft.Out(1)
		if checker == reflect.TypeOf(true) {
			// func({{ .TypeName }}) (T, bool)
			filter = func(values []reflect.Value) *reflect.Value {
				if values[1].Interface().(bool) {
					return &values[0]
				}
				return nil
			}
		} else if checker.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			// func({{ .TypeName }}) (T, error)
			filter = func(values []reflect.Value) *reflect.Value {
				if values[1].IsNil() {
					return &values[0]
				}
				return nil
			}
		} else {
			panic(expected)
		}
	} else {
		panic(expected)
	}

	result := reflect.MakeSlice(reflect.SliceOf(outType), 0, set.Len())
	set.ForEach(func(item {{ .TypeName }}) {
		ret := fVal.Call([]reflect.Value{reflect.ValueOf(item)})
		if val := filter(ret); val != nil {
			result = reflect.Append(result, *val)
		}
	})
	return result.Interface()
}

func (set *{{ .SetName }}) Reduce(f func({{ .TypeName }}, {{ .TypeName }}) {{ .TypeName }}) {{ .TypeName }} {
	if set.IsEmpty() {
		var defaultVal {{ .TypeName }}
		return defaultVal
	}
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	ret := set.elementSequence[0]
	for _, item := range set.elementSequence[1:] {
		ret = f(ret, item)
	}
	{{- else -}}
	var ret {{ .TypeName }}
	first := true
	for item := range set.elements {
		if first {
			ret = item
			first = false
			continue
		}
		ret = f(ret, item)
	}
	{{- end }}
	return ret
}

func (set *{{ .SetName }}) Fold(init {{ .TypeName }}, f func({{ .TypeName }}, {{ .TypeName }}) {{ .TypeName }}) {{ .TypeName }} {
	if set.IsEmpty() {
		return init
	}
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	for _, item := range set.elementSequence {
	{{- else -}}
	for item := range set.elements {
	{{- end }}
		init = f(init, item)
	}
	return init
}

func (set *{{ .SetName }}) String() string {
	{{ if or (eq .Order "Append") (eq .Order "Key") -}}
	return fmt.Sprint(set.elementSequence)
	{{- else -}}
	return fmt.Sprint(set.ToSlice())
	{{- end }}
}

func (set {{ .SetName }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *{{ .SetName }}) UnmarshalJSON(b []byte) error {
	{{ if (eq .Order "Key") -}}
	return fmt.Errorf("unsupported")
	{{- else -}}
	s := make([]{{ .TypeName }}, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *{{ .New }}{{ .CapitalizeSetName }}FromSlice(s)
	return nil
	{{- end }}
}
`

var tpl, _ = template.New("set").Parse(setTemplate)

type TemplateArgs struct {
	TypeName          string
	SetName           string
	CapitalizeSetName string
	Order             string
	IsComparable      bool
	New               string
}

func (ta TemplateArgs) GenerateTo(w io.Writer) error {
	return tpl.Execute(w, ta)
}
