package rfc9839

import (
	"unicode"
	"unicode/utf8"
)

// Exported functions

type Subset int

const (
	UnicodeScalar Subset = iota
	XmlChar
	UnicodeAssignable
)

func IsRuneUnicodeScalar(r rune) bool {
	return unicode.Is(unicodeScalars, r)
}
func IsRuneXmlChar(r rune) bool {
	return unicode.Is(xmlChars, r)
}
func IsRuneUnicodeAssignable(r rune) bool {
	return unicode.Is(unicodeAssignables, r)
}

func IsStringUnicodeScalars(s string) bool {
	return isStringInSubset(s, unicodeScalars)
}
func IsStringXmlChars(s string) bool {
	return isStringInSubset(s, xmlChars)
}
func IsStringUnicodeAssignables(s string) bool {
	return isStringInSubset(s, unicodeAssignables)
}

func IsUTF8UnicodeScalars(u []byte) bool {
	return isUTF8InSubset(u, unicodeScalars)
}
func IsUTF8XmlChars(u []byte) bool {
	return isUTF8InSubset(u, xmlChars)
}
func IsUTF8UnicodeAssignables(u []byte) bool {
	return isUTF8InSubset(u, unicodeAssignables)
}

// Internal

// These subset ranges are not sorted by order; the ranges most likely to
// contain runes being queried are moved to the front. "Most likely" is
// strictly based on Tim's intuition, there's no quantitative data behind it.

var unicodeScalars = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 0x0000, Hi: 0xD7FF, Stride: 1}, // most of the BMP
	},
	R32: []unicode.Range32{
		{Lo: 0xE000, Hi: 0x10FFFF, Stride: 1}, // mostly astral planes
	},
}

var xmlChars = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 0x0009, Hi: 0x0009, Stride: 1}, // Tab
		{Lo: 0x000A, Hi: 0x000A, Stride: 1}, // newline
		{Lo: 0x000D, Hi: 0x000D, Stride: 1}, // CR
		{Lo: 0x0020, Hi: 0xD7FF, Stride: 1}, // most of the BMP
		{Lo: 0xE000, Hi: 0xFFFD, Stride: 1}, // BMP after surrogates
	},
	R32: []unicode.Range32{
		{Lo: 0x10000, Hi: 0x10FFFF, Stride: 1}, // astral planes
	},
}

var unicodeAssignables = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 0x0009, Hi: 0x0009, Stride: 1}, // Tab
		{Lo: 0x000A, Hi: 0x000A, Stride: 1}, // newline
		{Lo: 0x000D, Hi: 0x000D, Stride: 1}, // CR
		{Lo: 0x0020, Hi: 0x007E, Stride: 1}, // ASCII
		{Lo: 0x00A0, Hi: 0xD7FF, Stride: 1}, // most of the BMP
		{Lo: 0xE000, Hi: 0xFDCF, Stride: 1}, // BMP after surrogates
		{Lo: 0xFDF0, Hi: 0xFFFD, Stride: 1}, // BMP after noncharacters block
	},
	R32: []unicode.Range32{
		{Lo: 0x10000, Hi: 0x1FFFD, Stride: 1}, // astral planes from here down
		{Lo: 0x20000, Hi: 0x2FFFD, Stride: 1},
		{Lo: 0x30000, Hi: 0x3FFFD, Stride: 1},
		{Lo: 0x40000, Hi: 0x4FFFD, Stride: 1},
		{Lo: 0x50000, Hi: 0x5FFFD, Stride: 1},
		{Lo: 0x60000, Hi: 0x6FFFD, Stride: 1},
		{Lo: 0x70000, Hi: 0x7FFFD, Stride: 1},
		{Lo: 0x80000, Hi: 0x8FFFD, Stride: 1},
		{Lo: 0x90000, Hi: 0x9FFFD, Stride: 1},
		{Lo: 0xA0000, Hi: 0xAFFFD, Stride: 1},
		{Lo: 0xB0000, Hi: 0xBFFFD, Stride: 1},
		{Lo: 0xC0000, Hi: 0xCFFFD, Stride: 1},
		{Lo: 0xD0000, Hi: 0xDFFFD, Stride: 1},
		{Lo: 0xE0000, Hi: 0xEFFFD, Stride: 1},
		{Lo: 0xF0000, Hi: 0xFFFFD, Stride: 1},
		{Lo: 0x100000, Hi: 0x10FFFD, Stride: 1},
	},
}

// …which raises issues of how you might optimize the current "simplest thing
// that could possibly work" implementation. I refuse to touch the code until
// I see evidence of noticeable slowdown.
// Bit arrays offer constant time but no matter how tight you pack, you're
// looking at ~400K of data structures.
// The worst case would probably be checking a large volume of Chinese text for
// UnicodeAssignable, because those code-points are at the third position in the
// subset just above, so there will be 4 wasted comparison operations for each
// code point.
// Having said that, the current implementation has excellent data locality and will
// probably be very CPU-cache-friendly.
// The craziest idea (and my favorite) is to build up a sample of which runePairs are
// being used the most and then from time to time shuffle the most-used to the front of the
// subset slice. You'd need an AtomicPointer update to be thread-safe.
// But like I said, I refuse to do anything without a sample corpus that is causing
// a measurable delay.

func isStringInSubset(s string, subset *unicode.RangeTable) bool {
	return isUTF8InSubset([]byte(s), subset)
}

func isUTF8InSubset(u []byte, subset *unicode.RangeTable) bool {
	index := 0
	for index < len(u) {
		r, width := utf8.DecodeRune(u[index:])
		if r == 0xFFFD && width == 1 {
			// this is how the utf8 pkg signals invalid UTF8 bytes, notably
			// including surrogate values
			return false
		}
		if !unicode.Is(subset, r) {
			return false
		}
		index += width
	}
	return true
}
