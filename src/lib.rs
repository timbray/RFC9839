//! RFC 9839 Unicode subset validation
//!
//! This is a Rust port of Tim Bray's original Go implementation of RFC 9839
//! Unicode character validation.
//!
//! Original repository: <https://github.com/timbray/RFC9839>
//!
//! This library provides functions to validate whether Unicode characters, strings,
//! or UTF-8 byte sequences contain only characters from specific RFC 9839-defined
//! subsets:
//! - Unicode Scalar Values (excludes surrogates)
//! - XML 1.0 Characters
//! - Unicode Assignable Characters (excludes noncharacters)
//!
//! # API Compatibility
//!
//! This port maintains full API compatibility with the Go implementation, providing
//! the same function names and behaviour. However, the Rust implementation leverages
//! the language's built-in guarantees where possible:
//!
//! - For Unicode scalar validation, we use Rust's built-in validation since `char`
//!   can only represent valid Unicode scalar values
//! - The [`is_string_unicode_scalars`] function always returns `true` for valid Rust
//!   strings, as they're guaranteed to be valid UTF-8 without surrogates
//! - The [`is_utf8_unicode_scalars`] function simply checks if bytes are valid UTF-8
//!   using [`str::from_utf8`], which already rejects surrogates
//!
//! The XML character and Unicode assignable validation functions perform actual
//! subset checking, as these are more restrictive than what Rust's type system
//! guarantees.
//!
//! # Implementation Notes
//!
//! 1. **Simplified Unicode scalar validation**: Since Rust's `char` type can only
//!    represent valid Unicode scalar values (automatically excluding surrogates),
//!    Unicode scalar validation uses Rust's built-in checks. The [`is_rune_unicode_scalar`]
//!    function simply uses [`char::from_u32()`] and `is_some()`.
//!
//! 2. **Optimisations**: Range arrays for XML and Assignable subsets are defined
//!    as `const` for compile-time optimisation. The Unicode Scalars array was
//!    removed entirely as Rust's type system handles this validation.

/// Unicode subset types as defined in RFC 9839
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Subset {
    UnicodeScalar,
    XmlChar,
    UnicodeAssignable,
}

/// A Unicode code point range (inclusive)
#[derive(Debug, Clone, Copy)]
struct RunePair {
    lo: u32,
    hi: u32,
}

impl RunePair {
    const fn new(lo: u32, hi: u32) -> Self {
        RunePair { lo, hi }
    }

    #[inline]
    fn contains(&self, ch: char) -> bool {
        let code = ch as u32;
        code >= self.lo && code <= self.hi
    }
}

// These subset ranges are not sorted by order; the ranges most likely to
// contain code points being queried are moved to the front.

// Note: UNICODE_SCALARS array removed as Rust's char type already enforces this

const XML_CHARS: &[RunePair] = &[
    RunePair::new(0x20, 0xD7FF),      // most of the BMP
    RunePair::new(0xA, 0xA),          // newline
    RunePair::new(0xE000, 0xFFFD),    // BMP after surrogates
    RunePair::new(0x9, 0x9),          // Tab
    RunePair::new(0xD, 0xD),          // CR
    RunePair::new(0x10000, 0x10FFFF), // astral planes
];

const UNICODE_ASSIGNABLES: &[RunePair] = &[
    RunePair::new(0x20, 0x7E),       // ASCII
    RunePair::new(0xA, 0xA),         // newline
    RunePair::new(0xA0, 0xD7FF),     // most of the BMP
    RunePair::new(0xE000, 0xFDCF),   // BMP after surrogates
    RunePair::new(0xFDF0, 0xFFFD),   // BMP after noncharacters block
    RunePair::new(0x9, 0x9),         // Tab
    RunePair::new(0xD, 0xD),         // CR
    RunePair::new(0x10000, 0x1FFFD), // astral planes from here down
    RunePair::new(0x20000, 0x2FFFD),
    RunePair::new(0x30000, 0x3FFFD),
    RunePair::new(0x40000, 0x4FFFD),
    RunePair::new(0x50000, 0x5FFFD),
    RunePair::new(0x60000, 0x6FFFD),
    RunePair::new(0x70000, 0x7FFFD),
    RunePair::new(0x80000, 0x8FFFD),
    RunePair::new(0x90000, 0x9FFFD),
    RunePair::new(0xA0000, 0xAFFFD),
    RunePair::new(0xB0000, 0xBFFFD),
    RunePair::new(0xC0000, 0xCFFFD),
    RunePair::new(0xD0000, 0xDFFFD),
    RunePair::new(0xE0000, 0xEFFFD),
    RunePair::new(0xF0000, 0xFFFFD),
    RunePair::new(0x100000, 0x10FFFD),
];

