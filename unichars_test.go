package rfc9839

import (
	"testing"
	"unicode/utf8"
)

var inverseScalars = []rune{
	0xD800, 0xDFFF, // surrogate pairs
}

func TestEmpties(t *testing.T) {
	var emptyU1 []byte
	var emptyU2 []byte = nil
	var emptyS = ""

	if !Scalars.ValidUtf8(emptyU1) {
		t.Error("zero-length []byte failure")
	}
	if !XmlChars.ValidUtf8(emptyU1) {
		t.Error("zero-length []byte failure")
	}
	if !Assignables.ValidUtf8(emptyU1) {
		t.Error("zero-length []byte failure")
	}

	if !Scalars.ValidUtf8(emptyU2) {
		t.Error("nil []byte failure")
	}
	if !XmlChars.ValidUtf8(emptyU2) {
		t.Error("nil []byte failure")
	}
	if !Assignables.ValidUtf8(emptyU2) {
		t.Error("nil []byte failure")
	}

	if !Scalars.ValidString(emptyS) {
		t.Error("empty string failure")
	}
	if !XmlChars.ValidString(emptyS) {
		t.Error("empty string failure")
	}
	if !Assignables.ValidString(emptyS) {
		t.Error("empty string failure")
	}
}

func TestScalars(t *testing.T) {
	// Test that all runes in our unicodeScalars table are accepted
	for _, r16 := range Scalars.table.R16 {
		for r := r16.Lo; r <= r16.Hi; r += r16.Stride {
			if !Scalars.ValidRune(rune(r)) {
				t.Errorf("%x should be unicode scalar", r)
			}
		}
	}
	for _, r32 := range Scalars.table.R32 {
		for r := r32.Lo; r <= r32.Hi; r += r32.Stride {
			if !Scalars.ValidRune(rune(r)) {
				t.Errorf("%x should be unicode scalar", r)
			}
		}
	}

	// Test that surrogate pairs are rejected
	for _, r := range inverseScalars {
		if Scalars.ValidRune(r) {
			t.Errorf("%x should not be unicode scalar", r)
		}
	}

	if Scalars.ValidRune(-1) {
		t.Error("-1 should not be scalar")
	}
	if Scalars.ValidRune(0x10FFFF + 1) {
		t.Error("0x10FFFF+1 should not be scalar")
	}

	// a Go string can't contain a surrogate
	badUTF8 := []byte{0xED, 0xBA, 0xAD} // U+DEAD
	bad := []byte{'a'}
	bad = append(bad, badUTF8...)
	bad = append(bad, 'z')
	if Scalars.ValidUtf8(bad) {
		t.Error("accepted invalid UTF8")
	}
	if Scalars.ValidString(string(bad)) {
		t.Error("accepted invalid UTF8")
	}
}

var inverseXML = []rune{
	0x0000, 0x0008, // control characters
	0x000B, 0x000C, // vertical tab, form feed
	0x000E, 0x001F, // control characters
	0xD800, 0xDFFF, // surrogate pairs
	0xFFFE, 0xFFFF, // noncharacters
}

