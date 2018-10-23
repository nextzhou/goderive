// Code generated by https://github.com/nextzhou/goderive. DO NOT EDIT.

package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type StrSet struct {
	elements        map[string]uint32
	elementSequence []string
}

func NewStrSet(capacity int) *StrSet {
	set := new(StrSet)
	if capacity > 0 {
		set.elements = make(map[string]uint32, capacity)
		set.elementSequence = make([]string, 0, capacity)
	} else {
		set.elements = make(map[string]uint32)
	}
	return set
}

func NewStrSetFromSlice(items []string) *StrSet {
	set := NewStrSet(len(items))
	for _, item := range items {
		set.Append(item)
	}
	return set
}

func (set *StrSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *StrSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *StrSet) ToSlice() []string {
	if set == nil {
		return nil
	}
	s := make([]string, set.Len())
	copy(s, set.elementSequence)
	return s
}

// NOTICE: efficient but unsafe
func (set *StrSet) ToSliceRef() []string {
	return set.elementSequence
}

func (set *StrSet) Append(keys ...string) {
	for _, key := range keys {
		if _, ok := set.elements[key]; !ok {
			set.elements[key] = uint32(len(set.elementSequence))
			set.elementSequence = append(set.elementSequence, key)
		}
	}
}

func (set *StrSet) Clear() {
	set.elements = make(map[string]uint32)
	set.elementSequence = set.elementSequence[:0]
}

func (set *StrSet) Clone() *StrSet {
	cloned := NewStrSet(set.Len())
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
	}
	return cloned
}

func (set *StrSet) Difference(another *StrSet) *StrSet {
	difference := NewStrSet(0)
	set.ForEach(func(item string) {
		if !another.Contains(item) {
			difference.Append(item)
		}
	})
	return difference
}

func (set *StrSet) Equal(another *StrSet) bool {
	if set.Len() != another.Len() {
		return false
	}
	return set.ContainsAll(another.elementSequence...)
}

// TODO keep order
func (set *StrSet) Intersect(another *StrSet) *StrSet {
	intersection := NewStrSet(0)
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

func (set *StrSet) Union(another *StrSet) *StrSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *StrSet) InPlaceUnion(another *StrSet) {
	another.ForEach(func(item string) {
		set.Append(item)
	})
}

