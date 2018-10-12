package tests

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	Convey("access-get", t, func() {
		var a *AA
		So(func() { _ = a.b.c.Def }, ShouldPanic) // nil pointer
		So(a, ShouldBeNil)
		So(a.getB(), ShouldBeNil)
		So(a.getB().getC(), ShouldBeNil)
		So(a.getB().getC().GetDef(), ShouldBeNil)
		So(a.GetB().GetC().getAbc(), ShouldBeZeroValue)
		a = &AA{c: &c{}}
		So(a.getC(), ShouldNotBeNil)
		a.c.abc = "some value"
		So(a.getC().getAbc(), ShouldEqual, "some value")
	})
}
