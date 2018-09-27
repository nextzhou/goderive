package set

import (
	"io"
	"text/template"
)

var setTemplate = `
type {{ .SetName }} map[{{ .TypeName }}]struct{}


func New{{ .CapitalizeSetName }}(capacity int) *{{ .SetName }} {
	var set {{ .SetName }}
	if capacity > 0 {
		set = make(map[{{ .TypeName }}]struct{}, capacity)
	} else {
		set = make(map[{{ .TypeName }}]struct{})
	}
	return (*{{ .SetName }})(&set)
}

func New{{ .CapitalizeSetName }}FromSlice(items []{{ .TypeName }}) *{{ .SetName }} {
	set := make(map[{{ .TypeName }}]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return (*{{ .SetName }})(&set)
}

func (set *{{ .SetName }}) Extend(items ...{{ .TypeName }}) {
	for _, item := range items {
		(*set)[item] = struct{}{}
	}
}

func (set *{{ .SetName }}) Len() int {
	if set == nil {
		return 0
	}
	return len(*set)
}

func (set *{{ .SetName }}) IsEmpty() bool {
	return set == nil || set.Len() == 0
}

func (set *{{ .SetName }}) ToSlice() []{{ .TypeName }} {
	if set == nil {
		return nil
	}
	s := make([]{{ .TypeName }}, 0, set.Len())
	set.ForEach(func(item {{.TypeName}}) {
		s = append(s, item)
	})
	return s
}

func (set *{{ .SetName }}) Put(key {{ .TypeName }}) {
	(*set)[key] = struct{}{}
}

func (set *{{ .SetName }}) Clear() {
	*set = make(map[{{ .TypeName }}]struct{})
}

func (set *{{ .SetName }}) Clone() *{{ .SetName }} {
	cloned := New{{ .SetName }}(set.Len())
	for item := range *set {
		(*cloned)[item] = struct{}{}
	}
	return cloned
}

func (set *{{ .SetName }}) Difference(another *{{ .SetName }}) *{{ .SetName }} {
	difference := New{{ .SetName }}(0)
	for item := range *set {
		if !another.Contains(item) {
			difference.Put(item)
		}
	}
	return difference
}

func (set *{{ .SetName }}) Equal(another *{{ .SetName }}) bool {
	if set.Len() != another.Len() {
		return false
	}
	for item := range *set {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *{{ .SetName }}) Intersect(another *{{ .SetName }}) *{{ .SetName }} {
	intersection := New{{ .SetName }}(0)
	if set.Len() < another.Len() {
		for item := range *set {
			if another.Contains(item) {
				intersection.Put(item)
			}
		}
	} else {
		for item := range *another {
			if set.Contains(item) {
				intersection.Put(item)
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
		set.Put(item)
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
	for item := range *set {
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
	for item := range *set {
		f(item)
	}
}

func (set *{{ .SetName }}) Filter(f func({{ .TypeName }}) bool) *{{ .SetName }} {
	result := New{{ .SetName }}(0)
	set.ForEach(func(item {{ .TypeName }}) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set {{ .SetName }}) Remove(key {{ .TypeName }}) {
	delete(set, key)
}

func (set {{ .SetName }}) Contains(key {{ .TypeName }}) bool {
	_, ok := set[key]
	return ok
}

func (set {{ .SetName }}) ContainsAny(keys ...{{ .TypeName }}) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set {{ .SetName }}) ContainsAll(keys ...{{ .TypeName }}) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *{{ .SetName }}) String() string {
	return fmt.Sprint(set.ToSlice())
}

func (set *{{ .SetName }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *{{ .SetName }}) UnmarshalJSON(b []byte) error {
	s := make([]{{ .TypeName }}, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *New{{ .SetName }}FromSlice(s)
	return nil
}
`

var tpl, _ = template.New("set").Parse(setTemplate)

type TemplateArgs struct {
	RawTypeName       string
	TypeName          string
	SetName           string
	CapitalizeSetName string
}

func (ta TemplateArgs) GenerateTo(w io.Writer) error {
	return tpl.Execute(w, ta)
}
