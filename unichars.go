package rfc9839

import (
	"unicode"
	"unicode/utf8"
)

// exported functions

type Subset struct {
	table unicode.RangeTable
}

func (sub *Subset) ValidRune(r rune) bool {
	return unicode.Is(&sub.table, r)
}
func (sub *Subset) ValidString(s string) bool {
	return isUTF8InSubset([]byte(s), &sub.table)
}
func (sub *Subset) ValidUtf8(u []byte) bool {
	return isUTF8InSubset(u, &sub.table)
}

// exported variables, a call looks like rfc9839.Scalars.ValidRune(r)

var Scalars = &Subset{
	unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0000, Hi: 0xD7FF, Stride: 1}, // most of the BMP
		},
		R32: []unicode.Range32{
			{Lo: 0xE000, Hi: 0x10FFFF, Stride: 1}, // mostly astral planes
		},
	},
}

var XmlChars = &Subset{
	unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0009, Hi: 0x000A, Stride: 1}, // tab & newline
			{Lo: 0x000D, Hi: 0x000D, Stride: 1}, // CR
			{Lo: 0x0020, Hi: 0xD7FF, Stride: 1}, // most of the BMP
			{Lo: 0xE000, Hi: 0xFFFD, Stride: 1}, // BMP after surrogates
		},
		R32: []unicode.Range32{
			{Lo: 0x10000, Hi: 0x10FFFF, Stride: 1}, // astral planes
		},
	},
}

var Assignables = &Subset{
	unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0009, Hi: 0x000A, Stride: 1}, // tab & newline
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
	},
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
