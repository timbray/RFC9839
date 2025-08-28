# RFC9839

[![Tests](https://github.com/timbray/RFC9839/actions/workflows/go-unit-tests.yaml/badge.svg)](https://github.com/timbray/RFC9839/actions/workflows/go-unit-tests.yaml)
[![codecov](https://codecov.io/gh/timbray/RFC9839/graph/badge.svg?token=6V5I17FTIM)](https://codecov.io/gh/timbray/RFC9839)
[![0 dependencies!](https://0dependencies.dev/0dependencies.svg)](https://0dependencies.dev)

Go-language library to check for problematic Unicode code points.

This is based on the Unicode code-point subsets specified in [RFC9839](https://www.rfc-editor.org/rfc/rfc9839.html).

The package defines a `Subset` type and exports three instances, named `Scalars`,
`XmlChars`, and `Assignables`. It exports three functions:

```go
func (sub *Subset) ValidRune(r rune) bool
func (sub *Subset) ValidString(s string) bool
func (sub *Subset) ValidUtf8(u []byte) bool
```

A typical call might look like:

```go
if !rfc9839.Assignables.ValidRune(r) {
	return t.Error("invalid rune")
}
```