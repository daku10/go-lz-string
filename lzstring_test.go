package lzstring

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	t.Run("", func(t *testing.T) {
		reader, err := Compress(strings.NewReader("Hello, World"))
		assert.NotNil(t, err)
		b, err := io.ReadAll(reader)
		assert.Equal(t, b, []byte{0x85, 0x04, 0x36, 0x30, 0xf6, 0x60, 0x40, 0x03, 0x0e, 0x04, 0x01, 0xa9, 0x80, 0x39, 0x40, 0x26})
	})
}

func TestDecompress(t *testing.T) {

}
