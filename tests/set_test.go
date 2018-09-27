package examples

import (
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

		set.Put(4)
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
	})
}
