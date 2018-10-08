package tests

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntSet(t *testing.T) {
	Convey("int set", t, func() {
		set := NewIntSetFromSlice([]int{1, 2, 3})
		So(set.Contains(1), ShouldBeTrue)
		So(set.Contains(2), ShouldBeTrue)
		So(set.Contains(3), ShouldBeTrue)
		So(set.Contains(4), ShouldBeFalse)
		So(set.ContainsAny(9, 6, 3, 0), ShouldBeTrue)
		So(set.ContainsAny(9, 6, 0), ShouldBeFalse)
		So(set.ContainsAll(9, 6, 3, 0), ShouldBeFalse)
		So(set.ContainsAll(1, 2, 3), ShouldBeTrue)

		set.Remove(1)
		So(set.Contains(1), ShouldBeFalse)
		So(set.ContainsAll(2, 3), ShouldBeTrue)

		So(set.Len(), ShouldEqual, 2)
		So(set.IsEmpty(), ShouldBeFalse)

		set.Append(4)
		So(set.Len(), ShouldEqual, 3)
		So(set.Contains(4), ShouldBeTrue)

		even := set.Filter(func(i Int) bool { return i%2 == 0 })
		So(even.Len(), ShouldEqual, 2)
		So(even.ContainsAll(2, 4), ShouldBeTrue)
		So(even.Contains(3), ShouldBeFalse)

		set.Clear()
		So(set.Len(), ShouldEqual, 0)
		So(set.IsEmpty(), ShouldBeTrue)
		So(set.ContainsAny(1, 2, 3, 4), ShouldBeFalse)

		set = &IntSet{}
		err := json.Unmarshal([]byte(`[3,7,4,7]`), set)
		So(err, ShouldBeNil)
		So(set.ContainsAll(3, 7, 4), ShouldBeTrue)
		So(set.Len(), ShouldEqual, 3)
		found := set.FindBy(func(i Int) bool {
			return i > 5
		})
		So(found, ShouldNotBeNil)
		So(*found, ShouldEqual, 7)
		*found = 5
		So(set.Contains(7), ShouldBeTrue)
		So(set.Contains(5), ShouldBeFalse)
		found = set.FindBy(func(i Int) bool {
			return i < 0
		})
		So(found, ShouldBeNil)

		So(set.All(func(i Int) bool { return i >= 3 }), ShouldBeTrue)
		So(set.All(func(i Int) bool { return i > 3 }), ShouldBeFalse)
		So(set.Any(func(i Int) bool { return i%2 == 0 }), ShouldBeTrue)
		So(set.Any(func(i Int) bool { return i%5 == 0 }), ShouldBeFalse)

		So(set.CountBy(func(i Int) bool { return i%2 == 1 }), ShouldEqual, 2)
	})
}

func TestAppendOrderIntSet(t *testing.T) {
	Convey("int append set", t, func() {
		set := newIntOrderSetFromSlice([]int{1, 2, 3})
		So(set.Contains(1), ShouldBeTrue)
		So(set.Contains(2), ShouldBeTrue)
		So(set.Contains(3), ShouldBeTrue)
		So(set.Contains(4), ShouldBeFalse)
		So(set.ContainsAny(9, 6, 3, 0), ShouldBeTrue)
		So(set.ContainsAny(9, 6, 0), ShouldBeFalse)
		So(set.ContainsAll(9, 6, 3, 0), ShouldBeFalse)
		So(set.ContainsAll(1, 2, 3), ShouldBeTrue)
		So(set.String(), ShouldEqual, "[1 2 3]")

		set.Remove(1)
		So(set.Contains(1), ShouldBeFalse)
		So(set.ContainsAll(2, 3), ShouldBeTrue)
		So(set.String(), ShouldEqual, "[2 3]")

		So(set.Len(), ShouldEqual, 2)
		So(set.IsEmpty(), ShouldBeFalse)

		set.Append(4)
		So(set.Len(), ShouldEqual, 3)
		So(set.Contains(4), ShouldBeTrue)
		So(set.String(), ShouldEqual, "[2 3 4]")

		set.Append(3)
		So(set.String(), ShouldEqual, "[2 3 4]")
		set.Append(1)
		So(set.String(), ShouldEqual, "[2 3 4 1]")

		even := set.Filter(func(i Int) bool { return i%2 == 0 })
		So(even.Len(), ShouldEqual, 2)
		So(even.ContainsAll(2, 4), ShouldBeTrue)
		So(even.Contains(3), ShouldBeFalse)
		So(even.String(), ShouldEqual, "[2 4]")

		set.Clear()
		So(set.Len(), ShouldEqual, 0)
		So(set.IsEmpty(), ShouldBeTrue)
		So(set.ContainsAny(1, 2, 3, 4), ShouldBeFalse)
		So(set.String(), ShouldEqual, "[]")

		set = &intOrderSet{}
		err := json.Unmarshal([]byte(`[3,7,4,7,5]`), set)
		So(err, ShouldBeNil)
		So(set.ContainsAll(3, 7, 4, 5), ShouldBeTrue)
		So(set.Len(), ShouldEqual, 4)
		data, err := json.Marshal(set)
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, `[3,7,4,5]`)
		found := set.FindBy(func(i Int2) bool {
			return i%2 == 0
		})
		So(found, ShouldNotBeNil)
		So(*found, ShouldEqual, 4)
		*found = 6
		So(set.String(), ShouldEqual, `[3 7 4 5]`)
		So(fmt.Sprint(set.ToSlice()), ShouldEqual, `[3 7 4 5]`)
		So(fmt.Sprint(set.ToSliceRef()), ShouldEqual, `[3 7 4 5]`)
		So(set.CountBy(func(i int) bool { return i%2 == 1 }), ShouldEqual, 3)

		found = set.FindLastBy(func(i Int2) bool {
			return i < 5
		})
		So(found, ShouldNotBeNil)
		So(*found, ShouldEqual, 4)
		found = set.FindLastBy(func(_ Int2) bool {
			return false
		})
		So(found, ShouldBeNil)
	})
}