#[inline]
fn subset_contains(subset: &[RunePair], ch: char) -> bool {
    subset.iter().any(|pair| pair.contains(ch))
}

fn subset_contains_u32(subset: &[RunePair], code: u32) -> bool {
    subset.iter().any(|pair| code >= pair.lo && code <= pair.hi)
}

// Character validation functions
// Note: is_char_unicode_scalar removed as Rust chars are always Unicode scalars

/// Check if a character is a valid XML character
///
/// Valid XML 1.0 characters exclude most control characters below U+0020,
/// except for tab (U+0009), line feed (U+000A), and carriage return (U+000D).
///
/// # Examples
///
/// ```
/// use rfc9839::is_char_xml_char;
///
/// assert!(is_char_xml_char('A'));
/// assert!(is_char_xml_char('\n'));  // Line feed is valid
/// assert!(is_char_xml_char('\t'));  // Tab is valid
/// assert!(is_char_xml_char('\r'));  // Carriage return is valid
/// assert!(!is_char_xml_char('\u{0008}'));  // Backspace is invalid
/// assert!(!is_char_xml_char('\u{0000}'));  // Null is invalid
/// ```
pub fn is_char_xml_char(ch: char) -> bool {
    subset_contains(XML_CHARS, ch)
}

/// Check if a character is a Unicode assignable character
///
/// Unicode assignable characters exclude noncharacters like U+FFFE, U+FFFF,
/// and the range U+FDD0 to U+FDEF, as well as the last two code points of
/// each plane (e.g., U+1FFFE, U+1FFFF).
///
/// # Examples
///
/// ```
/// use rfc9839::is_char_unicode_assignable;
///
/// assert!(is_char_unicode_assignable('A'));
/// assert!(is_char_unicode_assignable('中'));
/// assert!(is_char_unicode_assignable('\n'));
/// assert!(!is_char_unicode_assignable('\u{FFFE}'));  // Noncharacter
/// assert!(!is_char_unicode_assignable('\u{FFFF}'));  // Noncharacter
/// assert!(!is_char_unicode_assignable('\u{1FFFE}')); // Plane 1 noncharacter
/// ```
pub fn is_char_unicode_assignable(ch: char) -> bool {
    subset_contains(UNICODE_ASSIGNABLES, ch)
}

// Rune (u32) validation functions to match Go's rune type

/// Check if a u32 code point is a Unicode scalar value
///
/// This function matches the Go implementation's rune-based API.
/// Returns true if the u32 can be converted to a valid Rust char
/// (i.e., it's a valid Unicode scalar value).
///
/// # Examples
///
/// ```
/// use rfc9839::is_rune_unicode_scalar;
///
/// assert!(is_rune_unicode_scalar(0x41));  // 'A'
/// assert!(is_rune_unicode_scalar(0x1F980));  // 🦀
/// assert!(!is_rune_unicode_scalar(0xD800));  // Surrogate
/// assert!(!is_rune_unicode_scalar(0xDFFF));  // Surrogate
/// assert!(!is_rune_unicode_scalar(0x110000)); // Beyond Unicode
/// ```
pub fn is_rune_unicode_scalar(r: u32) -> bool {
    char::from_u32(r).is_some()
}

