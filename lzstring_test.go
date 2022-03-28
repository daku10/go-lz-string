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
		want []rune
	}{
		{
			arg:  "",
			want: []rune{},
		},
		{
			arg:  "H",
			want: []rune("Ґ"),
		},
		{
			arg:  "HelloHello",
			want: []rune{0x485, 0x3036, 0x60f6, 0xa194, 0x0},
		},
		{
			arg:  "ababcabcdabcde",
			want: []rune{0x2182, 0x3518, 0xc204, 0xda14, 0xc800},
		},
		{
			arg:  "Hello, world",
			want: []rune("҅〶惶̀Ў㦀♀"),
		},
		{
			arg:  "あいうえお",
			want: []rune("邃ℐ摢ಁ⃈刌䀀"),
		},
		{
			arg:  "🍎",
			want: []rune{'輆', '', 0xda00},
		},
		{
			arg:  "🍎🍇",
			want: []rune{0x8f06, 0xe397, 0xde9c, 0x5f68},
		},
		{
			arg:  "aあ🍎bい🍇c",
			want: []rune{0x21a2, 0x1064, 0x3c1b, 0x872f, 0xb046, 0x8220, 0xc68e, 0x2fb0, 0x6320},
		},
		{
			arg:  string([]rune{0x9c}),
			want: []rune{0xe50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := Compress(tt.arg)
			assert.Nil(t, err)
			b, err := reader, err
			assert.Equal(t, tt.want, []rune(b))
		})
	}
}

func TestDecompress(t *testing.T) {
	tests := []struct {
		arg  []rune
		want string
	}{
		{
			arg:  []rune{},
			want: "",
		},
		{
			arg:  []rune("Ґ"),
			want: "H",
		},
		{
			arg:  []rune{0x485, 0x3036, 0x60f6, 0xa194, 0x0},
			want: "HelloHello",
		},
		{
			arg:  []rune{0x2182, 0x3518, 0xc204, 0xda14, 0xc800},
			want: "ababcabcdabcde",
		},
		{
			arg:  []rune("҅〶惶̀Ў㦀♀"),
			want: "Hello, world",
		},
		{
			arg:  []rune("邃ℐ摢ಁ⃈刌䀀"),
			want: "あいうえお",
		},
		{
			arg:  []rune{'輆', '', 0xda00},
			want: "🍎",
		},
		{
			arg:  []rune{0x8f06, 0xe397, 0xde9c, 0x5f68},
			want: "🍎🍇",
		},
		{
			arg:  []rune{0x21a2, 0x1064, 0x3c1b, 0x872f, 0xb046, 0x8220, 0xc68e, 0x2fb0, 0x6320},
			want: "aあ🍎bい🍇c",
		},
		{
			arg:  []rune{0xe50},
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
