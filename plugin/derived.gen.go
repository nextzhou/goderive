package plugin

import (
	"encoding/json"
	"fmt"
)

type ValueSet struct {
	elements map[Value]struct{}
}

func NewValueSet(capacity int) *ValueSet {
	set := new(ValueSet)
	if capacity > 0 {
		set.elements = make(map[Value]struct{}, capacity)
	} else {
		set.elements = make(map[Value]struct{})
	}
	return set
}

func NewValueSetFromSlice(items []Value) *ValueSet {
	set := NewValueSet(len(items))
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func (set *ValueSet) Extend(items ...Value) {
	for _, item := range items {
		set.Put(item)
	}
}

func (set *ValueSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *ValueSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *ValueSet) ToSlice() []Value {
	if set == nil {
		return nil
	}
	s := make([]Value, 0, set.Len())
	set.ForEach(func(item Value) {
		s = append(s, item)
	})
	return s
}

func (set *ValueSet) Put(key Value) {
	set.elements[key] = struct{}{}
}

func (set *ValueSet) Clear() {
	set.elements = make(map[Value]struct{})
}

func (set *ValueSet) Clone() *ValueSet {
	cloned := NewValueSet(set.Len())
	for item := range set.elements {
		cloned.elements[item] = struct{}{}
	}
	return cloned
}

func (set *ValueSet) Difference(another *ValueSet) *ValueSet {
	difference := NewValueSet(0)
	set.ForEach(func(item Value) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *ValueSet) Equal(another *ValueSet) bool {
	if set.Len() != another.Len() {
		return false
	}
	for item := range set.elements {
		if !another.Contains(item) {
			return false
		}
	}
	return true
}

func (set *ValueSet) Intersect(another *ValueSet) *ValueSet {
	intersection := NewValueSet(0)
	if set.Len() < another.Len() {
		for item := range set.elements {
			if another.Contains(item) {
				intersection.Put(item)
			}
		}
	} else {
		for item := range another.elements {
			if set.Contains(item) {
				intersection.Put(item)
			}
		}
	}
	return intersection
}

func (set *ValueSet) Union(another *ValueSet) *ValueSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *ValueSet) InPlaceUnion(another *ValueSet) {
	another.ForEach(func(item Value) {
		set.Put(item)
	})
}

func (set *ValueSet) IsProperSubsetOf(another *ValueSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *ValueSet) IsProperSupersetOf(another *ValueSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *ValueSet) IsSubsetOf(another *ValueSet) bool {
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

func (set *ValueSet) IsSupersetOf(another *ValueSet) bool {
	return another.IsSubsetOf(set)
}

func (set *ValueSet) ForEach(f func(Value)) {
	if set.IsEmpty() {
		return
	}
	for item := range set.elements {
		f(item)
	}
}

func (set *ValueSet) Filter(f func(Value) bool) *ValueSet {
	result := NewValueSet(0)
	set.ForEach(func(item Value) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *ValueSet) Remove(key Value) {
	delete(set.elements, key)
}

func (set ValueSet) Contains(key Value) bool {
	_, ok := set.elements[key]
	return ok
}

func (set ValueSet) ContainsAny(keys ...Value) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set ValueSet) ContainsAll(keys ...Value) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *ValueSet) String() string {
	return fmt.Sprint(set.ToSlice())
}

func (set *ValueSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *ValueSet) UnmarshalJSON(b []byte) error {
	s := make([]Value, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewValueSetFromSlice(s)
	return nil
}
