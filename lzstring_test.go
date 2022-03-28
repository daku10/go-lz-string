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
			res := Decompress(tt.arg)
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
		repair := Decompress(compressed)
		assert.Equal(t, s, repair)
	})
}
