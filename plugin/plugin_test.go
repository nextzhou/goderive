package plugin

import (
	"testing"

	"github.com/nextzhou/goderive/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseOptions(t *testing.T) {
	Convey("parse options", t, func() {
		Convey("trim space", func() {
			s := ""
			opts, err := ParseOptions(s)
			So(err, ShouldBeNil)
			So(opts.ExistingOption, ShouldBeEmpty)

			s = " \t "
			opts, err = ParseOptions(s)
			So(err, ShouldBeNil)
			So(opts.ExistingOption, ShouldBeEmpty)

			s = "  \t  flag1    ;  \t flag2 \t\n; key  = val1 , val2\n"
			opts, err = ParseOptions(s)
			So(err, ShouldBeNil)
			So(opts.Flags["flag1"], ShouldEqual, utils.TriBoolTrue)
			So(opts.Flags["flag2"], ShouldEqual, utils.TriBoolTrue)
			So(opts.Args["key"].Values, ShouldContain, Value("val1"))
			So(opts.Args["key"].Values, ShouldContain, Value("val2"))
		})
		Convey("flag", func() {
			s := "!flag"
			opts, err := ParseOptions(s)
			So(err, ShouldBeNil)
			So(opts.ExistingOption["flag"], ShouldEqual, OptionTypeFlag)
			So(opts.WithFlag("flag"), ShouldBeFalse)
			So(opts.WithNegativeFlag("flag"), ShouldBeTrue)
			So(opts.WithFlag("NotExistedFlag"), ShouldBeFalse)
			So(opts.WithNegativeFlag("NotExistedFlag"), ShouldBeFalse)

			s = "flag!"
			_, err = ParseOptions(s)
			So(err, ShouldBeError, `invalid flag "flag!"`)

			s = "flag;flag"
			opts, err = ParseOptions(s)
			So(err, ShouldBeError, `already existed flag "flag"`)
			s = "flag;!flag"
			opts, err = ParseOptions(s)
			So(err, ShouldBeError, `already existed flag "flag"`)
		})
		Convey("arg", func() {
			s := "key=val"
			opts, err := ParseOptions(s)
			So(err, ShouldBeNil)
			So(opts.ExistingOption["key"], ShouldEqual, OptionTypeArgKey)
			So(opts.GetValue("key").IsNil(), ShouldBeFalse)
			So(opts.GetValue("key").Str(), ShouldEqual, "val")

			s = "key=   val\n"
			opts, err = ParseOptions(s)
			So(err, ShouldBeNil)
			So(opts.GetValue("key").Str(), ShouldEqual, "val")

			s = "key;key=val"
			opts, err = ParseOptions(s)
			So(err, ShouldBeError, `already existed flag "key"`)
			s = "key=val1;key=val2"
			opts, err = ParseOptions(s)
			So(err, ShouldBeError, `already existed arg key "key"`)
			s = "val;key=val"
			_, err = ParseOptions(s)
			So(err, ShouldBeNil)
		})
	})

}

func TestIdent(t *testing.T) {
	Convey("ident validate", t, func() {
		f := utils.ValidateIdentName
		Convey("valid", func() {
			So(f("a"), ShouldBeTrue)
			So(f("_"), ShouldBeTrue)
			So(f("___"), ShouldBeTrue)
			So(f("_a"), ShouldBeTrue)
			So(f("A"), ShouldBeTrue)
			So(f("asdf"), ShouldBeTrue)
			So(f("a1"), ShouldBeTrue)
			So(f("a1___"), ShouldBeTrue)
			So(f("a1___1234_fla"), ShouldBeTrue)
		})
		Convey("invalid", func() {
			So(f(""), ShouldBeFalse)
			So(f(" "), ShouldBeFalse)
			So(f("\t"), ShouldBeFalse)
			So(f("哈哈"), ShouldBeFalse)
			So(f("!!!"), ShouldBeFalse)
			So(f("abc?"), ShouldBeFalse)
			So(f("abc def"), ShouldBeFalse)
			So(f(" a"), ShouldBeFalse)
			So(f("a "), ShouldBeFalse)
		})
	})
}
