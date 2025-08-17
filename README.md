# RFC9839

[![Tests](https://github.com/timbray/quamina/actions/workflows/go-unit-tests.yaml/badge.svg)](https://github.com/timbray/quamina/actions/workflows/go-unit-tests.yaml)
[![codecov](https://codecov.io/gh/timbray/quamina/branch/main/graph/badge.svg?token=TC7MW723JO)](https://codecov.io/gh/timbray/quamina)
[![0 dependencies!](https://0dependencies.dev/0dependencies.svg)](https://0dependencies.dev)

Go-language library to check for problematic Unicode code points.

This is based on the Unicode code-point subsets specified in RFC9839.

There are three named subsets:
- Unicode Scalars
- XML Characters
- Unicode Assignables

The exported function names are descriptive of what they do.


