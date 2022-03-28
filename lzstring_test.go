package lzstring

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		arg  string
		want []uint16
	}{
		{
			arg:  "",
			want: []uint16{},
		},
		{
			arg:  "H",
			want: []uint16{0x490},
		},
		{
			arg:  "HelloHello",
			want: []uint16{0x485, 0x3036, 0x60f6, 0xa194, 0x0},
		},
		{
			arg:  "ababcabcdabcde",
			want: []uint16{0x2182, 0x3518, 0xc204, 0xda14, 0xc800},
		},
		{
			arg:  "Hello, world",
			want: []uint16{0x485, 0x3036, 0x60f6, 0x340, 0x40e, 0xe901, 0x3980, 0x2640},
		},
		{
			arg:  "あいうえお",
			want: []uint16{0x9083, 0x2110, 0x6462, 0xc81, 0x20c8, 0x520c, 0x4000},
		},
		{
			arg:  "🍎",
			want: []uint16{'輆', '', 0xda00},
		},
		{
			arg:  "🍎🍇",
			want: []uint16{0x8f06, 0xe397, 0xde9c, 0x5f68},
		},
		{
			arg:  "aあ🍎bい🍇c",
			want: []uint16{0x21a2, 0x1064, 0x3c1b, 0x872f, 0xb046, 0x8220, 0xc68e, 0x2fb0, 0x6320},
		},
		{
			arg:  string([]rune{0x9c}),
			want: []uint16{0xe50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := Compress(tt.arg)
			assert.Nil(t, err)
			b, err := reader, err
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestCompressToBase64(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{
			arg:  "",
			want: "Q===",
		},
		{
			arg:  "H",
			want: "BJA=",
		},
		{
			arg:  "HelloHello",
			want: "BIUwNmD2oZQ=",
		},
		{
			arg:  "ababcabcdabcde",
			want: "IYI1GMIE2hTI",
		},
		{
			arg:  "Hello, world",
			want: "BIUwNmD2A0AEDukBOYAmQ===",
		},
		{
			arg:  "あいうえお",
			want: "kIMhEGRiDIEgyFIMQ===",
		},
		{
			arg:  "🍎",
			want: "jwbjl9o=",
		},
		{
			arg:  "🍎🍇",
			want: "jwbjl96cX2g=",
		},
		{
			arg:  "aあ🍎bい🍇c",
			want: "IaIQZDwbhy+wRoIgxo4vsGMg",
		},
		{
			arg:  string([]rune{0x9c}),
			want: "DlA=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := CompressToBase64(tt.arg)
			assert.Nil(t, err)
			b, err := reader, err
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestCompressToUTF16(t *testing.T) {
	tests := []struct {
		arg  string
		want []uint16
	}{
		{
			arg:  "",
			want: []uint16{0x2020, 0x20},
		},
		{
			arg:  "H",
			want: []uint16{0x268, 0x20},
		},
		{
			arg:  "HelloHello",
			want: []uint16{0x262, 0x4c2d, 0x4c3e, 0x6a39, 0x2020, 0x20},
		},
		{
			arg:  "ababcabcdabcde",
			want: []uint16{0x10e1, 0xd66, 0x1860, 0x4dc1, 0x2660, 0x20},
		},
		{
			arg:  "Hello, world",
			want: []uint16{0x262, 0x4c2d, 0x4c3e, 0x6054, 0x40, 0x3bc4, 0x293, 0x46, 0x2020, 0x20},
		},
		{
			arg:  "あいうえお",
			want: []uint16{0x4861, 0x4864, 0xcac, 0x20e8, 0x926, 0x2168, 0x18a0, 0x20},
		},
		{
			arg:  "🍎",
			want: []uint16{0x47a3, 0x3905, 0x7b60, 0x20},
		},
		{
			arg:  "🍎🍇",
			want: []uint16{0x47a3, 0x3905, 0x7bf3, 0x4616, 0x4020, 0x20},
		},
		{
			arg:  "aあ🍎bい🍇c",
			want: []uint16{0x10f1, 0x439, 0x7a3, 0x3892, 0x7da2, 0x1a28, 0x41ad, 0xe4f, 0x5851, 0x4820, 0x20},
		},
		{
			arg:  string([]rune{0x9c}),
			want: []uint16{0x748, 0x20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := CompressToUTF16(tt.arg)
			assert.Nil(t, err)
			b, err := reader, err
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestDecompress(t *testing.T) {
	tests := []struct {
		arg  []uint16
		want string
	}{
		{
			arg:  []uint16{},
			want: "",
		},
		{
			arg:  []uint16{0x490},
			want: "H",
		},
		{
			arg:  []uint16{0x485, 0x3036, 0x60f6, 0xa194, 0x0},
			want: "HelloHello",
		},
		{
			arg:  []uint16{0x2182, 0x3518, 0xc204, 0xda14, 0xc800},
			want: "ababcabcdabcde",
		},
		{
			arg:  []uint16{0x485, 0x3036, 0x60f6, 0x340, 0x40e, 0xe901, 0x3980, 0x2640},
			want: "Hello, world",
		},
		{
			arg:  []uint16{0x9083, 0x2110, 0x6462, 0xc81, 0x20c8, 0x520c, 0x4000},
			want: "あいうえお",
		},
		{
			arg:  []uint16{0x8f06, 0xe397, 0xda00},
			want: "🍎",
		},
		{
			arg:  []uint16{0x8f06, 0xe397, 0xde9c, 0x5f68},
			want: "🍎🍇",
		},
		{
			arg:  []uint16{0x21a2, 0x1064, 0x3c1b, 0x872f, 0xb046, 0x8220, 0xc68e, 0x2fb0, 0x6320},
			want: "aあ🍎bい🍇c",
		},
		{
			arg:  []uint16{0xe50},
			want: string([]rune{0x9c}),
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.arg), func(t *testing.T) {
			res, err := Decompress(tt.arg)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestDecompressFromBase64(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{
			arg:  "Q===",
			want: "",
		},
		{
			arg:  "BJA=",
			want: "H",
		},
		{
			arg:  "BIUwNmD2oZQ=",
			want: "HelloHello",
		},
		{
			arg:  "IYI1GMIE2hTI",
			want: "ababcabcdabcde",
		},
		{
			arg:  "BIUwNmD2A0AEDukBOYAmQ===",
			want: "Hello, world",
		},
		{
			arg:  "kIMhEGRiDIEgyFIMQ===",
			want: "あいうえお",
		},
		{
			arg:  "jwbjl9o=",
			want: "🍎",
		},
		{
			arg:  "jwbjl96cX2g=",
			want: "🍎🍇",
		},
		{
			arg:  "IaIQZDwbhy+wRoIgxo4vsGMg",
			want: "aあ🍎bい🍇c",
		},
		{
			arg:  "DlA=",
			want: string([]rune{0x9c}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := DecompressFromBase64(tt.arg)
			assert.Nil(t, err)
			b, err := reader, err
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestDecompressFromUTF16(t *testing.T) {
	tests := []struct {
		arg  []uint16
		want string
	}{
		{
			arg:  []uint16{0x2020, 0x20},
			want: "",
		},
		{
			arg:  []uint16{0x268, 0x20},
			want: "H",
		},
		{
			arg:  []uint16{0x262, 0x4c2d, 0x4c3e, 0x6a39, 0x2020, 0x20},
			want: "HelloHello",
		},
		{
			arg:  []uint16{0x10e1, 0xd66, 0x1860, 0x4dc1, 0x2660, 0x20},
			want: "ababcabcdabcde",
		},
		{
			arg:  []uint16{0x262, 0x4c2d, 0x4c3e, 0x6054, 0x40, 0x3bc4, 0x293, 0x46, 0x2020, 0x20},
			want: "Hello, world",
		},
		{
			arg:  []uint16{0x4861, 0x4864, 0xcac, 0x20e8, 0x926, 0x2168, 0x18a0, 0x20},
			want: "あいうえお",
		},
		{
			arg:  []uint16{0x47a3, 0x3905, 0x7b60, 0x20},
			want: "🍎",
		},
		{
			arg:  []uint16{0x47a3, 0x3905, 0x7bf3, 0x4616, 0x4020, 0x20},
			want: "🍎🍇",
		},
		{
			arg:  []uint16{0x10f1, 0x439, 0x7a3, 0x3892, 0x7da2, 0x1a28, 0x41ad, 0xe4f, 0x5851, 0x4820, 0x20},
			want: "aあ🍎bい🍇c",
		},
		{
			arg:  []uint16{0x748, 0x20},
			want: string([]rune{0x9c}),
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.arg), func(t *testing.T) {
			res, err := DecompressFromUTF16(tt.arg)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, res)
		})
	}
}

func FuzzIntegrity(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := Compress(s)
		if !utf8.ValidString(s) {
			assert.NotNil(t, err)
			return
		}
		assert.Nil(t, err)
		repair, err := Decompress(compressed)
		assert.Nil(t, err)
		assert.Equal(t, s, repair)
	})
}
