package examples

import (
	"encoding/json"
	"fmt"
)

type IntSet struct {
	elements map[Int]struct{}
}

func NewIntSet(capacity int) *IntSet {
	set := new(IntSet)
	if capacity > 0 {
		set.elements = make(map[Int]struct{}, capacity)
	} else {
		set.elements = make(map[Int]struct{})
	}
	return set
}

func NewIntSetFromSlice(items []Int) *IntSet {
	set := new(IntSet)
	set.elements = make(map[Int]struct{}, len(items))
	for _, item := range items {
		set.elements[item] = struct{}{}
	}
	return set
}

func (set *IntSet) Extend(items ...Int) {
	for _, item := range items {
		set.elements[item] = struct{}{}
	}
}

func (set *IntSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *IntSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *IntSet) ToSlice() []Int {
	if set == nil {
		return nil
	}
	s := make([]Int, 0, set.Len())
	set.ForEach(func(item Int) {
		s = append(s, item)
	})
	return s
}

func (set *IntSet) Put(key Int) {
	set.elements[key] = struct{}{}
}

func (set *IntSet) Clear() {
	set.elements = make(map[Int]struct{})
}

func (set *IntSet) Clone() *IntSet {
	cloned := NewIntSet(set.Len())
	for item := range set.elements {
		cloned.Put(item)
	}
	return cloned
}

func (set *IntSet) Difference(another *IntSet) *IntSet {
	difference := NewIntSet(0)
	for item := range set.elements {
		if !another.Contains(item) {
			difference.Put(item)
		}
	}
	return difference
}

func (set *IntSet) Equal(another *IntSet) bool {
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

func (set *IntSet) Intersect(another *IntSet) *IntSet {
	intersection := NewIntSet(0)
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

func (set *IntSet) Union(another *IntSet) *IntSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *IntSet) InPlaceUnion(another *IntSet) {
	another.ForEach(func(item Int) {
		set.Put(item)
	})
}

func (set *IntSet) IsProperSubsetOf(another *IntSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *IntSet) IsProperSupersetOf(another *IntSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *IntSet) IsSubsetOf(another *IntSet) bool {
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

func (set *IntSet) IsSupersetOf(another *IntSet) bool {
	return another.IsSubsetOf(set)
}

func (set *IntSet) ForEach(f func(Int)) {
	if set.IsEmpty() {
		return
	}
	for item := range set.elements {
		f(item)
	}
}

func (set *IntSet) Filter(f func(Int) bool) *IntSet {
	result := NewIntSet(0)
	set.ForEach(func(item Int) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set IntSet) Remove(key Int) {
	delete(set.elements, key)
}

func (set IntSet) Contains(key Int) bool {
	_, ok := set.elements[key]
	return ok
}

func (set IntSet) ContainsAny(keys ...Int) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set IntSet) ContainsAll(keys ...Int) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *IntSet) String() string {
	return fmt.Sprint(set.ToSlice())
}

func (set *IntSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *IntSet) UnmarshalJSON(b []byte) error {
	s := make([]Int, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewIntSetFromSlice(s)
	return nil
}
