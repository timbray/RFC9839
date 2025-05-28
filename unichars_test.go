package unichars

import (
	"testing"
	"unicode/utf8"
)

var inverseScalars = []*runePair{
	{0xD800, 0xDFFF},
}

func TestEmpties(t *testing.T) {
	var emptyU1 = []byte{}
	var emptyU2 []byte = nil
	var emptyS = ""

	if !IsUTF8UnicodeScalars(emptyU1) {
		t.Error("zero-length []byte failrue")
	}
	if !IsUTF8XmlChars(emptyU1) {
		t.Error("zero-length []byte failrue")
	}
	if !IsUTF8UnicodeAssignables(emptyU1) {
		t.Error("zero-length []byte failrue")
	}

	if !IsUTF8UnicodeScalars(emptyU2) {
		t.Error("nil []byte failrue")
	}
	if !IsUTF8XmlChars(emptyU2) {
		t.Error("nil []byte failrue")
	}
	if !IsUTF8UnicodeAssignables(emptyU2) {
		t.Error("nil []byte failrue")
	}

	if !IsStringUnicodeScalars(emptyS) {
		t.Error("empty string failrue")
	}
	if !IsStringXmlChars(emptyS) {
		t.Error("empty string failrue")
	}
	if !IsStringUnicodeAssignables(emptyS) {
		t.Error("empty string failrue")
	}
}

func TestScalars(t *testing.T) {
	for _, pair := range unicodeScalars {
		for r := pair.lo; r <= pair.hi; r++ {
			if !IsRuneUnicodeScalar(r) {
				t.Errorf("%x should be unicode scalar", r)
			}
		}
	}
	for _, pair := range inverseScalars {
		for r := pair.lo; r <= pair.hi; r++ {
			if IsRuneUnicodeScalar(r) {
				t.Errorf("%x should not be unicode scalar", r)
			}
		}
	}
	if IsRuneUnicodeScalar(-1) {
		t.Error("-1 should not be scalar")
	}
	if IsRuneUnicodeScalar(0x10FFFF + 1) {
		t.Error("0x10FFFF+1 should not be scalar")
	}
	// a Go string can't contain a surrogate
	badUTF8 := []byte{0xED, 0xBA, 0xAD} // U+DEAD
	bad := []byte{'a'}
	bad = append(bad, badUTF8...)
	bad = append(bad, 'z')
	if IsUTF8UnicodeScalars(bad) {
		t.Error("accepted invalid UTF8")
	}
	if IsStringUnicodeScalars(string(bad)) {
		t.Error("accepted invalid UTF8")
	}
}

/*
	var xmlChars = []*runePair{
		{0x20, 0xD7FF},      // most of the BMP
		{0xA, 0xA},          // newline
		{0xE000, 0xFFFD},    // BMP after surrogates
		{0x9, 0x9},          // Tab
		{0xD, 0xD},          // CR
		{0x10000, 0x10FFFF}, // astral planes
	}
*/
var inverseXML = []*runePair{
	{0, 0x8},
	{0xB, 0xC},
	{0xE, 0x1F},
	{0xD800, 0xDFFF},
	{0xFFFE, 0xFFFF},
}

func TestXmlChars(t *testing.T) {
	for _, pair := range xmlChars {
		for r := pair.lo; r <= pair.hi; r++ {
			if !IsRuneXmlChar(r) {
				t.Errorf("%x should be XML", r)
			}
		}
	}
	for _, pair := range inverseXML {
		for r := pair.lo; r <= pair.hi; r++ {
			if IsRuneXmlChar(r) {
				t.Errorf("%x should not be XML", r)
			}
		}
	}
	if IsRuneXmlChar(-1) {
		t.Error("-1 should not be scalar")
	}
	if IsRuneXmlChar(0x10FFFF + 1) {
		t.Error("0x10FFFF+1 should not be scalar")
	}
	badUTF8 := []byte{0xED, 0xBA, 0xAD} // U+DEAD
	bad := []byte{'a'}
	bad = append(bad, badUTF8...)
	bad = append(bad, 'z')
	if IsUTF8XmlChars(bad) {
		t.Error("accepted invalid UTF8")
	}
	if IsStringXmlChars(string(bad)) {
		t.Error("accepted invalid UTF8")
	}
	goodS := []rune{}
	goodU := []byte{}
	for _, pair := range xmlChars {
		goodS = append(goodS, pair.lo, pair.hi)
		u := make([]byte, utf8.RuneLen(pair.lo))
		utf8.EncodeRune(u, pair.lo)
		goodU = append(goodU, u...)
		u = make([]byte, utf8.RuneLen(pair.hi))
		utf8.EncodeRune(u, pair.hi)
		goodU = append(goodU, u...)
	}
	if !IsStringXmlChars(string(goodS)) {
		t.Error("good string rejected")
	}
	if !IsUTF8XmlChars(goodU) {
		t.Error("good UTF8 rejected")
	}
	for _, pair := range inverseXML {
		if pair.lo == 0xD800 { // no surrogates
			continue
		}
		bad := []byte{'a'}
		u := make([]byte, utf8.RuneLen(pair.lo))
		utf8.EncodeRune(u, pair.lo)
		u = append(bad, u...)
		u = append(u, 'z')
		if IsUTF8XmlChars(u) {
			t.Errorf("accepted utf8 containing %x", pair.lo)
		}
		if IsStringXmlChars(string(u)) {
			t.Errorf("accepted utf8 containing %x", pair.lo)
		}
	}

}

