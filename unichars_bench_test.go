package rfc9839

import (
	"os"
	"testing"
)

func BenchmarkValidUtf8(b *testing.B) {
	file, err := os.ReadFile("testdata/sample.txt")
	if err != nil {
		b.Error(err)
	}

	bytes := len(file)
	b.SetBytes(int64(bytes))
	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Assignables.ValidUtf8(file)
		Scalars.ValidUtf8(file)
		XmlChars.ValidUtf8(file)
	}
}

func BenchmarkValidStringPrevious(b *testing.B) {
	file, err := os.ReadFile("testdata/sample.txt")
	if err != nil {
		b.Error(err)
	}

	bytes := len(file)
	b.SetBytes(int64(bytes))
	s := string(file)
	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Assignables.ValidStringPrevious(s)
		Scalars.ValidStringPrevious(s)
		XmlChars.ValidStringPrevious(s)
	}
}

func BenchmarkValidString(b *testing.B) {
	file, err := os.ReadFile("testdata/sample.txt")
	if err != nil {
		b.Error(err)
	}

	bytes := len(file)
	b.SetBytes(int64(bytes))
	s := string(file)
	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Assignables.ValidString(s)
		Scalars.ValidString(s)
		XmlChars.ValidString(s)
	}
}
