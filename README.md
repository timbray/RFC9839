# RFC9839
(??)

Rust library to check for problematic Unicode code points.

This is a Rust port of [Tim Bray's original Go implementation](https://github.com/timbray/RFC9839) of RFC 9839 Unicode character validation.

## Overview

This library validates Unicode code-point subsets as specified in RFC 9839.

There are three named subsets:
- **Unicode Scalars** - All Unicode code points except surrogates (U+D800 to U+DFFF)
- **XML Characters** - Valid XML 1.0 characters
- **Unicode Assignables** - Unicode characters excluding noncharacters

## Usage

```rust
use rfc9839::{is_char_unicode_scalar, is_string_xml_chars, is_utf8_unicode_assignables};

// Check individual characters
assert!(is_char_unicode_scalar('A'));

// Check strings
assert!(is_string_xml_chars("Valid XML text\n"));

// Check UTF-8 byte slices
assert!(is_utf8_unicode_assignables(b"Hello, world!"));
```

## API Compatibility

This port maintains full API compatibility with the Go implementation, providing the same function names and behaviour. However, the Rust implementation leverages the language's built-in guarantees where possible:

- **Unicode scalar validation**: Since Rust's `char` type can only represent valid Unicode scalar values, the Unicode scalar validation functions use Rust's built-in validation behind the scenes. For example, `is_rune_unicode_scalar` simply checks whether `char::from_u32()` succeeds.

- **String validation**: Rust strings are guaranteed to be valid UTF-8 without surrogates, so `is_string_unicode_scalars` always returns `true` for any valid `&str`.

- **UTF-8 validation**: The `is_utf8_unicode_scalars` function leverages Rust's `str::from_utf8`, which already rejects surrogates and invalid UTF-8.

The XML character and Unicode assignable validation functions perform actual subset checking, as these are more restrictive than what Rust's type system guarantees.

The exported function names are descriptive of what they do.
