package lines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_RawLinifyString(t *testing.T) {

	var (
		max   = 20
		sline = "This is a short line"
		lline = "This is a longer line in need of truncation, please"
		linll = "This is a longer lin\ne in need of truncat\nion, please"
	)

	Convey("", t, FailureContinues, func() {
		So(len(sline), ShouldEqual, max)
		So(RawLinifyString(sline, max), ShouldEqual, sline)
		So(len(lline), ShouldBeGreaterThan, max)
		So(RawLinifyString(lline, max), ShouldEqual, linll)
	})

}

func Test_LinifyString(t *testing.T) {

	var (
		max   = 20
		sline = "This is a short line"
		lline = "This is a longer line in need of truncation, please"
		linll = "This is a longer\nline in need of\ntruncation, please"
		lword = "areyoukiddingmerightnow"
		linlw = "areyoukiddingmeright\nnow\n"
	)

	Convey("", t, FailureContinues, func() {
		So(len(sline), ShouldEqual, max)
		So(LinifyString(sline, max), ShouldEqual, sline)

		So(len(lline), ShouldBeGreaterThan, max)
		So(LinifyString(lline, max), ShouldEqual, linll)

		So(len(lword), ShouldBeGreaterThan, max)
		So(LinifyString(lword, max), ShouldEqual, linlw)
	})

}