func TestXmlChars(t *testing.T) {
	// Test that all runes in our xmlChars table are accepted
	for _, r16 := range XmlChars.table.R16 {
		for r := r16.Lo; r <= r16.Hi; r += r16.Stride {
			if !XmlChars.ValidRune(rune(r)) {
				t.Errorf("%x should be XML", r)
			}
		}
	}
	for _, r32 := range XmlChars.table.R32 {
		for r := r32.Lo; r <= r32.Hi; r += r32.Stride {
			if !XmlChars.ValidRune(rune(r)) {
				t.Errorf("%x should be XML", r)
			}
		}
	}

	// Test that inverse ranges are rejected
	for _, r := range inverseXML {
		if XmlChars.ValidRune(rune(r)) {
			t.Errorf("%x should not be XML", r)
		}
	}

	if XmlChars.ValidRune(-1) {
		t.Error("-1 should not be scalar")
	}
	if XmlChars.ValidRune(0x10FFFF + 1) {
		t.Error("0x10FFFF+1 should not be scalar")
	}

	badUTF8 := []byte{0xED, 0xBA, 0xAD} // U+DEAD
	bad := []byte{'a'}
	bad = append(bad, badUTF8...)
	bad = append(bad, 'z')
	if XmlChars.ValidUtf8(bad) {
		t.Error("accepted invalid UTF8")
	}
	if XmlChars.ValidString(string(bad)) {
		t.Error("accepted invalid UTF8")
	}

	// Test good strings
	var goodS []rune
	var goodU []byte
	for _, r16 := range XmlChars.table.R16 {
		goodS = append(goodS, rune(r16.Lo), rune(r16.Hi))
		loLen := utf8.RuneLen(rune(r16.Lo))
		if loLen > 0 {
			u := make([]byte, loLen)
			utf8.EncodeRune(u, rune(r16.Lo))
			goodU = append(goodU, u...)
		}
		hiLen := utf8.RuneLen(rune(r16.Hi))
		if hiLen > 0 {
			u := make([]byte, hiLen)
			utf8.EncodeRune(u, rune(r16.Hi))
			goodU = append(goodU, u...)
		}
	}
	for _, r32 := range XmlChars.table.R32 {
		goodS = append(goodS, rune(r32.Lo), rune(r32.Hi))
		loLen := utf8.RuneLen(rune(r32.Lo))
		if loLen > 0 {
			u := make([]byte, loLen)
			utf8.EncodeRune(u, rune(r32.Lo))
			goodU = append(goodU, u...)
		}
		hiLen := utf8.RuneLen(rune(r32.Hi))
		if hiLen > 0 {
			u := make([]byte, hiLen)
			utf8.EncodeRune(u, rune(r32.Hi))
			goodU = append(goodU, u...)
		}
	}

	if !XmlChars.ValidString(string(goodS)) {
		t.Error("good string rejected")
	}
	if !XmlChars.ValidUtf8(goodU) {
		t.Error("good UTF8 rejected")
	}

	// Test bad strings
	for _, r := range inverseXML {
		if r == 0xD800 { // no surrogates
			continue
		}
		bad := []byte{'a'}
		runeLen := utf8.RuneLen(r)
		if runeLen > 0 {
			u := make([]byte, runeLen)
			utf8.EncodeRune(u, r)
			u = append(bad, u...)
			u = append(u, 'z')
			if XmlChars.ValidUtf8(u) {
				t.Errorf("accepted utf8 containing %x", r)
			}
			if XmlChars.ValidString(string(u)) {
				t.Errorf("accepted utf8 containing %x", r)
			}
		}
	}
}

var inverseAssignables = []rune{
	0x0000, 0x0008, // control characters
	0x000B, 0x000C, // vertical tab, form feed
	0x000E, 0x001F, // control characters
	0x007F, 0x009F, // control characters
	0xD800, 0xDFFF, // surrogate pairs
	0xFDD0, 0xFDEF, // noncharacters
	0xFFFE, 0xFFFF, // noncharacters
	0x1FFFE, 0x1FFFF, // noncharacters
	0x2FFFE, 0x2FFFF, // noncharacters
	0x3FFFE, 0x3FFFF, // noncharacters
	0x4FFFE, 0x4FFFF, // noncharacters
	0x5FFFE, 0x5FFFF, // noncharacters
	0x6FFFE, 0x6FFFF, // noncharacters
	0x7FFFE, 0x7FFFF, // noncharacters
	0x8FFFE, 0x8FFFF, // noncharacters
	0x9FFFE, 0x9FFFF, // noncharacters
	0xAFFFE, 0xAFFFF, // noncharacters
	0xBFFFE, 0xBFFFF, // noncharacters
	0xCFFFE, 0xCFFFF, // noncharacters
	0xDFFFE, 0xDFFFF, // noncharacters
	0xEFFFE, 0xEFFFF, // noncharacters
	0xFFFFE, 0xFFFFF, // noncharacters
	0x10FFFE, 0x10FFFF, // noncharacters
}

