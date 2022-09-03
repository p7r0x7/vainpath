package vainpath

import (
	"fmt"
	"math/rand"
	. "path/filepath"
	"strconv"
	"strings"
	"testing"
	"unicode"
	. "unicode/utf8"
)

// Copyright © 2022 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

// randPath produces random valid strings of a given byte length.
func randPath(size int) string {
	var path strings.Builder
	for path.Grow(size); size > 0; {
		r := rand.Int31() % 100
		switch {
		case r > 10:
			path.WriteByte(byte(rand.Int31()%95 + 32))
			size--
			continue
		case r > 5 && size >= 2:
			for !unicode.IsPrint(r) {
				r = rand.Int31() >> 16 & '\u07FF'
			}
		case r > 1 && size >= 3:
			for !unicode.IsPrint(r) {
				r = rand.Int31() >> 8 & '\uFFFF'
			}
		case size >= 4:
			for !unicode.IsPrint(r) {
				r = rand.Int31() & '\U0010FFFF'
			}
		default:
			continue
		}
		size -= RuneLen(r)
		path.WriteRune(r)
	}
	return path.String()
}

// refSimplify defines the expected behavior of Simplify.
func refSimplify(path string) string {
	if !ValidString(path) {
		return path
	}
	segments := strings.Split(Clean(path), string(Separator))

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

	return strings.Join(segments, string(Separator))
}

func TestForPanic(t *testing.T) {
	t.Parallel()
	for i := 0; i < 1024; i++ {
		bytes := make([]byte, 1024)
		rand.Read(bytes)
		for i2 := range bytes[:len(bytes)-8] {
			Trim(string(bytes[i2:i2+8]), "", 8)
			Simplify(string(bytes[i2 : i2+8]))
		}
	}
}

// TestValidate tests if Simplify behaves identically to refSimplify for a sufficiently large number
// of valid, variably-sized inputs.
func TestIfValid(t *testing.T) {
	t.Parallel()
	space, log := strings.Repeat(" ", 40), strings.Builder{}

	for i := 0; i < 1e5; i++ {
		path := randPath(rand.Int() % 60)
		if s, r := Simplify(path), refSimplify(path); s != r {
			if log.Len() < 4*1024 {
				dex := strconv.Itoa(i)
				fmt.Fprint(&log, "\n", space[:5-len(dex)], dex, space[:2],
					s, space[:40-min(40, RuneCountInString(s))], " \t",
					r, space[:40-min(40, RuneCountInString(r))], " \t",
					path)
				t.Fail()
			} else {
				break
			}
		}
	}
	if t.Failed() {
		t.Log(log.String())
	}
}

func TestAllocRate(t *testing.T) {
	t.Parallel()
	for _, v := range [...]int{64, 128, 256, 512, 1024} {
		var once, alloced, actual int
		for i := 1e4; i > 0; i-- {
			path := randPath(v)
			dex := strings.LastIndexByte(path, Separator) + 1
			alloc, ln := min(v>>3+v-dex, v), len(Simplify(path))
			alloced += alloc
			actual += ln
			if alloc >= ln {
				once++
			}
		}
		t.Log(strconv.FormatFloat(float64(once)/1e2, 'f', 2, 64)+"%",
			strconv.FormatFloat(float64(alloced)/float64(actual), 'f', 4, 64))
	}
}

func BenchmarkAvgTrim(b *testing.B) {
	const bestPath, worstPath = "\U0010FFFF", " "
	fast := strings.Repeat(bestPath, b.N)
	slow := strings.Repeat(worstPath, b.N/len(worstPath)+1)

	b.SetBytes(int64(b.N/2 + 2))
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Trim(fast[:i], "", b.N)
		Trim(slow[:i], " ", 2)
	}
}

func BenchmarkRandTrim(b *testing.B) {
	path := randPath(b.N)
	n := RuneCountInString(path) - 1
	b.SetBytes(int64(len(Trim(path, "…", n))))
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Trim(path, "…", n)
	}
}

func BenchmarkAvgSimplify(b *testing.B) {
	const bestPath = string(Separator)
	const worstPath = string(Separator) + "\u07FF "
	fast := strings.Repeat(bestPath, b.N)
	slow := strings.Repeat(worstPath, b.N/len(worstPath)+1)

	b.SetBytes(int64(b.N))
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Simplify(fast[:i])
		Simplify(slow[:i])
	}
}

func BenchmarkRandSimplify(b *testing.B) {
	path := randPath(b.N)
	b.SetBytes(int64(b.N))
	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Simplify(path)
	}
}

func BenchmarkRandPath(b *testing.B) {
	b.SetBytes(1)
	randPath(b.N)
}
