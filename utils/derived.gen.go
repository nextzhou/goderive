package utils

import (
	"encoding/json"
	"fmt"
)

type StrSet struct {
	elements map[Str]struct{}
}

func NewStrSet(capacity int) *StrSet {
	set := new(StrSet)
	if capacity > 0 {
		set.elements = make(map[Str]struct{}, capacity)
	} else {
		set.elements = make(map[Str]struct{})
	}
	return set
}

func NewStrSetFromSlice(items []Str) *StrSet {
	set := NewStrSet(len(items))
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func (set *StrSet) Extend(items ...Str) {
	for _, item := range items {
		set.Put(item)
	}
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

func (set *StrSet) ToSlice() []Str {
	if set == nil {
		return nil
	}
	s := make([]Str, 0, set.Len())
	set.ForEach(func(item Str) {
		s = append(s, item)
	})
	return s
}

func (set *StrSet) Put(key Str) {
	set.elements[key] = struct{}{}
}

func (set *StrSet) Clear() {
	set.elements = make(map[Str]struct{})
}

func (set *StrSet) Clone() *StrSet {
	cloned := NewStrSet(set.Len())
	for item := range set.elements {
		cloned.elements[item] = struct{}{}
	}
	return cloned
}

func (set *StrSet) Difference(another *StrSet) *StrSet {
	difference := NewStrSet(0)
	set.ForEach(func(item Str) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *StrSet) Equal(another *StrSet) bool {
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

func (set *StrSet) Intersect(another *StrSet) *StrSet {
	intersection := NewStrSet(0)
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

func (set *StrSet) Union(another *StrSet) *StrSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *StrSet) InPlaceUnion(another *StrSet) {
	another.ForEach(func(item Str) {
		set.Put(item)
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

func (set *StrSet) ForEach(f func(Str)) {
	if set.IsEmpty() {
		return
	}
	for item := range set.elements {
		f(item)
	}
}

func (set *StrSet) Filter(f func(Str) bool) *StrSet {
	result := NewStrSet(0)
	set.ForEach(func(item Str) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *StrSet) Remove(key Str) {
	delete(set.elements, key)
}

func (set *StrSet) Contains(key Str) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *StrSet) ContainsAny(keys ...Str) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *StrSet) ContainsAll(keys ...Str) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *StrSet) FindBy(f func(Str) bool) *Str {
	for item := range set.elements {
		if f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrSet) String() string {
	return fmt.Sprint(set.ToSlice())
}

func (set *StrSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *StrSet) UnmarshalJSON(b []byte) error {
	s := make([]Str, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewStrSetFromSlice(s)
	return nil
}
