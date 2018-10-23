package tests

import (
	"json"
	"testing"

	"strconv"

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

		Convey("for each", func() {
			var s *IntSlice
			sum := 0
			s.ForEach(func(i int) { sum += i })
			So(sum, ShouldEqual, 0)

			s = NewIntSliceFromSlice([]int{1, 2, 3})
			s.ForEach(func(i int) { sum += i })
			So(sum, ShouldEqual, 6)

			s.Clear()
			s.ForEach(func(i int) { sum += i })
			So(sum, ShouldEqual, 6)
		})

		Convey("index", func() {
			ss := []int{1, 2, 3, 4, 5}
			s := NewIntSliceFromSlice(ss)

			for idx := range ss {
				So(*s.Index(idx), ShouldEqual, ss[idx])
			}

			So(*s.Index(-1), ShouldEqual, 5)

			*s.Index(0) = 100
			So(s.String(), ShouldEqual, "[100 2 3 4 5]")

			ss = s.IndexRange(1, -1)
			So(NewIntSliceFromSlice(ss).String(), ShouldEqual, "[2 3 4]")

			ss = s.IndexFrom(2)
			So(NewIntSliceFromSlice(ss).String(), ShouldEqual, "[3 4 5]")

			ss = s.IndexTo(3)
			So(NewIntSliceFromSlice(ss).String(), ShouldEqual, "[100 2 3]")
		})

		Convey("find", func() {
			s := NewIntSliceFromSlice([]int{1, 2, 3, 4, 5})
			So(s.Find(100), ShouldEqual, -1)
			So(s.FindBy(func(i int) bool { return true }), ShouldEqual, 0)
			So(s.FindBy(func(i int) bool { return false }), ShouldEqual, -1)

			So(s.Find(3), ShouldEqual, 2)
			So(s.FindBy(func(i int) bool { return i%2 == 0 }), ShouldEqual, 1)

			s.Append(4, 3, 2, 1)
			So(s.Find(3), ShouldEqual, 2)
			So(s.FindLast(3), ShouldEqual, 6)
			So(s.FindBy(func(i int) bool { return i%2 == 0 }), ShouldEqual, 1)
			So(s.FindLastBy(func(i int) bool { return i%2 == 0 }), ShouldEqual, 7)
		})

		Convey("count", func() {
			var s *IntSlice
			So(s.Count(1), ShouldEqual, 0)
			So(s.CountBy(func(i int) bool { return true }), ShouldEqual, 0)
			So(s.CountBy(func(i int) bool { return false }), ShouldEqual, 0)

			s = NewIntSliceFromSlice([]int{1, 2, 3, 4, 5})
			So(s.Count(0), ShouldEqual, 0)
			So(s.Count(1), ShouldEqual, 1)
			So(s.Count(2), ShouldEqual, 1)
			s.Append(2)
			So(s.Count(2), ShouldEqual, 2)

			So(s.CountBy(func(i int) bool { return true }), ShouldEqual, s.Len())
			So(s.CountBy(func(i int) bool { return false }), ShouldEqual, 0)
			So(s.CountBy(func(i int) bool { return i%2 == 1 }), ShouldEqual, 3)
		})

		Convey("group", func() {
			Convey("group by bool", func() {
				s := NewIntSliceFromSlice([]int{1, 2, 3, 4, 5})
				odd, even := s.GroupByBool(func(i int) bool { return i%2 == 1 })
				So(odd.String(), ShouldEqual, "[1 3 5]")
				So(even.String(), ShouldEqual, "[2 4]")

				s = nil
				odd, even = s.GroupByBool(func(i int) bool { return i%2 == 1 })
				So(odd.IsEmpty(), ShouldBeTrue)
				So(even.IsEmpty(), ShouldBeTrue)
			})

			Convey("group by int", func() {
				s := NewIntSliceFromSlice([]int{1, 2, 3, 4, 5})
				groups := s.GroupByInt(func(i int) int { return i % 3 })
				So(groups[0].String(), ShouldEqual, "[3]")
				So(groups[1].String(), ShouldEqual, "[1 4]")
				So(groups[2].String(), ShouldEqual, "[2 5]")
				So(groups[3].IsEmpty(), ShouldBeTrue)
			})

			Convey("group by str", func() {
				s := NewIntSliceFromSlice([]int{1, 12, 2, 21, 22, 3})
				groups := s.GroupByStr(func(i int) string {
					return string(strconv.Itoa(i)[0])
				})
				So(groups[""].IsEmpty(), ShouldBeTrue)
				So(groups["0"].IsEmpty(), ShouldBeTrue)
				So(groups["12"].IsEmpty(), ShouldBeTrue)
				So(groups["1"].String(), ShouldEqual, "[1 12]")
				So(groups["2"].String(), ShouldEqual, "[2 21 22]")
				So(groups["3"].String(), ShouldEqual, "[3]")
			})

			Convey("group by interface", func() {
				s := NewIntSliceFromSlice([]int{1, 12, 2, 21, 22, 3})
				groups := s.GroupBy(func(i int) interface{} {
					s := strconv.Itoa(i)
					if len(s) == 1 {
						return i
					}
					return s[:1]
				})
				So(groups[1].String(), ShouldEqual, "[1]")
				So(groups["1"].String(), ShouldEqual, "[12]")
				So(groups[2].String(), ShouldEqual, "[2]")
				So(groups["2"].String(), ShouldEqual, "[21 22]")
				So(groups[3].String(), ShouldEqual, "[3]")
				So(groups["3"].IsEmpty(), ShouldBeTrue)
			})
		})
	})
}
