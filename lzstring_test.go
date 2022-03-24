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