/// Check if a u32 code point is a valid XML character
///
/// # Examples
///
/// ```
/// use rfc9839::is_rune_xml_char;
///
/// assert!(is_rune_xml_char(0x41));  // 'A'
/// assert!(is_rune_xml_char(0x09));  // Tab
/// assert!(is_rune_xml_char(0x0A));  // Line feed
/// assert!(!is_rune_xml_char(0x00));  // Null
/// assert!(!is_rune_xml_char(0x08));  // Backspace
/// ```
pub fn is_rune_xml_char(r: u32) -> bool {
    subset_contains_u32(XML_CHARS, r)
}

/// Check if a u32 code point is a Unicode assignable character
///
/// # Examples
///
/// ```
/// use rfc9839::is_rune_unicode_assignable;
///
/// assert!(is_rune_unicode_assignable(0x41));  // 'A'
/// assert!(!is_rune_unicode_assignable(0xFFFE));  // Noncharacter
/// assert!(!is_rune_unicode_assignable(0xFDD0));  // Noncharacter
/// assert!(!is_rune_unicode_assignable(0x10FFFE)); // Plane 16 noncharacter
/// ```
pub fn is_rune_unicode_assignable(r: u32) -> bool {
    subset_contains_u32(UNICODE_ASSIGNABLES, r)
}

// String validation functions

/// Check if a string contains only Unicode scalar values
///
/// Since Rust strings are guaranteed to be valid UTF-8 and cannot contain
/// surrogates, this function will always return true for valid Rust strings.
/// This function exists for API compatibility with the Go implementation.
///
/// # Examples
///
/// ```
/// use rfc9839::is_string_unicode_scalars;
///
/// assert!(is_string_unicode_scalars("Hello, world!"));
/// assert!(is_string_unicode_scalars("Hello, 世界!"));
/// assert!(is_string_unicode_scalars("🦀 Rust 🦀"));
/// assert!(is_string_unicode_scalars(""));  // Empty string is valid
/// ```
pub fn is_string_unicode_scalars(_s: &str) -> bool {
    // Rust strings are always valid Unicode scalars
    true
}

/// Check if a string contains only valid XML characters
///
/// # Examples
///
/// ```
/// use rfc9839::is_string_xml_chars;
///
/// assert!(is_string_xml_chars("Valid XML text"));
/// assert!(is_string_xml_chars("Line 1\nLine 2"));  // Newlines are valid
/// assert!(is_string_xml_chars("Tab\there"));  // Tabs are valid
/// assert!(!is_string_xml_chars("Null\u{0000}char"));  // Null is invalid
/// assert!(!is_string_xml_chars("Bell\u{0007}char"));  // Bell is invalid
/// ```
pub fn is_string_xml_chars(s: &str) -> bool {
    s.chars().all(is_char_xml_char)
}

/// Check if a string contains only Unicode assignable characters
///
/// # Examples
///
/// ```
/// use rfc9839::is_string_unicode_assignables;
///
/// assert!(is_string_unicode_assignables("Hello, world!"));
/// assert!(is_string_unicode_assignables("Hello, 世界!"));
/// assert!(!is_string_unicode_assignables("Has nonchar\u{FFFE}"));
/// assert!(!is_string_unicode_assignables("Has\u{FDD0}nonchar"));
/// ```
pub fn is_string_unicode_assignables(s: &str) -> bool {
    s.chars().all(is_char_unicode_assignable)
}

// UTF-8 byte slice validation functions

/// Check if a UTF-8 byte slice contains only Unicode scalar values
///
/// Returns false if the bytes are not valid UTF-8. Since Rust's UTF-8 validation
/// already rejects surrogates, valid UTF-8 always contains only Unicode scalars.
///
/// # Examples
///
/// ```
/// use rfc9839::is_utf8_unicode_scalars;
///
/// assert!(is_utf8_unicode_scalars(b"Hello, world!"));
/// assert!(is_utf8_unicode_scalars("Hello, 世界!".as_bytes()));
///
/// // Invalid UTF-8 sequences return false
/// assert!(!is_utf8_unicode_scalars(&[0xFF, 0xFE]));
///
/// // Surrogate-containing UTF-8 returns false
/// assert!(!is_utf8_unicode_scalars(&[0xED, 0xA0, 0x80]));  // U+D800
/// ```
pub fn is_utf8_unicode_scalars(bytes: &[u8]) -> bool {
    std::str::from_utf8(bytes).is_ok()
}

