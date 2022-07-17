package vainpath

import (
	"math/rand"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode"
	. "unicode/utf8"
)

// Copyright © 2022 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

var worstCase = strings.Repeat(string([]rune{'߿', filepath.Separator}), 1e4)
var bestCase = strings.Repeat(string(filepath.Separator), 3e4)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func makePath(utf, asc, seg int) string {
	var path strings.Builder
	path.Grow((1 + utf*4 + asc) * seg)

	for i := seg; i > 0; i-- {
		for i2 := utf; i2 > 0; i2-- {
			r := rand.Int31()
			for !ValidRune(r) {
				r = rand.Int31()
			}
			path.WriteRune(r)
		}
		for i2 := asc; i2 > 0; i2-- {
			path.WriteByte(byte(rand.Int()%95 + 32))
		}
		path.WriteByte(filepath.Separator)
	}

	return path.String()
}

func TestValidate(t *testing.T) {
	const count int = 1e4
	for i := count; i > 0; i-- {
		path := makePath(2, 20, 4)
		if r, c := refClean(path), Clean(path); r != c {
			t.Log("path: " + path + " ref: " + r + " clean: " + c)
			t.Fail()
		}
	}
}

func BenchmarkAvgClean(b *testing.B) {
	b.SetBytes(int64(b.N))
	for i := b.N; i > 0; i-- {
		refClean(bestCase[:i])
		refClean(worstCase[:i])
	}
}

func BenchmarkRandClean(b *testing.B) {
	path := makePath(2, 200, 4)
	b.SetBytes(int64(len(path)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		refClean(path)
	}
}

func refClean(path string) string {
	segments := strings.Split(filepath.Clean(path), string(filepath.Separator))

	/* Skips final index */
	for i, v := range segments[:len(segments)-1] {
		if RuneCountInString(v) < 2 {
			continue
		}

		r, w0 := DecodeRuneInString(v)
		if unicode.IsLetter(r) {
			segments[i] = v[:w0]
		} else {
			_, w1 := DecodeRuneInString(v[w0:])
			segments[i] = v[:w0+w1]
		}
	}

	return strings.Join(segments, string(filepath.Separator))
}
