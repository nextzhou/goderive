package plugin

import (
	"encoding/json"
	"fmt"
)

type PluginSet struct {
	elements        map[Plugin]uint32
	elementSequence []Plugin
}

func NewPluginSet(capacity int) *PluginSet {
	set := new(PluginSet)
	if capacity > 0 {
		set.elements = make(map[Plugin]uint32, capacity)
		set.elementSequence = make([]Plugin, 0, capacity)
	} else {
		set.elements = make(map[Plugin]uint32)
	}
	return set
}

func NewPluginSetFromSlice(items []Plugin) *PluginSet {
	set := NewPluginSet(len(items))
	for _, item := range items {
		set.Put(item)
	}
	return set
}

func (set *PluginSet) Extend(items ...Plugin) {
	for _, item := range items {
		set.Put(item)
	}
}

func (set *PluginSet) Len() int {
	if set == nil {
		return 0
	}
	return len(set.elements)
}

func (set *PluginSet) IsEmpty() bool {
	return set.Len() == 0
}

func (set *PluginSet) ToSlice() []Plugin {
	if set == nil {
		return nil
	}
	s := make([]Plugin, set.Len())
	for idx, item := range set.elementSequence {
		s[idx] = item
	}
	return s
}

// NOTICE: efficient but unsafe
func (set *PluginSet) ToSliceRef() []Plugin {
	return set.elementSequence
}

func (set *PluginSet) Put(key Plugin) {
	if _, ok := set.elements[key]; !ok {
		set.elements[key] = uint32(len(set.elementSequence))
		set.elementSequence = append(set.elementSequence, key)
	}
}

func (set *PluginSet) Clear() {
	set.elements = make(map[Plugin]uint32)
	set.elementSequence = set.elementSequence[:0]
}

func (set *PluginSet) Clone() *PluginSet {
	cloned := NewPluginSet(set.Len())
	for idx, item := range set.elementSequence {
		cloned.elements[item] = uint32(idx)
		cloned.elementSequence = append(cloned.elementSequence, item)
	}
	return cloned
}

func (set *PluginSet) Difference(another *PluginSet) *PluginSet {
	difference := NewPluginSet(0)
	set.ForEach(func(item Plugin) {
		if !another.Contains(item) {
			difference.Put(item)
		}
	})
	return difference
}

func (set *PluginSet) Equal(another *PluginSet) bool {
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
func (set *PluginSet) Intersect(another *PluginSet) *PluginSet {
	intersection := NewPluginSet(0)
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

func (set *PluginSet) Union(another *PluginSet) *PluginSet {
	union := set.Clone()
	union.InPlaceUnion(another)
	return union
}

func (set *PluginSet) InPlaceUnion(another *PluginSet) {
	another.ForEach(func(item Plugin) {
		set.Put(item)
	})
}

func (set *PluginSet) IsProperSubsetOf(another *PluginSet) bool {
	return !set.Equal(another) && set.IsSubsetOf(another)
}

func (set *PluginSet) IsProperSupersetOf(another *PluginSet) bool {
	return !set.Equal(another) && set.IsSupersetOf(another)
}

func (set *PluginSet) IsSubsetOf(another *PluginSet) bool {
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

func (set *PluginSet) IsSupersetOf(another *PluginSet) bool {
	return another.IsSubsetOf(set)
}

func (set *PluginSet) ForEach(f func(Plugin)) {
	if set.IsEmpty() {
		return
	}
	for _, item := range set.elementSequence {
		f(item)
	}
}

func (set *PluginSet) Filter(f func(Plugin) bool) *PluginSet {
	result := NewPluginSet(0)
	set.ForEach(func(item Plugin) {
		if f(item) {
			result.Put(item)
		}
	})
	return result
}

func (set *PluginSet) Remove(key Plugin) {
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

func (set *PluginSet) Contains(key Plugin) bool {
	_, ok := set.elements[key]
	return ok
}

func (set *PluginSet) ContainsAny(keys ...Plugin) bool {
	for _, key := range keys {
		if set.Contains(key) {
			return true
		}
	}
	return false
}

func (set *PluginSet) ContainsAll(keys ...Plugin) bool {
	for _, key := range keys {
		if !set.Contains(key) {
			return false
		}
	}
	return true
}

func (set *PluginSet) DoUntil(f func(Plugin) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *PluginSet) DoWhile(f func(Plugin) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *PluginSet) FindBy(f func(Plugin) bool) *Plugin {
	for _, item := range set.elementSequence {
		if f(item) {
			return &item
		}
	}
	return nil
}

func (set *PluginSet) FindLastBy(f func(Plugin) bool) *Plugin {
	for i := set.Len() - 1; i >= 0; i-- {
		if item := set.elementSequence[i]; f(item) {
			return &item
		}
	}
	return nil
}

func (set *PluginSet) String() string {
	return fmt.Sprint(set.elementSequence)
}

func (set *PluginSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.ToSlice())
}

func (set *PluginSet) UnmarshalJSON(b []byte) error {
	s := make([]Plugin, 0)
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*set = *NewPluginSetFromSlice(s)
	return nil
}

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

func (set *ValueSet) DoUntil(f func(Value) bool) int {
	for idx, item := range set.elementSequence {
		if f(item) {
			return idx
		}
	}
	return -1
}

func (set *ValueSet) DoWhile(f func(Value) bool) int {
	for idx, item := range set.elementSequence {
		if !f(item) {
			return idx
		}
	}
	return -1
}

func (set *ValueSet) FindBy(f func(Value) bool) *Value {
	for _, item := range set.elementSequence {
		if f(item) {
			return &item
		}
	}
	return nil
}

func (set *ValueSet) FindLastBy(f func(Value) bool) *Value {
	for i := set.Len() - 1; i >= 0; i-- {
		if item := set.elementSequence[i]; f(item) {
			return &item
		}
	}
	return nil
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
