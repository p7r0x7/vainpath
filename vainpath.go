package vainpath

import (
	"path/filepath"
	"regexp"
	"strings"
)

// Copyright © 2021 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

// Clean returns a formatted version of `path` based off of the fish shell's method to shorten
// paths as used in `fish/functions/prompt_pwd.fish`; it will almost certainly not return valid
// paths and should be used for vanity purposes only.
func Clean(path string) string {
	path = filepath.ToSlash(filepath.Clean(path)) /* Windows-sensitive */
	segments := strings.Split(path, "/")

	exp := regexp.MustCompile("[ [:punct:]]?.")
	/* Skips final index */
	for i := len(segments) - 2; i >= 0; i-- {
		segments[i] = exp.FindString(segments[i])
	}

	if path[:1] == "/" {
		segments[0] = "/" + segments[0]
	}
	return filepath.Join(segments...) /* Windows-sensitive */
}
