package vainpath

import (
	"encoding/ascii85"
	"math/rand"
	"testing"
)

// Copyright © 2021 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

func BenchmarkClean(b *testing.B) {
	bytes1 := make([]byte, 100)
	rand.Read(bytes1)
	bytes2 := make([]byte, ascii85.MaxEncodedLen(100))
	ascii85.Encode(bytes2, bytes1)
	path := string(bytes2)

	b.SetBytes(int64(ascii85.MaxEncodedLen(100)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Clean(path)
	}
}
