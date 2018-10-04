package plugin

import (
	"encoding/json"
	"fmt"
)

type ValueSet struct {
	elements        map[Value]uint32
	elementSequence []Value
}

func NewValueSet(capacity int) *ValueSet {
	set := new(ValueSet)
	if capacity > 0 {
		set.elements = make(map[Value]uint32, capacity)
		set.elementSequence = make([]Value, 0, capacity)
	} else {
		set.elements = make(map[Value]uint32)
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
	s := make([]Value, set.Len())
	for idx, item := range set.elementSequence {
		s[idx] = item
	}
	return s
}

// NOTICE: efficient but unsafe
func (set *ValueSet) ToSliceRef() []Value {
	return set.elementSequence
}

func (set *ValueSet) Put(key Value) {
	if _, ok := set.elements[key]; !ok {
		set.elements[key] = uint32(len(set.elementSequence))
		set.elementSequence = append(set.elementSequence, key)
	}
}

func (set *ValueSet) Clear() {
	set.elements = make(map[Value]uint32)
	set.elementSequence = set.elementSequence[:0]
}

func (set *ValueSet) Clone() *ValueSet {
	cloned := NewValueSet(set.Len())
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
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

// TODO keep order
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
	for _, item := range set.elementSequence {
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

func (set *ValueSet) Contains(key Value) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *ValueSet) ContainsAny(keys ...Value) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *ValueSet) ContainsAll(keys ...Value) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *ValueSet) String() string {
	return fmt.Sprint(set.elementSequence)
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
