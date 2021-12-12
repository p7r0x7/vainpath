package vainpath

import (
	"encoding/ascii85"
	"math/rand"
	"testing"
)

// Copyright © 2021 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

func BenchmarkClean(b *testing.B) {
	const length, segments = 25, 4
	var path string

	for i := segments; i > 0; i-- {
		bytes := make([]byte, length)
		rand.Read(bytes)
		str := make([]byte, ascii85.MaxEncodedLen(length))
		ascii85.Encode(str, bytes)
		path += string(str) + "/"
	}

	b.SetBytes(int64(len(path)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Clean(path)
	}
}
