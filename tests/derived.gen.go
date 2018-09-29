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
	set := NewIntSet(len(items))
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func (set *IntSet) Extend(items ...Int) {
	for _, item := range items {
		set.Put(item)
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
		cloned.elements[item] = struct{}{}
	}
	return cloned
}

func (set *IntSet) Difference(another *IntSet) *IntSet {
	difference := NewIntSet(0)
	set.ForEach(func(item Int) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
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

func (set *IntSet) Remove(key Int) {
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

type intOrderSet struct {
	elements        map[Int2]uint32
	elementSequence []Int2
}

func NewIntOrderSet(capacity int) *intOrderSet {
	set := new(intOrderSet)
	if capacity > 0 {
		set.elements = make(map[Int2]uint32, capacity)
		set.elementSequence = make([]Int2, 0, capacity)
	} else {
		set.elements = make(map[Int2]uint32)
	}
	return set
}

func NewIntOrderSetFromSlice(items []Int2) *intOrderSet {
	set := NewIntOrderSet(len(items))
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func (set *intOrderSet) Extend(items ...Int2) {
	for _, item := range items {
		set.Put(item)
	}
}

func (set *intOrderSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *intOrderSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *intOrderSet) ToSlice() []Int2 {
	if set == nil {
		return nil
	}
	s := make([]Int2, set.Len())
	for idx, item := range set.elementSequence {
		s[idx] = item
	}
	return s
}

// NOTICE: efficient but unsafe
func (set *intOrderSet) ToSliceRef() []Int2 {
	return set.elementSequence
}

func (set *intOrderSet) Put(key Int2) {
	if _, ok := set.elements[key]; !ok {
		set.elements[key] = uint32(len(set.elementSequence))
		set.elementSequence = append(set.elementSequence, key)
	}
}

func (set *intOrderSet) Clear() {
	set.elements = make(map[Int2]uint32)
	set.elementSequence = set.elementSequence[:0]
}

func (set *intOrderSet) Clone() *intOrderSet {
	cloned := NewIntOrderSet(set.Len())
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
	}
	return cloned
}

func (set *intOrderSet) Difference(another *intOrderSet) *intOrderSet {
	difference := NewIntOrderSet(0)
	set.ForEach(func(item Int2) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *intOrderSet) Equal(another *intOrderSet) bool {
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
func (set *intOrderSet) Intersect(another *intOrderSet) *intOrderSet {
	intersection := NewIntOrderSet(0)
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

func (set *intOrderSet) Union(another *intOrderSet) *intOrderSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *intOrderSet) InPlaceUnion(another *intOrderSet) {
	another.ForEach(func(item Int2) {
		set.Put(item)
	})
}

func (set *intOrderSet) IsProperSubsetOf(another *intOrderSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *intOrderSet) IsProperSupersetOf(another *intOrderSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *intOrderSet) IsSubsetOf(another *intOrderSet) bool {
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

func (set *intOrderSet) IsSupersetOf(another *intOrderSet) bool {
	return another.IsSubsetOf(set)
}

func (set *intOrderSet) ForEach(f func(Int2)) {
	if set.IsEmpty() {
		return
	}
	for _, item := range set.elementSequence {
		f(item)
	}
}

func (set *intOrderSet) Filter(f func(Int2) bool) *intOrderSet {
	result := NewIntOrderSet(0)
	set.ForEach(func(item Int2) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *intOrderSet) Remove(key Int2) {
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

func (set intOrderSet) Contains(key Int2) bool {
	_, ok := set.elements[key]
	return ok
}

func (set intOrderSet) ContainsAny(keys ...Int2) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set intOrderSet) ContainsAll(keys ...Int2) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *intOrderSet) String() string {
	return fmt.Sprint(set.elementSequence)
}

func (set *intOrderSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *intOrderSet) UnmarshalJSON(b []byte) error {
	s := make([]Int2, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewIntOrderSetFromSlice(s)
	return nil
}