/*
var unicodeAssignables = []*runePair{
	{0x9, 0x9},         // Tab
	{0xA, 0xA},         // newline
	{0xD, 0xD},         // CR
	{0x20, 0x7E},       // ASCII
	{0xA0, 0xD7FF},     // most of the BMP
	{0xE000, 0xFDCF},   // BMP after surrogates
	{0xFDF0, 0xFFFD},   // BMP after noncharacters block
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
*/

var inverseAssignables = []*runePair{
	{0, 0x8},
	{0xB, 0xC},
	{0xE, 0x1F},
	{0x7F, 0x9F},
	{0xD800, 0xDFFF},
	{0xFDD0, 0xFDEF},
	{0xFFFE, 0xFFFF},
	{0x1FFFE, 0x1FFFF},
	{0x2FFFE, 0x2FFFF},
	{0x3FFFE, 0x3FFFF},
	{0x4FFFE, 0x4FFFF},
	{0x5FFFE, 0x5FFFF},
	{0x6FFFE, 0x6FFFF},
	{0x7FFFE, 0x7FFFF},
	{0x8FFFE, 0x8FFFF},
	{0x9FFFE, 0x9FFFF},
	{0xAFFFE, 0xAFFFF},
	{0xBFFFE, 0xBFFFF},
	{0xCFFFE, 0xCFFFF},
	{0xDFFFE, 0xDFFFF},
	{0xEFFFE, 0xEFFFF},
	{0xFFFFE, 0xFFFFF},
	{0x10FFFE, 0xFFFFF},
}

func TestAssignables(t *testing.T) {
	for _, pair := range unicodeAssignables {
		for r := pair.lo; r <= pair.hi; r++ {
			if !IsRuneUnicodeAssignable(r) {
				t.Errorf("%x should be Assignable", r)
			}
		}
	}
	for _, pair := range inverseAssignables {
		for r := pair.lo; r <= pair.hi; r++ {
			if IsRuneUnicodeAssignable(r) {
				t.Errorf("%x should not be Assignable", r)
			}
		}
	}
	if IsRuneUnicodeAssignable(-1) {
		t.Error("-1 should not be scalar")
	}
	if IsRuneUnicodeAssignable(0x10FFFF + 1) {
		t.Error("0x10FFFF+1 should not be scalar")
	}
	badUTF8 := []byte{0xED, 0xBA, 0xAD} // U+DEAD
	bad := []byte{'a'}
	bad = append(bad, badUTF8...)
	bad = append(bad, 'z')
	if IsUTF8UnicodeAssignables(bad) {
		t.Error("accepted invalid UTF8")
	}
	if IsStringUnicodeAssignables(string(bad)) {
		t.Error("accepted invalid UTF8")
	}
	goodS := []rune{}
	goodU := []byte{}
	for _, pair := range unicodeAssignables {
		goodS = append(goodS, pair.lo, pair.hi)
		u := make([]byte, utf8.RuneLen(pair.lo))
		utf8.EncodeRune(u, pair.lo)
		goodU = append(goodU, u...)
		u = make([]byte, utf8.RuneLen(pair.hi))
		utf8.EncodeRune(u, pair.hi)
		goodU = append(goodU, u...)
	}
	if !IsStringUnicodeAssignables(string(goodS)) {
		t.Error("good string rejected")
	}
	if !IsUTF8UnicodeAssignables(goodU) {
		t.Error("good UTF8 rejected")
	}
	for _, pair := range inverseAssignables {
		if pair.lo == 0xD800 { // no surrogates
			continue
		}
		bad := []byte{'a'}
		u := make([]byte, utf8.RuneLen(pair.lo))
		utf8.EncodeRune(u, pair.lo)
		u = append(bad, u...)
		u = append(u, 'z')
		if IsUTF8UnicodeAssignables(u) {
			t.Errorf("accepted utf8 containing %x", pair.lo)
		}
		if IsStringUnicodeAssignables(string(u)) {
			t.Errorf("accepted utf8 containing %x", pair.lo)
		}
	}
}
