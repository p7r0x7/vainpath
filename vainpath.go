package vainpath

import (
	"path/filepath"
	"strings"
	"unicode"
	. "unicode/utf8"
)

// Copyright © 2022 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

// Clean formats inputs in a way similar to the fish shell's method of shortening paths in
// 'fish/functions/prompt_pwd.fish' and is properly Windows-sensitive; it will almost certainly not
// return valid paths and should be used for vanity purposes only. Inputs are assumed to be valid
// UTF-8-encoded strings; behavior is undefined for other inputs.
func Clean(path string) string {
	path = filepath.Clean(path)
	/* Bytes of ASCII and non-ASCII code points cannot be mistaken for each other. */
	dex := strings.LastIndexByte(path, filepath.Separator)
	if len(path) < 4 || dex < 2 {
		// The following cannot be any more shortened:
		// a) paths smaller than 4 bytes
		// b) paths without separators (dex of -1)
		// c) paths where the only separator is the first or second character (dex of 0 or 1)
		return path
	}

	/* Max memory requirement: root + sepCount * twoRunesAndSep + suffix. */
	out, suffix := strings.Builder{}, path[dex+1:] /* There will always be more bytes. */
	out.Grow(1 + strings.Count(path, string(filepath.Separator))*9 + len(suffix))
	path = path[:dex+1]

	/* For performance reasons, this section is more verbose than strictly necessary. */
	for len(path) > 0 {
		r, w := DecodeRuneInString(path)
		out.WriteString(path[:w])

		if !unicode.IsLetter(r) {
			path = path[w:]
			if path[0] == filepath.Separator {
				out.WriteByte(filepath.Separator)
				path = path[1:]
				continue
			}
			_, w = DecodeRuneInString(path)
			out.WriteString(path[:w])
		}

		for path[w] != filepath.Separator {
			w++
		}
		out.WriteByte(filepath.Separator)
		path = path[w+1:]
	}

	out.WriteString(suffix)
	return out.String()
}