/// Check if a UTF-8 byte slice contains only valid XML characters
///
/// Returns false if the bytes are not valid UTF-8 or contain invalid XML characters.
///
/// # Examples
///
/// ```
/// use rfc9839::is_utf8_xml_chars;
///
/// assert!(is_utf8_xml_chars(b"Valid XML"));
/// assert!(is_utf8_xml_chars(b"Line 1\nLine 2"));
///
/// // Control characters are invalid
/// assert!(!is_utf8_xml_chars(b"Null\x00char"));
///
/// // Invalid UTF-8 returns false
/// assert!(!is_utf8_xml_chars(&[0xFF, 0xFE]));
/// ```
pub fn is_utf8_xml_chars(bytes: &[u8]) -> bool {
    match std::str::from_utf8(bytes) {
        Ok(s) => is_string_xml_chars(s),
        Err(_) => false,
    }
}

/// Check if a UTF-8 byte slice contains only Unicode assignable characters
///
/// Returns false if the bytes are not valid UTF-8 or contain noncharacters.
///
/// # Examples
///
/// ```
/// use rfc9839::is_utf8_unicode_assignables;
///
/// assert!(is_utf8_unicode_assignables(b"Hello, world!"));
/// assert!(is_utf8_unicode_assignables("Hello, 世界!".as_bytes()));
///
/// // Invalid UTF-8 returns false
/// assert!(!is_utf8_unicode_assignables(&[0xFF, 0xFE]));
/// ```
pub fn is_utf8_unicode_assignables(bytes: &[u8]) -> bool {
    match std::str::from_utf8(bytes) {
        Ok(s) => is_string_unicode_assignables(s),
        Err(_) => false,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    // Test data matching Go test file
    const INVERSE_XML: &[RunePair] = &[
        RunePair::new(0, 0x8),
        RunePair::new(0xB, 0xC),
        RunePair::new(0xE, 0x1F),
        RunePair::new(0xD800, 0xDFFF),
        RunePair::new(0xFFFE, 0xFFFF),
    ];

    const INVERSE_ASSIGNABLES: &[RunePair] = &[
        RunePair::new(0, 0x8),
        RunePair::new(0xB, 0xC),
        RunePair::new(0xE, 0x1F),
        RunePair::new(0x7F, 0x9F),
        RunePair::new(0xD800, 0xDFFF),
        RunePair::new(0xFDD0, 0xFDEF),
        RunePair::new(0xFFFE, 0xFFFF),
        RunePair::new(0x1FFFE, 0x1FFFF),
        RunePair::new(0x2FFFE, 0x2FFFF),
        RunePair::new(0x3FFFE, 0x3FFFF),
        RunePair::new(0x4FFFE, 0x4FFFF),
        RunePair::new(0x5FFFE, 0x5FFFF),
        RunePair::new(0x6FFFE, 0x6FFFF),
        RunePair::new(0x7FFFE, 0x7FFFF),
        RunePair::new(0x8FFFE, 0x8FFFF),
        RunePair::new(0x9FFFE, 0x9FFFF),
        RunePair::new(0xAFFFE, 0xAFFFF),
        RunePair::new(0xBFFFE, 0xBFFFF),
        RunePair::new(0xCFFFE, 0xCFFFF),
        RunePair::new(0xDFFFE, 0xDFFFF),
        RunePair::new(0xEFFFE, 0xEFFFF),
        RunePair::new(0xFFFFE, 0xFFFFF),
        RunePair::new(0x10FFFE, 0x10FFFF), // Note: Fixed typo from Go test (was 0xFFFFF)
    ];

    #[test]
    fn test_empties() {
        let empty_u1: &[u8] = &[];
        let empty_u2: Option<&[u8]> = None;
        let empty_s = "";

        // Test empty byte slice
        assert!(is_utf8_unicode_scalars(empty_u1));
        assert!(is_utf8_xml_chars(empty_u1));
        assert!(is_utf8_unicode_assignables(empty_u1));

        // Test None case (treating as empty)
        if let Some(bytes) = empty_u2 {
            assert!(is_utf8_unicode_scalars(bytes));
            assert!(is_utf8_xml_chars(bytes));
            assert!(is_utf8_unicode_assignables(bytes));
        } else {
            // None is treated as empty, which should pass
            assert!(is_utf8_unicode_scalars(&[]));
            assert!(is_utf8_xml_chars(&[]));
            assert!(is_utf8_unicode_assignables(&[]));
        }

        // Test empty string
        assert!(is_string_unicode_scalars(empty_s));
        assert!(is_string_xml_chars(empty_s));
        assert!(is_string_unicode_assignables(empty_s));
    }

    #[test]
    fn test_scalars() {
        // Test valid scalar values
        assert!(is_rune_unicode_scalar(0x41)); // 'A'
        assert!(is_rune_unicode_scalar(0x4E2D)); // '中'
        assert!(is_rune_unicode_scalar(0x10000)); // U+10000
        assert!(is_rune_unicode_scalar(0x10FFFF)); // Max valid Unicode

        // Test invalid scalar ranges (surrogates)
        for r in 0xD800..=0xDFFF {
            assert!(
                !is_rune_unicode_scalar(r),
                "{:x} should not be unicode scalar",
                r
            );
        }

        // Test boundary conditions
        assert!(!is_rune_unicode_scalar(0xFFFFFFFF)); // -1 in u32
        assert!(!is_rune_unicode_scalar(0x110000)); // 0x10FFFF + 1

        // Test that Rust strings are always Unicode scalars
        assert!(is_string_unicode_scalars("Hello, 世界!"));
        assert!(is_string_unicode_scalars("")); // Empty string

        // Test invalid UTF-8 (surrogate)
        let bad_utf8 = vec![0xED, 0xBA, 0xAD]; // U+DEAD (surrogate)
        let mut bad = vec![b'a'];
        bad.extend_from_slice(&bad_utf8);
        bad.push(b'z');

        assert!(!is_utf8_unicode_scalars(&bad));
    }

    #[test]
    fn test_xml_chars() {
        // Test all valid XML char ranges (with sampling for performance)
        for pair in XML_CHARS {
            // Test boundaries and a sample in the middle
            assert!(is_rune_xml_char(pair.lo), "{:x} should be XML", pair.lo);
            assert!(is_rune_xml_char(pair.hi), "{:x} should be XML", pair.hi);
            if pair.hi - pair.lo > 2 {
                let mid = (pair.lo + pair.hi) / 2;
                assert!(is_rune_xml_char(mid), "{:x} should be XML", mid);
            }
        }

        // Test invalid XML char ranges
        for pair in INVERSE_XML {
            for r in pair.lo..=pair.hi.min(pair.lo + 100) {
                // Sample for performance
                assert!(!is_rune_xml_char(r), "{:x} should not be XML", r);
            }
        }

        // Test boundary conditions
        assert!(!is_rune_xml_char(0xFFFFFFFF)); // -1 in u32
        assert!(!is_rune_xml_char(0x110000)); // 0x10FFFF + 1

        // Test invalid UTF-8 (surrogate)
        let bad_utf8 = vec![0xED, 0xBA, 0xAD]; // U+DEAD
        let mut bad = vec![b'a'];
        bad.extend_from_slice(&bad_utf8);
        bad.push(b'z');

        assert!(!is_utf8_xml_chars(&bad));

        // Test good strings
        let mut good_s = String::new();
        let mut good_u = Vec::new();

        for pair in XML_CHARS {
            if let Some(lo_char) = char::from_u32(pair.lo) {
                good_s.push(lo_char);
                let mut buf = [0; 4];
                let encoded = lo_char.encode_utf8(&mut buf);
                good_u.extend_from_slice(encoded.as_bytes());
            }
            if let Some(hi_char) = char::from_u32(pair.hi) {
                good_s.push(hi_char);
                let mut buf = [0; 4];
                let encoded = hi_char.encode_utf8(&mut buf);
                good_u.extend_from_slice(encoded.as_bytes());
            }
        }

        assert!(is_string_xml_chars(&good_s));
        assert!(is_utf8_xml_chars(&good_u));

        // Test invalid characters
        for pair in INVERSE_XML {
            if pair.lo == 0xD800 {
                // Skip surrogates (can't create valid Rust char)
                continue;
            }
            if let Some(ch) = char::from_u32(pair.lo) {
                let mut bad = vec![b'a'];
                let mut buf = [0; 4];
                let encoded = ch.encode_utf8(&mut buf);
                bad.extend_from_slice(encoded.as_bytes());
                bad.push(b'z');

                assert!(
                    !is_utf8_xml_chars(&bad),
                    "accepted utf8 containing {:x}",
                    pair.lo
                );
                assert!(
                    !is_string_xml_chars(&String::from_utf8_lossy(&bad)),
                    "accepted string containing {:x}",
                    pair.lo
                );
            }
        }
    }

    #[test]
    fn test_assignables() {
        // Test all valid assignable ranges (with sampling for performance)
        for pair in UNICODE_ASSIGNABLES {
            // Test boundaries
            assert!(
                is_rune_unicode_assignable(pair.lo),
                "{:x} should be Assignable",
                pair.lo
            );
            assert!(
                is_rune_unicode_assignable(pair.hi),
                "{:x} should be Assignable",
                pair.hi
            );
        }

        // Test invalid assignable ranges
        for pair in INVERSE_ASSIGNABLES {
            // Sample a few values from each range for performance
            let test_values = vec![
                pair.lo,
                pair.hi,
                if pair.hi > pair.lo {
                    (pair.lo + pair.hi) / 2
                } else {
                    pair.lo
                },
            ];
            for r in test_values {
                assert!(
                    !is_rune_unicode_assignable(r),
                    "{:x} should not be Assignable",
                    r
                );
            }
        }

        // Test boundary conditions
        assert!(!is_rune_unicode_assignable(0xFFFFFFFF)); // -1 in u32
        assert!(!is_rune_unicode_assignable(0x110000)); // 0x10FFFF + 1

        // Test invalid UTF-8 (surrogate)
        let bad_utf8 = vec![0xED, 0xBA, 0xAD]; // U+DEAD
        let mut bad = vec![b'a'];
        bad.extend_from_slice(&bad_utf8);
        bad.push(b'z');

        assert!(!is_utf8_unicode_assignables(&bad));

        // Test good strings
        let mut good_s = String::new();
        let mut good_u = Vec::new();

        for pair in UNICODE_ASSIGNABLES {
            if let Some(lo_char) = char::from_u32(pair.lo) {
                good_s.push(lo_char);
                let mut buf = [0; 4];
                let encoded = lo_char.encode_utf8(&mut buf);
                good_u.extend_from_slice(encoded.as_bytes());
            }
            if let Some(hi_char) = char::from_u32(pair.hi) {
                good_s.push(hi_char);
                let mut buf = [0; 4];
                let encoded = hi_char.encode_utf8(&mut buf);
                good_u.extend_from_slice(encoded.as_bytes());
            }
        }

        assert!(is_string_unicode_assignables(&good_s));
        assert!(is_utf8_unicode_assignables(&good_u));

        // Test invalid characters
        for pair in INVERSE_ASSIGNABLES {
            if pair.lo == 0xD800 {
                // Skip surrogates (can't create valid Rust char)
                continue;
            }
            if let Some(ch) = char::from_u32(pair.lo) {
                let mut bad = vec![b'a'];
                let mut buf = [0; 4];
                let encoded = ch.encode_utf8(&mut buf);
                bad.extend_from_slice(encoded.as_bytes());
                bad.push(b'z');

                assert!(
                    !is_utf8_unicode_assignables(&bad),
                    "accepted utf8 containing {:x}",
                    pair.lo
                );
                assert!(
                    !is_string_unicode_assignables(&String::from_utf8_lossy(&bad)),
                    "accepted string containing {:x}",
                    pair.lo
                );
            }
        }
    }
}
