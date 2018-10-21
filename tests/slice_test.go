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
	})
}
