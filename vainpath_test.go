package vainpath

import (
	"encoding/ascii85"
	"math/rand"
	"path/filepath"
	"testing"
	"time"
	"unicode/utf8"
)

// Copyright © 2022 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

const unilen, asclen, segments = 5, 20, 4

func init() {
	rand.Seed(time.Now().UnixNano())
}

func BenchmarkClean(b *testing.B) {
	var (
		path  string
		char  rune
		bytes = make([]byte, asclen/5*4-3) /* Encoding adds length */
		enc   = make([]byte, asclen)
	)

	for i := segments; i > 0; i-- {
		for i2 := unilen; i2 > 0; i2-- {
			char = rand.Int31()
			for !utf8.ValidRune(char) {
				char = rand.Int31()
			}
			path += string(char)
		}
		rand.Read(bytes)
		ascii85.Encode(enc, bytes)
		path += string(enc) + string(filepath.Separator)
	}

	// fmt.Println(path + "\n" + Clean(path))
	b.SetBytes(int64(len(path)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Clean(path)
	}
}
