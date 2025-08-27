package rfc9839

import (
	"bufio"
	"compress/gzip"
	"os"
	"sync"
	"testing"
)

func BenchmarkCityLots(b *testing.B) {
	lines := getCityLotsLines(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lineIndex := i
		if i >= len(lines) {
			lineIndex = 0
		}
		if !Assignables.ValidUtf8(lines[lineIndex]) {
			panic("OUCH!")
		}
	}
}

const oneMeg = 1024 * 1024

var (
	cityLotsLock      sync.Mutex
	cityLotsLines     [][]byte
	cityLotsLineCount int
)

func getCityLotsLines(tb testing.TB) [][]byte {
	tb.Helper()

	cityLotsLock.Lock()
	defer cityLotsLock.Unlock()
	if cityLotsLines != nil {
		return cityLotsLines
	}
	file, err := os.Open("testdata/citylots.jlines.gz")
	if err != nil {
		tb.Error("Can't open citlots.jlines.gz: " + err.Error())
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	zr, err := gzip.NewReader(file)
	if err != nil {
		tb.Error("Can't open zip reader: " + err.Error())
	}

	scanner := bufio.NewScanner(zr)
	buf := make([]byte, oneMeg)
	scanner.Buffer(buf, oneMeg)
	for scanner.Scan() {
		cityLotsLineCount++
		cityLotsLines = append(cityLotsLines, []byte(scanner.Text()))
	}
	return cityLotsLines
}
