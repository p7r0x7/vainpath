package vainpath

import (
	"math/bits"
	. "path/filepath"
	"strings"
	"unicode"
	. "unicode/utf8"
)

// Copyright © 2022 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

// Shorten formats inputs in a way similar to the fish shell's method of shortening paths in
// 'fish/functions/prompt_pwd.fish' and is properly Windows-sensitive; it will almost certainly not
// return valid paths and should be used for vanity purposes only. Inputs are assumed to be valid
// UTF-8-encoded strings; behavior is undefined for other inputs.
func Shorten(path string) string {
	path = Clean(path)
	/* Bytes of ASCII and non-ASCII code points cannot be mistaken for each other. */
	dex := strings.LastIndexByte(path, Separator)
	if len(path) < 4 || dex < 2 {
		// The following already-cleaned paths cannot be any more shortened:
		// a) paths smaller than 4 bytes
		// b) paths without separators (dex of -1)
		// c) paths where the only separator is the first or second character (dex of 0 or 1)
		return path
	}

	/* There will always be more bytes. */
	out, prefix, suffix, w := strings.Builder{}, path[:dex+1], len(path[dex+1:]), 0
	/* Max memory requirement: root + sepCount * twoRunesAndSep + suffix. */
	out.Grow(1 + strings.Count(path, string(Separator))*9 + suffix)

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
			if c := prefix[w]; c < RuneSelf {
				if c == Separator {
					break
				}
				w++
			} else {
				w += bits.LeadingZeros8(^c)
			}
		}
		prefix, w = prefix[w:], 1
	}

	out.WriteString(path[len(path)-w-suffix:])
	return out.String()
}