func TestKeyOrderIntSet(t *testing.T) {
	Convey("int key order set", t, func() {
		set := NewAscendingInt3Set(0)
		set.Append([]int{3, 8, 1, 5, 3, 5, 4}...)
		So(set.Len(), ShouldEqual, 5)
		So(set.String(), ShouldEqual, "[1 3 4 5 8]")
		So(fmt.Sprint(set.ToSlice()), ShouldEqual, "[1 3 4 5 8]")
		So(fmt.Sprint(set.ToSliceRef()), ShouldEqual, "[1 3 4 5 8]")
		set.Remove(5)
		So(set.Len(), ShouldEqual, 4)
		So(set.String(), ShouldEqual, "[1 3 4 8]")
		set.Append(2)
		So(set.String(), ShouldEqual, "[1 2 3 4 8]")
		union := set.Union(NewInt3SetFromSlice([]int{2, 4, 6, 8, 10}, func(i, j Int3) bool { return i < j }))
		So(union.String(), ShouldEqual, "[1 2 3 4 6 8 10]")

		descSet := NewDescendingInt3SetFromSlice(set.ToSlice())
		So(descSet.String(), ShouldEqual, "[8 4 3 2 1]")

		So(set.All(func(i int) bool { return i > 0 }), ShouldBeTrue)
		So(set.All(func(i int) bool { return i > 3 }), ShouldBeFalse)
		So(set.Any(func(i int) bool { return i > 3 }), ShouldBeTrue)
		So(set.Any(func(i int) bool { return i > 9 }), ShouldBeFalse)

		So(set.CountBy(func(i int) bool { return i%2 == 1 }), ShouldEqual, 2)
	})
}

func TestGroupBy(t *testing.T) {
	Convey("group by", t, func() {
		set := NewAscendingSSet(0)
		set.Append("bbbb", "abc", "bcd", "a", "defghi")
		containsA, notContainsA := set.GroupByBool(func(s string) bool { return strings.Contains(s, "a") })
		So(containsA.String(), ShouldEqual, `[a abc]`)
		So(notContainsA.String(), ShouldEqual, `[bbbb bcd defghi]`)

		lenGroups := set.GroupByInt(func(s string) int { return len(s) })
		So(lenGroups[1].String(), ShouldEqual, `[a]`)
		So(lenGroups[2].IsEmpty(), ShouldBeTrue)
		So(lenGroups[3].String(), ShouldEqual, `[abc bcd]`)
		So(lenGroups[4].String(), ShouldEqual, `[bbbb]`)
		So(lenGroups[6].String(), ShouldEqual, `[defghi]`)

		intialGroups := set.GroupByStr(func(s string) string { return string(s[0]) })
		So(intialGroups["a"].String(), ShouldEqual, `[a abc]`)
		So(intialGroups["b"].String(), ShouldEqual, `[bbbb bcd]`)

		lenOrSelfGroup := set.GroupBy(func(s string) interface{} {
			if len(s) <= 3 {
				return len(s)
			}
			return s
		})
		So(lenOrSelfGroup[1].String(), ShouldEqual, `[a]`)
		So(lenOrSelfGroup[3].String(), ShouldEqual, `[abc bcd]`)
		So(lenOrSelfGroup["bbbb"].String(), ShouldEqual, `[bbbb]`)
		So(lenOrSelfGroup["defghi"].String(), ShouldEqual, `[defghi]`)
	})
}
