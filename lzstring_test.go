package lzstring

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{
			arg:  "H",
			want: "Ґ",
		},
		{
			arg:  "HelloHello",
			want: string([]rune{0x485, 0x3036, 0x60f6, 0xa194, 0x0}),
		},
		{
			arg:  "ababcabcdabcde",
			want: string([]rune{0x2182, 0x3518, 0xc204, 0xda14, 0xc800}),
		},
		{
			arg:  "Hello, world",
			want: "҅〶惶̀Ў㦀♀",
		},
		{
			arg:  "あいうえお",
			want: "邃ℐ摢ಁ⃈刌䀀",
		},
		{
			arg:  "🍎",
			want: string([]rune{'輆', '', 0xda00}),
		},
		{
			arg:  "🍎🍇",
			want: string([]rune{0x8f06, 0xe397, 0xde9c, 0x5f68}),
		},
		{
			arg:  "aあ🍎bい🍇c",
			want: string([]rune{0x21a2, 0x1064, 0x3c1b, 0x872f, 0xb046, 0x8220, 0xc68e, 0x2fb0, 0x6320}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := Compress(strings.NewReader(tt.arg))
			assert.Nil(t, err)
			b, err := reader, err
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestDecompress(t *testing.T) {

}
