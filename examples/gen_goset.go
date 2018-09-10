package examples


type IntSet map[Int]struct{}


func NewIntSet (capacity int) IntSet {
	if capacity > 0 {
		return make(map[Int]struct{}, capacity)
	}
	return make(map[Int]struct{})
}

func NewIntSetFromSlice(items ...Int) IntSet {
	set := make(map[Int]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

func (set IntSet) Extend(items ...Int) {
	for _, item := range items {
		set[item] = struct{}{}
	}
}

func (set IntSet) Len() int {
	return len(set)
}

func (set IntSet) Put(key Int) {
	set[key] = struct{}{}
}

func (set IntSet) Delete(key Int) {
	delete(set, key)
}

func (set IntSet) Contains(key Int) bool {
	_, ok := set[key]
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
