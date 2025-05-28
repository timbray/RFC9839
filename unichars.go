package unichars

import "unicode/utf8"

type runePair struct {
	lo rune
	hi rune
}

// These subset ranges are not sorted by order; the ranges most likely to
// contain runes being queried are moved to the front. "most likely" is
// strictly based on Tim's intuition, there's no quantitative data behind it.

var unicodeScalars = []*runePair{
	{0, 0xD7FF},        // most of the BMP
	{0xE000, 0x10FFFF}, // mostly astral planes
}

var xmlChars = []*runePair{
	{0x20, 0xD7FF},      // most of the BMP
	{0xA, 0xA},          // newline
	{0xE000, 0xFFFD},    // BMP after surrogates
	{0x9, 0x9},          // Tab
	{0xD, 0xD},          // CR
	{0x10000, 0x10FFFF}, // astral planes
}

var unicodeAssignables = []*runePair{
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
}

func pairContains(pair *runePair, r rune) bool {
	return r >= pair.lo && r <= pair.hi
}
func rangeContains(subset []*runePair, r rune) bool {
	for _, pair := range subset {
		if pairContains(pair, r) {
			return true
		}
	}
	return false
}

func IsRuneUnicodeScalar(r rune) bool {
	return rangeContains(unicodeScalars, r)
}
func IsRuneXmlChar(r rune) bool {
	return rangeContains(xmlChars, r)
}
func IsRuneUnicodeAssignable(r rune) bool {
	return rangeContains(unicodeAssignables, r)
}
func isStringInSubset(s string, subset []*runePair) bool {
	return isUTF8InSubset([]byte(s), subset)
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
func isUTF8InSubset(u []byte, subset []*runePair) bool {
	index := 0
	for index < len(u) {
		r, width := utf8.DecodeRune(u[index:])
		if r == 0xFFFD && width == 1 {
			return false
		}
		if !rangeContains(subset, r) {
			return false
		}
		index += width
	}
	return true
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
