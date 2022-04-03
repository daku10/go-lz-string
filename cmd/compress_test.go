package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
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
			tmpOut, err := os.CreateTemp(t.TempDir(), "out")
			assert.Nil(t, err)
			tmpIn, err := os.CreateTemp(t.TempDir(), "in")
			assert.Nil(t, err)
			config := &Config{
				In:  os.Stdin,
				Out: os.Stdout,
				Err: os.Stderr,
			}
			cmd := newCompressCmd(config)
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
