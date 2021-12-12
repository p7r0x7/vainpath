package vainpath

import (
	"path/filepath"
	"strings"
)

// Copyright © 2021 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

// Clean returns a formatted version of `path` based off of the fish shell's method to shorten
// paths as used in `fish/functions/prompt_pwd.fish`; it will almost certainly not return valid
// paths and should be used for vanity purposes only.
func Clean(path string) string {
	path = filepath.ToSlash(filepath.Clean(path)) /* Windows-sensitive */
	segments := strings.Split(path, "/")

	/* Skips final index */
	for i := len(segments) - 2; i >= 0; i-- {
		switch segments[i][:1] {
		/* RegEx equivalent is [ [:punct:]] */
		case " ", "]", "[", "!", "\"", "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", ".",
			"/", ":", ";", "<", "=", ">", "?", "@", "\\", "^", "_", "`", "{", "|", "}", "~", "-":
			if len(segments[i]) > 1 {
				segments[i] = segments[i][:2]
			}
		default:
			segments[i] = segments[i][:1]
		}
	}

	if path[:1] == "/" {
		segments[0] = "/" + segments[0]
	}
	return strings.Join(segments, string(filepath.Separator)) /* Windows-sensitive */
}
