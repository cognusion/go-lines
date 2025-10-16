package lines

import (
	"os"
	"strings"
	"testing"

	"github.com/cognusion/go-recyclable"
	"github.com/fortytw2/leaktest"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_RawLinifyString(t *testing.T) {
	defer leaktest.Check(t)()

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
	defer leaktest.Check(t)()

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

func ExampleLinifyStream() {

	// You are going to have a scanner with an unknown,
	// never-ending buffer of words that need to be
	// assembled with spaces between, and newlines on or
	// before some max line length.
	fields := strings.FieldsSeq("This line is not that long. But imagine it is much longer.")

	// create a channel for those words to pipe over
	wordChan := make(chan string)

	// create a possibly never-ending stream to send words to the wordChan
	go func() {
		defer close(wordChan) // super important, if we ever want to end
		for word := range fields {
			wordChan <- word
		}
	}()

	// Linify the stream from wordChan, write to os.StdOut, each line max 20 characters, separate by space.
	err := LinifyStream(wordChan, os.Stdout, 20)
	if err != nil {
		// for real?!
		panic(err)
	}
	// Output: This line is not
	//that long. But
	//imagine it is much
	//longer.

}

func Test_LinifyStreamSeparator(t *testing.T) {
	defer leaktest.Check(t)()

	var (
		max    = 20
		lword  = "areyoukiddingmerightnow"
		linlw  = "areyoukiddingmeright\nnow\n"
		leline = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	)

	fields := strings.Fields(leline)
	schan := make(chan string)
	buff := &recyclable.Buffer{}

	// start the stream
	go func() {
		defer close(schan)
		for _, f := range fields {
			schan <- f
		}
	}()
	Convey("When a long line is streamed Linified, the results are expected.", t, FailureContinues, func() {
		err := LinifyStreamSeparator(schan, buff, max, "")
		So(err, ShouldBeNil)
		SoMsg("Incorrect number of newlines added", buff.Len(), ShouldEqual, 445)

		Convey("When a stupid word is streamed, larger than max, the results are expected", FailureContinues, func() {
			buff.Reset(make([]byte, 0))
			schan = make(chan string, 1)
			schan <- lword // queue up the word
			close(schan)
			err := LinifyStreamSeparator(schan, buff, max, "")
			So(err, ShouldBeNil)
			So(buff.String(), ShouldEqual, linlw)

		})

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
