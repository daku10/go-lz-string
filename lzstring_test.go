package lzstring

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf16"
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
			want: []uint16{0x4000},
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
		{
			arg:  "邊󠄆",
			want: []uint16{0x9442, 0x6016, 0xdc60, 0xbb40},
		},
		{
			arg:  "👨‍👨‍👦",
			want: []uint16{0xaf06, 0xe0b1, 0xdcb0, 0x4a1, 0x8663, 0xb400},
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
		{
			arg:  "邊󠄆",
			want: []uint16{0x4a41, 0x1825, 0x5bac, 0xbd4, 0x20, 0x20},
		},
		{
			arg:  "👨‍👨‍👦",
			want: []uint16{0x57a3, 0x384c, 0x3bb6, 0x6a, 0xc53, 0xef0, 0x20},
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

func TestCompressToUint8Array(t *testing.T) {
	tests := []struct {
		arg  string
		want []byte
	}{
		{
			arg:  "",
			want: []byte{0x40, 0x0},
		},
		{
			arg:  "H",
			want: []byte{0x4, 0x90},
		},
		{
			arg:  "HelloHello",
			want: []byte{0x4, 0x85, 0x30, 0x36, 0x60, 0xf6, 0xa1, 0x94, 0x0, 0x0},
		},
		{
			arg:  "ababcabcdabcde",
			want: []byte{0x21, 0x82, 0x35, 0x18, 0xc2, 0x4, 0xda, 0x14, 0xc8, 0x0},
		},
		{
			arg:  "Hello, world",
			want: []byte{0x4, 0x85, 0x30, 0x36, 0x60, 0xf6, 0x3, 0x40, 0x4, 0xe, 0xe9, 0x1, 0x39, 0x80, 0x26, 0x40},
		},
		{
			arg:  "あいうえお",
			want: []byte{0x90, 0x83, 0x21, 0x10, 0x64, 0x62, 0xc, 0x81, 0x20, 0xc8, 0x52, 0xc, 0x40, 0x0},
		},
		{
			arg:  "🍎",
			want: []byte{0x8f, 0x6, 0xe3, 0x97, 0xda, 0x0},
		},
		{
			arg:  "🍎🍇",
			want: []byte{0x8f, 0x6, 0xe3, 0x97, 0xde, 0x9c, 0x5f, 0x68},
		},
		{
			arg:  "aあ🍎bい🍇c",
			want: []byte{0x21, 0xa2, 0x10, 0x64, 0x3c, 0x1b, 0x87, 0x2f, 0xb0, 0x46, 0x82, 0x20, 0xc6, 0x8e, 0x2f, 0xb0, 0x63, 0x20},
		},
		{
			arg:  string([]rune{0x9c}),
			want: []byte{0xe, 0x50},
		},
		{
			arg:  "邊󠄆",
			want: []byte{0x94, 0x42, 0x60, 0x16, 0xdc, 0x60, 0xbb, 0x40},
		},
		{
			arg:  "👨‍👨‍👦",
			want: []byte{0xaf, 0x6, 0xe0, 0xb1, 0xdc, 0xb0, 0x4, 0xa1, 0x86, 0x63, 0xb4, 0x0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := CompressToUint8Array(tt.arg)
			assert.Nil(t, err)
			b, err := reader, err
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestCompressToEncodedURIComponent(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{
			arg:  "",
			want: "Q",
		},
		{
			arg:  "H",
			want: "BJA",
		},
		{
			arg:  "HelloHello",
			want: "BIUwNmD2oZQ",
		},
		{
			arg:  "ababcabcdabcde",
			want: "IYI1GMIE2hTI",
		},
		{
			arg:  "Hello, world",
			want: "BIUwNmD2A0AEDukBOYAmQ",
		},
		{
			arg:  "あいうえお",
			want: "kIMhEGRiDIEgyFIMQ",
		},
		{
			arg:  "🍎",
			want: "jwbjl9o",
		},
		{
			arg:  "🍎🍇",
			want: "jwbjl96cX2g",
		},
		{
			arg:  "aあ🍎bい🍇c",
			want: "IaIQZDwbhy+wRoIgxo4vsGMg",
		},
		{
			arg:  string([]rune{0x9c}),
			want: "DlA",
		},
		{
			arg:  "邊󠄆",
			want: "lEJgFtxgu0A",
		},
		{
			arg:  "👨‍👨‍👦",
			want: "rwbgsdywBKGGY7Q",
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			reader, err := CompressToEncodedURIComponent(tt.arg)
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
		{
			arg:  []uint16{0x9442, 0x6016, 0xdc60, 0xbb40},
			want: "邊󠄆",
		},
		{
			arg:  []uint16{0xaf06, 0xe0b1, 0xdcb0, 0x4a1, 0x8663, 0xb400},
			want: "👨‍👨‍👦",
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

func TestDecompressFromUint8Array(t *testing.T) {
	tests := []struct {
		arg  []byte
		want string
	}{
		{
			arg:  []byte{0x40, 0x0},
			want: "",
		},
		{
			arg:  []byte{0x4, 0x90},
			want: "H",
		},
		{
			arg:  []byte{0x4, 0x85, 0x30, 0x36, 0x60, 0xf6, 0xa1, 0x94, 0x0, 0x0},
			want: "HelloHello",
		},
		{
			arg:  []byte{0x21, 0x82, 0x35, 0x18, 0xc2, 0x4, 0xda, 0x14, 0xc8, 0x0},
			want: "ababcabcdabcde",
		},
		{
			arg:  []byte{0x4, 0x85, 0x30, 0x36, 0x60, 0xf6, 0x3, 0x40, 0x4, 0xe, 0xe9, 0x1, 0x39, 0x80, 0x26, 0x40},
			want: "Hello, world",
		},
		{
			arg:  []byte{0x90, 0x83, 0x21, 0x10, 0x64, 0x62, 0xc, 0x81, 0x20, 0xc8, 0x52, 0xc, 0x40, 0x0},
			want: "あいうえお",
		},
		{
			arg:  []byte{0x8f, 0x6, 0xe3, 0x97, 0xda, 0x0},
			want: "🍎",
		},
		{
			arg:  []byte{0x8f, 0x6, 0xe3, 0x97, 0xde, 0x9c, 0x5f, 0x68},
			want: "🍎🍇",
		},
		{
			arg:  []byte{0x21, 0xa2, 0x10, 0x64, 0x3c, 0x1b, 0x87, 0x2f, 0xb0, 0x46, 0x82, 0x20, 0xc6, 0x8e, 0x2f, 0xb0, 0x63, 0x20},
			want: "aあ🍎bい🍇c",
		},
		{
			arg:  []byte{0xe, 0x50},
			want: string([]rune{0x9c}),
		},
		{
			arg:  []byte{0x94, 0x42, 0x60, 0x16, 0xdc, 0x60, 0xbb, 0x40},
			want: "邊󠄆",
		},
		{
			arg:  []byte{0xaf, 0x6, 0xe0, 0xb1, 0xdc, 0xb0, 0x4, 0xa1, 0x86, 0x63, 0xb4, 0x0},
			want: "👨‍👨‍👦",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.arg), func(t *testing.T) {
			res, err := DecompressFromUint8Array(tt.arg)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestDecompressFromEncodedURIComponent(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{
			arg:  "Q",
			want: "",
		},
		{
			arg:  "BJA",
			want: "H",
		},
		{
			arg:  "BIUwNmD2oZQ",
			want: "HelloHello",
		},
		{
			arg:  "IYI1GMIE2hTI",
			want: "ababcabcdabcde",
		},
		{
			arg:  "BIUwNmD2A0AEDukBOYAmQ",
			want: "Hello, world",
		},
		{
			arg:  "kIMhEGRiDIEgyFIMQ",
			want: "あいうえお",
		},
		{
			arg:  "jwbjl9o",
			want: "🍎",
		},
		{
			arg:  "jwbjl96cX2g",
			want: "🍎🍇",
		},
		{
			arg:  "IaIQZDwbhy+wRoIgxo4vsGMg",
			want: "aあ🍎bい🍇c",
		},
		{
			arg:  "DlA",
			want: string([]rune{0x9c}),
		},
		{
			arg:  "lEJgFtxgu0A",
			want: "邊󠄆",
		},
		{
			arg:  "rwbgsdywBKGGY7Q",
			want: "👨‍👨‍👦",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.arg), func(t *testing.T) {
			res, err := DecompressFromEncodedURIComponent(tt.arg)
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
