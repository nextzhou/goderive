package tests

import (
	"json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntSlice(t *testing.T) {
	Convey("int slice", t, func() {
		var s *IntSlice
		So(s, ShouldBeNil)
		So(s.Len(), ShouldEqual, 0)

		s = NewIntSlice(0)
		So(s.Len(), ShouldEqual, 0)
		So(s.String(), ShouldEqual, "[]")

		s.Append(1, 2, 3)
		So(s.Len(), ShouldEqual, 3)
		So(s.String(), ShouldEqual, "[1 2 3]")

		cloned := s.Clone()
		cloned.Append(4, 5, 6)
		So(cloned.Len(), ShouldEqual, 6)
		So(cloned.String(), ShouldEqual, "[1 2 3 4 5 6]")
		So(s.Len(), ShouldEqual, 3)
		So(s.String(), ShouldEqual, "[1 2 3]")

		s.Append(2)
		So(s.Len(), ShouldEqual, 4)
		So(s.String(), ShouldEqual, "[1 2 3 2]")

		s.ToSlice()[0] = 5
		So(s.String(), ShouldEqual, "[1 2 3 2]")
		s.ToSliceRef()[0] = 5
		So(s.String(), ShouldEqual, "[5 2 3 2]")

		j, err := json.Marshal(s)
		So(err, ShouldBeNil)
		So(string(j), ShouldEqual, "[5,2,3,2]")
		j, err = json.Marshal(*s)
		So(err, ShouldBeNil)
		So(string(j), ShouldEqual, "[5,2,3,2]")

		j = []byte("[3,2,1]")
		err = json.Unmarshal(j, s)
		So(err, ShouldBeNil)
		So(s.String(), ShouldEqual, "[3 2 1]")

		Convey("insert", func() {
			s := NewIntSliceFromSlice([]int{1, 3, 5})

			// "-1" equal "s.Len()-1"
			s.Insert(-1, 4)
			So(s.String(), ShouldEqual, "[1 3 4 5]")

			s.Insert(1, 2)
			So(s.String(), ShouldEqual, "[1 2 3 4 5]")

			s.Insert(1, []int{1, 1, 2, 2}...)
			So(s.String(), ShouldEqual, "[1 1 1 2 2 2 3 4 5]")

			s.Insert(0, []int{0, 0, 0}...)
			So(s.String(), ShouldEqual, "[0 0 0 1 1 1 2 2 2 3 4 5]")
		})

		Convey("remove", func() {
			s := NewIntSliceFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8})
			So(func() { s.Remove(10) }, ShouldPanic)

			// remove first element
			s.Remove(0)
			So(s.String(), ShouldEqual, "[2 3 4 5 6 7 8]")

			// remove last element
			s.Remove(-1)
			So(s.String(), ShouldEqual, "[2 3 4 5 6 7]")

			// remove last 3 elements
			s.RemoveFrom(-3)
			So(s.String(), ShouldEqual, "[2 3 4]")

			s = NewIntSliceFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8})
			s.RemoveTo(2)
			So(s.String(), ShouldEqual, "[4 5 6 7 8]")

			s.RemoveRange(1, -2)
			So(s.String(), ShouldEqual, "[4 8]")
		})

		Convey("filter", func() {
			s := NewIntSliceFromSlice([]int{1, 2, 3, 4, 5})
			So(s.Filter(func(i int) bool { return false }).String(), ShouldEqual, "[]")
			So(s.Filter(func(i int) bool { return true }).String(), ShouldEqual, "[1 2 3 4 5]")
			So(s.Filter(func(i int) bool { return i%2 == 1 }).String(), ShouldEqual, "[1 3 5]")

			So(s.String(), ShouldEqual, "[1 2 3 4 5]")
		})

		Convey("concat", func() {
			s := NewIntSliceFromSlice([]int{1, 2, 3, 4, 5})

			So(s.Concat(nil).String(), ShouldEqual, "[1 2 3 4 5]")
			So(s.Concat(NewIntSlice(0)).String(), ShouldEqual, "[1 2 3 4 5]")

			So(s.Concat(NewIntSliceFromSlice([]int{4, 5})).String(), ShouldEqual, "[1 2 3 4 5 4 5]")

			So(s.String(), ShouldEqual, "[1 2 3 4 5]")

			s.InPlaceConcat(nil)
			So(s.String(), ShouldEqual, "[1 2 3 4 5]")

			s.InPlaceConcat(NewIntSliceFromSlice([]int{6, 7, 8}))
			So(s.String(), ShouldEqual, "[1 2 3 4 5 6 7 8]")
		})
	})
}
