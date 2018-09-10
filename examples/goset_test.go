package examples

import "testing"

func TestBaseType(t *testing.T) {
	set := NewIntSet(0)
	if set.Len() != 0 {
		t.Errorf("expected 0, got %d", set.Len())
	}
	set = NewIntSet(10)
	if set.Len() != 0 {
		t.Errorf("expected 0, got %d", set.Len())
	}
	set.Put(1)
	if set.Len() != 1 {
		t.Errorf("expected 1, got %d", set.Len())
	}
	set.Put(2)
	if set.Len() != 2 {
		t.Errorf("expected 2, got %d", set.Len())
	}
	set.Put(1)
	if set.Len() != 2 {
		t.Errorf("expected 2, got %d", set.Len())
	}
	if !set.Contains(1) {
		t.Errorf("set should contains 1")
	}
	if !set.Contains(2) {
		t.Errorf("set should contains 2")
	}
	if set.Contains(3) {
		t.Errorf("set shouldn't contains 3")
	}
	if set.ContainsAll(1,2,3) {
		t.Errorf("set shouldn't contains 3")
	}
	if !set.ContainsAll(1,2) {
		t.Errorf("set shouldn contains 1,2")
	}
	if set.ContainsAny(3,4,5) {
		t.Errorf("set shouldn't contains anyone in 3,4,5")
	}
	if !set.ContainsAny(2,3,4,5) {
		t.Errorf("set should contains someone in 2,3,4,5")
	}
	if !set.ContainsAny(3,2,4,5) {
		t.Errorf("set should contains someone in 3,2,4,5")
	}
	if !set.ContainsAny(3,4,2,5) {
		t.Errorf("set should contains someone in 3,4,2,5")
	}
	if !set.ContainsAny(3,4,5,2) {
		t.Errorf("set should contains someone in 3,4,5,2")
	}
	set.Delete(1)
	if set.Len() != 1 {
		t.Errorf("expected 1, got %d", set.Len())
	}
	if set.ContainsAll(1,2) {
		t.Errorf("set shouldn't contains 1,2")
	}
	if !set.ContainsAny(1,2) {
		t.Errorf("set shouldn contains someone in 1,2")
	}
	set.Extend(3,4,5)
	if set.Len() != 4 {
		t.Errorf("expected 4, got %d", set.Len())
	}
}