func TestAssignables(t *testing.T) {
	// Test that all runes in our unicodeAssignables table are accepted
	for _, r16 := range Assignables.table.R16 {
		for r := r16.Lo; r <= r16.Hi; r += r16.Stride {
			if !Assignables.ValidRune(rune(r)) {
				t.Errorf("%x should be Assignable", r)
			}
		}
	}
	for _, r32 := range Assignables.table.R32 {
		for r := r32.Lo; r <= r32.Hi; r += r32.Stride {
			if !Assignables.ValidRune(rune(r)) {
				t.Errorf("%x should be Assignable", r)
			}
		}
	}

	// Test that inverse ranges are rejected
	for _, r := range inverseAssignables {
		if Assignables.ValidRune(r) {
			t.Errorf("%x should not be Assignable", r)
		}
	}

	if Assignables.ValidRune(-1) {
		t.Error("-1 should not be scalar")
	}
	if Assignables.ValidRune(0x10FFFF + 1) {
		t.Error("0x10FFFF+1 should not be scalar")
	}

	badUTF8 := []byte{0xED, 0xBA, 0xAD} // U+DEAD
	bad := []byte{'a'}
	bad = append(bad, badUTF8...)
	bad = append(bad, 'z')
	if Assignables.ValidUtf8(bad) {
		t.Error("accepted invalid UTF8")
	}
	if Assignables.ValidString(string(bad)) {
		t.Error("accepted invalid UTF8")
	}

	// Test good strings
	var goodS []rune
	var goodU []byte
	for _, r16 := range Assignables.table.R16 {
		goodS = append(goodS, rune(r16.Lo), rune(r16.Hi))
		loLen := utf8.RuneLen(rune(r16.Lo))
		if loLen > 0 {
			u := make([]byte, loLen)
			utf8.EncodeRune(u, rune(r16.Lo))
			goodU = append(goodU, u...)
		}
		hiLen := utf8.RuneLen(rune(r16.Hi))
		if hiLen > 0 {
			u := make([]byte, hiLen)
			utf8.EncodeRune(u, rune(r16.Hi))
			goodU = append(goodU, u...)
		}
	}
	for _, r32 := range Assignables.table.R32 {
		goodS = append(goodS, rune(r32.Lo), rune(r32.Hi))
		loLen := utf8.RuneLen(rune(r32.Lo))
		if loLen > 0 {
			u := make([]byte, loLen)
			utf8.EncodeRune(u, rune(r32.Lo))
			goodU = append(goodU, u...)
		}
		hiLen := utf8.RuneLen(rune(r32.Hi))
		if hiLen > 0 {
			u := make([]byte, hiLen)
			utf8.EncodeRune(u, rune(r32.Hi))
			goodU = append(goodU, u...)
		}
	}

	if !Assignables.ValidString(string(goodS)) {
		t.Error("good string rejected")
	}
	if !Assignables.ValidUtf8(goodU) {
		t.Error("good UTF8 rejected")
	}

	// Test bad strings
	for _, r := range inverseAssignables {
		if r == 0xD800 { // no surrogates
			continue
		}
		bad := []byte{'a'}
		runeLen := utf8.RuneLen(r)
		if runeLen > 0 {
			u := make([]byte, runeLen)
			utf8.EncodeRune(u, r)
			u = append(bad, u...)
			u = append(u, 'z')
			if Assignables.ValidUtf8(u) {
				t.Errorf("accepted utf8 containing %x", r)
			}
			if Assignables.ValidString(string(u)) {
				t.Errorf("accepted utf8 containing %x", r)
			}
		}
	}
}