func (set *StrSet) IsProperSubsetOf(another *StrSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *StrSet) IsProperSupersetOf(another *StrSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *StrSet) IsSubsetOf(another *StrSet) bool {
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

func (set *StrSet) IsSupersetOf(another *StrSet) bool {
	return another.IsSubsetOf(set)
}

func (set *StrSet) ForEach(f func(string)) {
	if set.IsEmpty() {
		return
	}
	for _, item := range set.elementSequence {
		f(item)
	}
}

func (set *StrSet) Filter(f func(string) bool) *StrSet {
	result := NewStrSet(0)
	set.ForEach(func(item string) {
		if f(item) {
			result.Append(item)
		}
	})
	return result
}

func (set *StrSet) Remove(key string) {
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
}

func (set *StrSet) Contains(key string) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *StrSet) ContainsAny(keys ...string) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *StrSet) ContainsAll(keys ...string) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *StrSet) DoUntil(f func(string) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrSet) DoWhile(f func(string) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrSet) DoUntilError(f func(string) error) error {
	for _, item := range set.elementSequence {
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

func (set *StrSet) All(f func(string) bool) bool {
	for item := range set.elements {
		if !f(item) {
			return false
		}
	}
	return true
}

func (set *StrSet) Any(f func(string) bool) bool {
	for item := range set.elements {
		if f(item) {
			return true
		}
	}
	return false
}

func (set *StrSet) FindBy(f func(string) bool) *string {
	for _, item := range set.elementSequence {
		if f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrSet) FindLastBy(f func(string) bool) *string {
	for i := set.Len() - 1; i >= 0; i-- {
		if item := set.elementSequence[i]; f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrSet) CountBy(f func(string) bool) int {
	count := 0
	set.ForEach(func(item string) {
		if f(item) {
			count++
		}
	})
	return count
}

func (set *StrSet) GroupByBool(f func(string) bool) (trueGroup *StrSet, falseGroup *StrSet) {
	trueGroup, falseGroup = NewStrSet(0), NewStrSet(0)
	set.ForEach(func(item string) {
		if f(item) {
			trueGroup.Append(item)
		} else {
			falseGroup.Append(item)
		}
	})
	return trueGroup, falseGroup
}

func (set *StrSet) GroupByStr(f func(string) string) map[string]*StrSet {
	groups := make(map[string]*StrSet)
	set.ForEach(func(item string) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = NewStrSet(0)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (set *StrSet) GroupByInt(f func(string) int) map[int]*StrSet {
	groups := make(map[int]*StrSet)
	set.ForEach(func(item string) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = NewStrSet(0)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (set *StrSet) GroupBy(f func(string) interface{}) map[interface{}]*StrSet {
	groups := make(map[interface{}]*StrSet)
	set.ForEach(func(item string) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = NewStrSet(0)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

// f: func(string) T
// return: []T
func (set *StrSet) Map(f interface{}) interface{} {
	expected := "f should be func(string)T"
	ft := reflect.TypeOf(f)
	fVal := reflect.ValueOf(f)
	if ft.Kind() != reflect.Func {
		panic(expected)
	}
	if ft.NumIn() != 1 {
		panic(expected)
	}
	elemType := reflect.TypeOf(new(string)).Elem()
	if ft.In(0) != elemType {
		panic(expected)
	}
	if ft.NumOut() != 1 {
		panic(expected)
	}
	outType := ft.Out(0)
	result := reflect.MakeSlice(reflect.SliceOf(outType), 0, set.Len())
	set.ForEach(func(item string) {
		result = reflect.Append(result, fVal.Call([]reflect.Value{reflect.ValueOf(item)})[0])
	})
	return result.Interface()
}

// f: func(string) *T
//    func(string) (T, bool)
//    func(string) (T, error)
// return: []T
func (set *StrSet) FilterMap(f interface{}) interface{} {
	expected := "f should be func(string) *T / func(string) (T, bool) / func(string) (T, error)"
	ft := reflect.TypeOf(f)
	fVal := reflect.ValueOf(f)
	if ft.Kind() != reflect.Func {
		panic(expected)
	}
	if ft.NumIn() != 1 {
		panic(expected)
	}
	in := ft.In(0)
	if in != reflect.TypeOf(new(string)).Elem() {
		panic(expected)
	}
	var outType reflect.Type
	var filter func([]reflect.Value) *reflect.Value
	if ft.NumOut() == 1 {
		// func(string) *T
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
			// func(string) (T, bool)
			filter = func(values []reflect.Value) *reflect.Value {
				if values[1].Interface().(bool) {
					return &values[0]
				}
				return nil
			}
		} else if checker.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			// func(string) (T, error)
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
	set.ForEach(func(item string) {
		ret := fVal.Call([]reflect.Value{reflect.ValueOf(item)})
		if val := filter(ret); val != nil {
			result = reflect.Append(result, *val)
		}
	})
	return result.Interface()
}

func (set *StrSet) String() string {
	return fmt.Sprint(set.elementSequence)
}

func (set StrSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *StrSet) UnmarshalJSON(b []byte) error {
	s := make([]string, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewStrSetFromSlice(s)
	return nil
}

type StrOrderSet struct {
	cmp             func(i, j string) bool
	elements        map[string]uint32
	elementSequence []string
}

func NewStrOrderSet(capacity int, cmp func(i, j string) bool) *StrOrderSet {
	set := new(StrOrderSet)
	if capacity > 0 {
		set.elements = make(map[string]uint32, capacity)
		set.elementSequence = make([]string, 0, capacity)
	} else {
		set.elements = make(map[string]uint32)
	}
	set.cmp = cmp
	return set
}

func NewStrOrderSetFromSlice(items []string, cmp func(i, j string) bool) *StrOrderSet {
	set := NewStrOrderSet(len(items), cmp)
	for _, item := range items {
		set.Append(item)
	}
	return set
}

func NewAscendingStrOrderSet(capacity int) *StrOrderSet {
	return NewStrOrderSet(capacity, func(i, j string) bool { return i < j })
}

func NewDescendingStrOrderSet(capacity int) *StrOrderSet {
	return NewStrOrderSet(capacity, func(i, j string) bool { return i > j })
}

func NewAscendingStrOrderSetFromSlice(items []string) *StrOrderSet {
	return NewStrOrderSetFromSlice(items, func(i, j string) bool { return i < j })
}

func NewDescendingStrOrderSetFromSlice(items []string) *StrOrderSet {
	return NewStrOrderSetFromSlice(items, func(i, j string) bool { return i > j })
}

func (set *StrOrderSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *StrOrderSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *StrOrderSet) ToSlice() []string {
	if set == nil {
		return nil
	}
	s := make([]string, set.Len())
	copy(s, set.elementSequence)
	return s
}

// NOTICE: efficient but unsafe
func (set *StrOrderSet) ToSliceRef() []string {
	return set.elementSequence
}

func (set *StrOrderSet) Append(keys ...string) {
	for _, key := range keys {
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
	}
}

func (set *StrOrderSet) Clear() {
	set.elements = make(map[string]uint32)
	set.elementSequence = set.elementSequence[:0]
}

func (set *StrOrderSet) Clone() *StrOrderSet {
	cloned := NewStrOrderSet(set.Len(), set.cmp)
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
	}
	return cloned
}

func (set *StrOrderSet) Difference(another *StrOrderSet) *StrOrderSet {
	difference := NewStrOrderSet(0, set.cmp)
	set.ForEach(func(item string) {
		if !another.Contains(item) {
			difference.Append(item)
		}
	})
	return difference
}

func (set *StrOrderSet) Equal(another *StrOrderSet) bool {
	if set.Len() != another.Len() {
		return false
	}
	return set.ContainsAll(another.elementSequence...)
}

func (set *StrOrderSet) Intersect(another *StrOrderSet) *StrOrderSet {
	intersection := NewStrOrderSet(0, set.cmp)
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

func (set *StrOrderSet) Union(another *StrOrderSet) *StrOrderSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *StrOrderSet) InPlaceUnion(another *StrOrderSet) {
	another.ForEach(func(item string) {
		set.Append(item)
	})
}

func (set *StrOrderSet) IsProperSubsetOf(another *StrOrderSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *StrOrderSet) IsProperSupersetOf(another *StrOrderSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *StrOrderSet) IsSubsetOf(another *StrOrderSet) bool {
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

func (set *StrOrderSet) IsSupersetOf(another *StrOrderSet) bool {
	return another.IsSubsetOf(set)
}

func (set *StrOrderSet) ForEach(f func(string)) {
	if set.IsEmpty() {
		return
	}
	for _, item := range set.elementSequence {
		f(item)
	}
}

func (set *StrOrderSet) Filter(f func(string) bool) *StrOrderSet {
	result := NewStrOrderSet(0, set.cmp)
	set.ForEach(func(item string) {
		if f(item) {
			result.Append(item)
		}
	})
	return result
}

func (set *StrOrderSet) Remove(key string) {
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
}

func (set *StrOrderSet) Contains(key string) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *StrOrderSet) ContainsAny(keys ...string) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *StrOrderSet) ContainsAll(keys ...string) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *StrOrderSet) DoUntil(f func(string) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrOrderSet) DoWhile(f func(string) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrOrderSet) DoUntilError(f func(string) error) error {
	for _, item := range set.elementSequence {
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

func (set *StrOrderSet) All(f func(string) bool) bool {
	for item := range set.elements {
		if !f(item) {
			return false
		}
	}
	return true
}

func (set *StrOrderSet) Any(f func(string) bool) bool {
	for item := range set.elements {
		if f(item) {
			return true
		}
	}
	return false
}

func (set *StrOrderSet) FindBy(f func(string) bool) *string {
	for _, item := range set.elementSequence {
		if f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrOrderSet) FindLastBy(f func(string) bool) *string {
	for i := set.Len() - 1; i >= 0; i-- {
		if item := set.elementSequence[i]; f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrOrderSet) CountBy(f func(string) bool) int {
	count := 0
	set.ForEach(func(item string) {
		if f(item) {
			count++
		}
	})
	return count
}

func (set *StrOrderSet) GroupByBool(f func(string) bool) (trueGroup *StrOrderSet, falseGroup *StrOrderSet) {
	trueGroup, falseGroup = NewStrOrderSet(0, set.cmp), NewStrOrderSet(0, set.cmp)
	set.ForEach(func(item string) {
		if f(item) {
			trueGroup.Append(item)
		} else {
			falseGroup.Append(item)
		}
	})
	return trueGroup, falseGroup
}

func (set *StrOrderSet) GroupByStr(f func(string) string) map[string]*StrOrderSet {
	groups := make(map[string]*StrOrderSet)
	set.ForEach(func(item string) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = NewStrOrderSet(0, set.cmp)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (set *StrOrderSet) GroupByInt(f func(string) int) map[int]*StrOrderSet {
	groups := make(map[int]*StrOrderSet)
	set.ForEach(func(item string) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = NewStrOrderSet(0, set.cmp)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

func (set *StrOrderSet) GroupBy(f func(string) interface{}) map[interface{}]*StrOrderSet {
	groups := make(map[interface{}]*StrOrderSet)
	set.ForEach(func(item string) {
		key := f(item)
		group := groups[key]
		if group == nil {
			group = NewStrOrderSet(0, set.cmp)
			groups[key] = group
		}
		group.Append(item)
	})
	return groups
}

// f: func(string) T
// return: []T
func (set *StrOrderSet) Map(f interface{}) interface{} {
	expected := "f should be func(string)T"
	ft := reflect.TypeOf(f)
	fVal := reflect.ValueOf(f)
	if ft.Kind() != reflect.Func {
		panic(expected)
	}
	if ft.NumIn() != 1 {
		panic(expected)
	}
	elemType := reflect.TypeOf(new(string)).Elem()
	if ft.In(0) != elemType {
		panic(expected)
	}
	if ft.NumOut() != 1 {
		panic(expected)
	}
	outType := ft.Out(0)
	result := reflect.MakeSlice(reflect.SliceOf(outType), 0, set.Len())
	set.ForEach(func(item string) {
		result = reflect.Append(result, fVal.Call([]reflect.Value{reflect.ValueOf(item)})[0])
	})
	return result.Interface()
}

// f: func(string) *T
//    func(string) (T, bool)
//    func(string) (T, error)
// return: []T
func (set *StrOrderSet) FilterMap(f interface{}) interface{} {
	expected := "f should be func(string) *T / func(string) (T, bool) / func(string) (T, error)"
	ft := reflect.TypeOf(f)
	fVal := reflect.ValueOf(f)
	if ft.Kind() != reflect.Func {
		panic(expected)
	}
	if ft.NumIn() != 1 {
		panic(expected)
	}
	in := ft.In(0)
	if in != reflect.TypeOf(new(string)).Elem() {
		panic(expected)
	}
	var outType reflect.Type
	var filter func([]reflect.Value) *reflect.Value
	if ft.NumOut() == 1 {
		// func(string) *T
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
			// func(string) (T, bool)
			filter = func(values []reflect.Value) *reflect.Value {
				if values[1].Interface().(bool) {
					return &values[0]
				}
				return nil
			}
		} else if checker.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			// func(string) (T, error)
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
	set.ForEach(func(item string) {
		ret := fVal.Call([]reflect.Value{reflect.ValueOf(item)})
		if val := filter(ret); val != nil {
			result = reflect.Append(result, *val)
		}
	})
	return result.Interface()
}

func (set *StrOrderSet) String() string {
	return fmt.Sprint(set.elementSequence)
}

func (set StrOrderSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *StrOrderSet) UnmarshalJSON(b []byte) error {
	return fmt.Errorf("unsupported")
}
