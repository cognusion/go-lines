// Package lines is a multipurpose line truncator, that ensures outputted lines do not exceed a specified length ever.
package lines

import "strings"

// RawLinifyString returns a string that has newlines inserted every max characters, irrespective of word boundaries
func RawLinifyString(s string, max int) string {
	// sanity
	if max < 1 {
		max = 128
	}

	var (
		cc        int
		newString string
	)

	for _, c := range s {
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

	for word := range strings.FieldsSeq(s) {
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
