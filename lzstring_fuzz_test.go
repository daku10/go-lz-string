package lzstring

import (
	"strings"
	"testing"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func FuzzIntegrityBase64(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToBase64(s)
		if !utf8.ValidString(s) {
			assert.NotNil(t, err)
			return
		}
		assert.Condition(t, func() (success bool) {
			for _, c := range compressed {
				if !strings.Contains(keyStrBase64, string(c)) {
					return false
				}
			}
			return true
		})
		repair, err := DecompressFromBase64(compressed)
		assert.Nil(t, err)
		assert.Equal(t, s, repair)
	})
}

func FuzzIntegrityUTF16(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToUTF16(s)
		if !utf8.ValidString(s) {
			assert.NotNil(t, err)
			return
		}
		isValid := utf8.ValidString(string(utf16.Decode(compressed)))
		assert.True(t, isValid)
		repair, err := DecompressFromUTF16(compressed)
		assert.Nil(t, err)
		assert.Equal(t, s, repair)
	})
}

func FuzzIntegrityUint8Array(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToUint8Array(s)
		if !utf8.ValidString(s) {
			assert.NotNil(t, err)
			return
		}
		repair, err := DecompressFromUint8Array(compressed)
		assert.Nil(t, err)
		assert.Equal(t, s, repair)
	})
}

func FuzzIntegrityEncodedURIComponent(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToEncodedURIComponent(s)
		if !utf8.ValidString(s) {
			assert.NotNil(t, err)
			return
		}
		if !utf8.ValidString(s) {
			assert.NotNil(t, err)
			return
		}
		assert.Condition(t, func() (success bool) {
			for _, c := range compressed {
				if !strings.Contains(keyStrUriSafe, string(c)) {
					return false
				}
			}
			return true
		})
		repair, err := DecompressFromEncodedURIComponent(compressed)
		assert.Nil(t, err)
		assert.Equal(t, s, repair)
	})
}
