package rfc9839

import (
	"unicode/utf8"
)

type runePair struct {
	lo rune
	hi rune
}

type Subset struct{ pairs []runePair }

func (sub *Subset) ValidRune(r rune) bool {
	return subsetContains(sub, r)
}
func (sub *Subset) ValidStringPrevious(s string) bool {
	return isUTF8InSubset([]byte(s), sub)
}
func (sub *Subset) ValidString(s string) bool {
	return isStringInSubset(s, sub)
}

func (sub *Subset) ValidUtf8(u []byte) bool {
	return isUTF8InSubset(u, sub)
}

// implementation note: the Subset could contain, instead of []runePair, a
// unicode.RangeTable, then subsetContains could be replaced by unicode.Is(). We
// implemented this, but it had a >2x performance penalty.

// usage will look something like rfc9839.Scalars.ValidRune('r')

// note that these are not sorted by numeric order, but by in descending order of
// estimated traffic, as measured by Tim's guesswork. The idea is that you'd like
// to minimize the number of PairContains calls.

var Scalars = &Subset{
	pairs: []runePair{
		{0, 0xD7FF},        // most of the BMP
		{0xE000, 0x10FFFF}, // mostly astral planes
	},
}

var XmlChars = &Subset{
	pairs: []runePair{
		{0x20, 0xD7FF},      // most of the BMP
		{0xA, 0xA},          // newline
		{0xE000, 0xFFFD},    // BMP after surrogates
		{0x9, 0x9},          // Tab
		{0xD, 0xD},          // CR
		{0x10000, 0x10FFFF}, // astral planes
	},
}

var Assignables = &Subset{
	pairs: []runePair{
		{0x20, 0x7E},       // ASCII
		{0xA, 0xA},         // newline
		{0xA0, 0xD7FF},     // most of the BMP
		{0xE000, 0xFDCF},   // BMP after surrogates
		{0xFDF0, 0xFFFD},   // BMP after noncharacters block
		{0x9, 0x9},         // Tab
		{0xD, 0xD},         // CR
		{0x10000, 0x1FFFD}, // astral planes from here down
		{0x20000, 0x2FFFD},
		{0x30000, 0x3FFFD},
		{0x40000, 0x4FFFD},
		{0x50000, 0x5FFFD},
		{0x60000, 0x6FFFD},
		{0x70000, 0x7FFFD},
		{0x80000, 0x8FFFD},
		{0x90000, 0x9FFFD},
		{0xA0000, 0xAFFFD},
		{0xB0000, 0xBFFFD},
		{0xC0000, 0xCFFFD},
		{0xD0000, 0xDFFFD},
		{0xE0000, 0xEFFFD},
		{0xF0000, 0xFFFFD},
		{0x100000, 0x10FFFD},
	},
}

func pairContains(pair runePair, r rune) bool {
	return r >= pair.lo && r <= pair.hi
}

func subsetContains(sub *Subset, r rune) bool {
	for _, pair := range sub.pairs {
		if pairContains(pair, r) {
			return true
		}
	}
	return false
}

func isUTF8InSubset(u []byte, sub *Subset) bool {
	index := 0
	for index < len(u) {
		r, width := utf8.DecodeRune(u[index:])
		if r == 0xFFFD && width == 1 {
			// this is how the utf8 pkg signals invalid UTF8 bytes, notably
			// including surrogate values
			return false
		}
		if !subsetContains(sub, r) {
			return false
		}
		index += width
	}
	return true
}

func isStringInSubset(s string, sub *Subset) bool {
	index := 0
	for index < len(s) {
		r, width := utf8.DecodeRuneInString(s[index:])
		if r == 0xFFFD && width == 1 {
			// this is how the utf8 pkg signals invalid UTF8 bytes, notably
			// including surrogate values
			return false
		}
		if !subsetContains(sub, r) {
			return false
		}
		index += width
	}
	return true
}
