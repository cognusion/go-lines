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

	Convey("When a series of lines are RawLinified, the results are expected.", t, FailureContinues, func() {
		So(len(sline), ShouldEqual, max)
		So(RawLinifyString(sline, max), ShouldEqual, sline)
		So(len(lline), ShouldBeGreaterThan, max)
		So(RawLinifyString(lline, max), ShouldEqual, linll)
	})

}

func Test_LinifyString(t *testing.T) {

	var (
		max    = 20
		sline  = "This is a short line"
		lline  = "This is a longer line in need of truncation, please"
		linll  = "This is a longer\nline in need of\ntruncation, please"
		lword  = "areyoukiddingmerightnow"
		linlw  = "areyoukiddingmeright\nnow\n"
		leline = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	)

	Convey("When a series of lines are Linified, the results are expected.", t, FailureContinues, func() {
		So(len(sline), ShouldEqual, max)
		So(LinifyString(sline, max), ShouldEqual, sline)

		So(len(lline), ShouldBeGreaterThan, max)
		So(LinifyString(lline, max), ShouldEqual, linll)

		So(len(lword), ShouldBeGreaterThan, max)
		So(LinifyString(lword, max), ShouldEqual, linlw)

		So(len(leline), ShouldBeGreaterThan, max)
		So(len(LinifyString(leline, max)), ShouldEqual, 445)
	})
}

func Benchmark_LongLinifyString(b *testing.B) {
	var (
		lline = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
		max   = 20
	)

	for b.Loop() {
		LinifyString(lline, max)
	}
}

func Benchmark_MidLinifyString(b *testing.B) {
	var (
		lline = "This is a longer line in need of truncation, please"
		max   = 20
	)

	for b.Loop() {
		LinifyString(lline, max)
	}
}

func Benchmark_ShortLinifyString(b *testing.B) {
	var (
		lline = "This is a short line"
		max   = 20
	)

	for b.Loop() {
		LinifyString(lline, max)
	}
}

func Benchmark_EmptyLinifyString(b *testing.B) {
	var (
		lline = ""
		max   = 20
	)

	for b.Loop() {
		LinifyString(lline, max)
	}
}
