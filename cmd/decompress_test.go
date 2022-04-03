package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompressCmd(t *testing.T) {
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
			tmpIn, err := os.CreateTemp(t.TempDir(), "in")
			assert.Nil(t, err)
			tmpOut, err := os.CreateTemp(t.TempDir(), "out")
			assert.Nil(t, err)
			cmd := newDecompressCmd()
			cmd.SetArgs([]string{"-m", "base64", "-o", tmpOut.Name(), tmpIn.Name()})
			ioutil.WriteFile(tmpIn.Name(), []byte(tt.arg), os.ModePerm)
			err = cmd.Execute()
			assert.Nil(t, err)
			res, err := ioutil.ReadFile(tmpOut.Name())
			assert.Nil(t, err)
			assert.Equal(t, tt.want, string(res))
		})
	}
}
