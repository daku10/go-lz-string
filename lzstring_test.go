package lzstring

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		arg  string
		want []byte
	}{
		{
			arg:  "Hello, world",
			want: []byte{0x85, 0x04, 0x36, 0x30, 0xf6, 0x60, 0x40, 0x03, 0x0e, 0x04, 0x01, 0xe9, 0x80, 0x39, 0x40, 0x26},
		},
		{
			arg:  "あいうえお",
			want: []byte{0x83, 0x90, 0x10, 0x21, 0x62, 0x64, 0x81, 0x0c, 0xc8, 0x20, 0x0c, 0x52, 0x00, 0x40},
		},
		{
			arg:  "🍎",
			want: []byte{0x06, 0x8f, 0x97, 0xe3, 0x00, 0xda},
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := Compress(strings.NewReader(tt.arg))
			assert.NotNil(t, err)
			b, err := io.ReadAll(reader)
			assert.Equal(t, b, tt.want)
		})
	}
}

func TestDecompress(t *testing.T) {

}
