package vainpath

import (
	"math/rand"
	. "path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"
	. "unicode/utf8"
)

// Copyright © 2022 Matthew R Bonnette. Licensed under a BSD-3-Clause license.

func init() {
	rand.Seed(time.Now().UnixNano())
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func randPath(size int) string {
	var path strings.Builder
	for path.Grow(size); size > 0; size-- {
		if r := rand.Int31(); r%32 == 0 {
			path.WriteByte(Separator)

		} else if r%8 == 0 && size > 1 {
			for r = 0; !unicode.IsPrint(r); {
				switch min(4, size) {
				case 4:
					r = rand.Int31() & '\U0010FFFF'
				case 3:
					r = rand.Int31() >> 8 & '\uFFFF'
				case 2:
					r = rand.Int31() >> 16 & '\u07FF'
				}
			}
			size -= min(4, size) - 1
			path.WriteRune(r)

		} else {
			path.WriteByte(byte(rand.Int31()%95 + 32))
		}
	}

	return path.String()
}

func TestValidate(t *testing.T) {
	space, log := strings.Repeat(" ", 40), strings.Builder{}

	for i := 0; i < 1e5; i++ {
		path := randPath(50)
		if s, r := Shorten(path), refShorten(path); s != r {
			if log.Len() < 4*1024 {
				dex := strconv.Itoa(i)
				for _, v := range []string{
					"\n", space[:5-len(dex)], dex, "  ",
					s, space[:40-min(40, RuneCountInString(s))], " \t",
					r, space[:40-min(40, RuneCountInString(r))], " \t",
					path,
				} {
					log.WriteString(v)
				}
			}
			t.Fail()
		}
	}
	if t.Failed() {
		t.Log(log.String())
	}
}

func BenchmarkAvgClean(b *testing.B) {
	const bestPath = string(Separator)
	const worstPath = "\u07FF " + string(Separator)
	fast := strings.Repeat(bestPath, b.N)
	slow := strings.Repeat(worstPath, b.N/len(worstPath)+1)

	b.SetBytes(int64(b.N))
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Shorten(fast[:i])
		Shorten(slow[:i])
	}
}

func BenchmarkRandClean(b *testing.B) {
	path := randPath(120)
	b.SetBytes(120)
	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		Shorten(path)
	}
}

func refShorten(path string) string {
	segments := strings.Split(Clean(path), string(Separator))

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

	return strings.Join(segments, string(Separator))
}
