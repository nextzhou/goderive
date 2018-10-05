package utils

import (
	"encoding/json"
	"fmt"
	"sort"
)

type StrSet struct {
	elements        map[Str]uint32
	elementSequence []Str
}

func NewStrSet(capacity int) *StrSet {
	set := new(StrSet)
	if capacity > 0 {
		set.elements = make(map[Str]uint32, capacity)
		set.elementSequence = make([]Str, 0, capacity)
	} else {
		set.elements = make(map[Str]uint32)
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
	s := make([]Str, set.Len())
	for idx, item := range set.elementSequence {
		s[idx] = item
	}
	return s
}

// NOTICE: efficient but unsafe
func (set *StrSet) ToSliceRef() []Str {
	return set.elementSequence
}

func (set *StrSet) Put(key Str) {
	if _, ok := set.elements[key]; !ok {
		set.elements[key] = uint32(len(set.elementSequence))
		set.elementSequence = append(set.elementSequence, key)
	}
}

func (set *StrSet) Clear() {
	set.elements = make(map[Str]uint32)
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

// TODO keep order
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
	for _, item := range set.elementSequence {
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

func (set *StrSet) DoUntil(f func(Str) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrSet) DoWhile(f func(Str) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrSet) FindBy(f func(Str) bool) *Str {
	for _, item := range set.elementSequence {
		if f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrSet) FindLastBy(f func(Str) bool) *Str {
	for i := set.Len() - 1; i >= 0; i-- {
		if item := set.elementSequence[i]; f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrSet) String() string {
	return fmt.Sprint(set.elementSequence)
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

type StrOrderSet struct {
	cmp             func(i, j Str2) bool
	elements        map[Str2]uint32
	elementSequence []Str2
}

func NewStrOrderSet(capacity int, cmp func(i, j Str2) bool) *StrOrderSet {
	set := new(StrOrderSet)
	if capacity > 0 {
		set.elements = make(map[Str2]uint32, capacity)
		set.elementSequence = make([]Str2, 0, capacity)
	} else {
		set.elements = make(map[Str2]uint32)
	}
	set.cmp = cmp
	return set
}

func NewStrOrderSetFromSlice(items []Str2, cmp func(i, j Str2) bool) *StrOrderSet {
	set := NewStrOrderSet(len(items), cmp)
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func NewAscendingStrOrderSet(capacity int) *StrOrderSet {
	return NewStrOrderSet(capacity, func(i, j Str2) bool { return i < j })
}

func NewDescendingStrOrderSet(capacity int) *StrOrderSet {
	return NewStrOrderSet(capacity, func(i, j Str2) bool { return i > j })
}

func NewAscendingStrOrderSetFromSlice(items []Str2) *StrOrderSet {
	return NewStrOrderSetFromSlice(items, func(i, j Str2) bool { return i < j })
}

func NewDescendingStrOrderSetFromSlice(items []Str2) *StrOrderSet {
	return NewStrOrderSetFromSlice(items, func(i, j Str2) bool { return i > j })
}

func (set *StrOrderSet) Extend(items ...Str2) {
	for _, item := range items {
		set.Put(item)
	}
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

func (set *StrOrderSet) ToSlice() []Str2 {
	if set == nil {
		return nil
	}
	s := make([]Str2, 0, set.Len())
	set.ForEach(func(item Str2) {
		s = append(s, item)
	})
	return s
}

func (set *StrOrderSet) Put(key Str2) {
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

func (set *StrOrderSet) Clear() {
	set.elements = make(map[Str2]uint32)
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
	set.ForEach(func(item Str2) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *StrOrderSet) Equal(another *StrOrderSet) bool {
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

func (set *StrOrderSet) Intersect(another *StrOrderSet) *StrOrderSet {
	intersection := NewStrOrderSet(0, set.cmp)
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

func (set *StrOrderSet) Union(another *StrOrderSet) *StrOrderSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *StrOrderSet) InPlaceUnion(another *StrOrderSet) {
	another.ForEach(func(item Str2) {
		set.Put(item)
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

func (set *StrOrderSet) ForEach(f func(Str2)) {
	if set.IsEmpty() {
		return
	}
	for _, item := range set.elementSequence {
		f(item)
	}
}

func (set *StrOrderSet) Filter(f func(Str2) bool) *StrOrderSet {
	result := NewStrOrderSet(0, set.cmp)
	set.ForEach(func(item Str2) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *StrOrderSet) Remove(key Str2) {
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

func (set *StrOrderSet) Contains(key Str2) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *StrOrderSet) ContainsAny(keys ...Str2) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *StrOrderSet) ContainsAll(keys ...Str2) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *StrOrderSet) DoUntil(f func(Str2) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrOrderSet) DoWhile(f func(Str2) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *StrOrderSet) FindBy(f func(Str2) bool) *Str2 {
	for _, item := range set.elementSequence {
		if f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrOrderSet) FindLastBy(f func(Str2) bool) *Str2 {
	for i := set.Len() - 1; i >= 0; i-- {
		if item := set.elementSequence[i]; f(item) {
			return &item
		}
	}
	return nil
}

func (set *StrOrderSet) String() string {
	return fmt.Sprint(set.elementSequence)
}

func (set *StrOrderSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *StrOrderSet) UnmarshalJSON(b []byte) error {
	return fmt.Errorf("unsupported")
}
