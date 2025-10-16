// Package lines is a multipurpose ASCII line truncator, that ensures outputted lines do not exceed a specified length ever.
// I was unimpressed with the packages I found for quality, speed, or both.
//
// There is a major blindspot here with regards to multibyte characters. That's not a thing I am working toward due to its inherent
// complexity and the purpose of this mod. PRs are welcome, but better packages probably exist
// elsewhere.
package lines

import (
	"io"
	"strings"
)

// RawLinifyString returns a string that has newlines inserted every max characters, irrespective of word boundaries.
func RawLinifyString(s string, max int) string {
	// sanity
	if max < 1 {
		max = 128
	}

	var (
		cc        int
		newString string
	)

	var c rune
	for _, c = range s {
		if cc >= max {
			newString += "\n" + string(c)
			cc = 1 // not zero as we are appending a rune there
		} else {
			newString += string(c)
			cc++
		}
	}
	return newString
}

// LinifyString returns a string that has newlines inserted between word boundaries, at most every max characters.
func LinifyString(s string, max int) string {
	// sanity
	if max < 1 {
		max = 128
	}

	var (
		newString string
		llen      int
	)

	// :easybutton: using Fields. We should do our own
	// tokenizing if this gets traction.
	var word string
	for word = range strings.FieldsSeq(s) {
		if len(word) > max {
			// are you KIDDING ME RIGHT NOW?!
			if llen > 0 {
				newString += "\n"
			}
			llen = 0
			newString += RawLinifyString(word, max) + "\n" // UGLYAF
		} else if llen > 0 && llen+1+len(word) > max {
			// there's a word already, and it blows out max
			llen = len(word)
			newString += "\n" + word
		} else if llen == 0 {
			// first word in the line, just append the word
			llen = len(word)
			newString += word
		} else {
			// llen > 0, so we append space and the word
			llen += 1 + len(word)
			newString += " " + word
		}
	}
	return newString
}

// LinifyStream consumes a string chan and pushes linified results to the specified io.StringWriter.
// An error is returned IFF the io.StringWriter returns an error.
// This is only meaningfully efficient for arbitrarily massive sets of strings. Unless you are
// linifying 'The Tommyknockers' or 'War and Peace', I doubt this is what you're looking for.
func LinifyStream(stream <-chan string, out io.StringWriter, max int) error {
	return LinifyStreamSeparator(stream, out, max, "")
}

// LinifyStreamSeparator consumes a string chan and pushes linified results to the specified io.StringWriter.
// The separator may specify what is used to separate words.
// An error is returned IFF the io.StringWriter returns an error.
// This is only meaningfully efficient for arbitrarily massive sets of strings. Unless you are
// linifying 'The Tommyknockers' or 'War and Peace', I doubt this is what you're looking for.
func LinifyStreamSeparator(stream <-chan string, out io.StringWriter, max int, separator string) error {
	// each string across stream is a word.
	// case and punctuation preserved.
	// We add spaces and newlines only
	var (
		llen int
	)

	var word string
	for word = range stream {
		if len(word) > max {
			// are you KIDDING ME RIGHT NOW?!
			if llen > 0 {
				if _, err := out.WriteString(separator + "\n"); err != nil {
					return err
				}
			}
			llen = 0
			if _, err := out.WriteString(RawLinifyString(word, max-len(separator)) + separator + "\n"); err != nil {
				return err
			}
		} else if llen > 0 && llen+len(separator)+1+len(word) > max {
			// there's a word already, and it blows out max
			llen = len(word)
			if _, err := out.WriteString(separator + "\n" + word); err != nil {
				return err
			}
		} else if llen == 0 {
			// first word in the line, just append the word
			llen = len(word)
			if _, err := out.WriteString(word); err != nil {
				return err
			}
		} else {
			// llen > 0, so we append a separator and the word
			llen += len(separator) + 1 + len(word)
			if _, err := out.WriteString(separator + " " + word); err != nil {
				return err
			}
		}
	}
	return nil
}
