package vainpath

import (
	. "path/filepath"
	"strings"
	"unicode"
	. "unicode/utf8"
)

// Copyright © 2022 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

const fullRuneWidth = "" +
	// Valid or Overlong-Encoded range (2-6 bytes)
	"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02" +
	"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02" +
	"\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03" +
	"\x04\x04\x04\x04\x04\x04\x04\x04\x05\x05\x05\x05\x06\x06\x01\x01"

// Shorten formats inputs in a way similar to the fish shell's method of shortening paths in
// 'fish/functions/prompt_pwd.fish' and is properly Windows-sensitive; it will almost certainly not
// return valid paths and should be used for vanity purposes only. Inputs are assumed to be valid
// UTF-8-encoded strings.
func Shorten(path string) string {
	path = Clean(path)
	/* Bytes of ASCII and non-ASCII code points cannot be mistaken for each other. */
	dex := strings.LastIndexByte(path, Separator) + 1 /* There will always be more bytes. */
	if len(path) < 4 || dex < 3 {
		// The following already-cleaned paths cannot be any more shortened:
		// a) paths smaller than 4 bytes
		// b) paths without separators (dex of 0)
		// c) paths where the only separator is the first or second character (dex of 1 or 2)
		return path
	}

	out, prefix, w := strings.Builder{}, path[:dex], 0
	out.Grow(len(path) / 3) /* Minimizes re-allocation without expensive calculation. */

	/* For performance reasons, this section is more verbose than strictly necessary. */
	if prefix[0] == Separator {
		w++
	}
	for len(prefix[w:]) > 0 {
		r, t := DecodeRuneInString(prefix[w:])
		if w += t; prefix[w] == Separator {
			w++
			continue
		}
		if !unicode.IsLetter(r) {
			_, t = DecodeRuneInString(prefix[w:])
			if w += t; prefix[w] == Separator {
				w++
				continue
			}
		}
		out.WriteString(prefix[:w])
		for {
			if c := prefix[w]; c < 192 {
				if c == Separator {
					break
				}
				w++ /* Invalid leading bytes are skipped. */
			} else {
				w += int(fullRuneWidth[c-192])
			}
		}
		prefix, w = prefix[w:], 1
	}

	out.WriteString(path[dex-w:])
	return out.String()
}
