package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			tmpIn, err := os.CreateTemp(t.TempDir(), "in")
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			config := &Config{
				In:  os.Stdin,
				Out: os.Stdout,
				Err: os.Stderr,
			}
			cmd := newCompressCmd(config)
			cmd.SetArgs([]string{"-m", "base64", "-o", tmpOut.Name(), tmpIn.Name()})
			err = ioutil.WriteFile(tmpIn.Name(), []byte(tt.arg), os.ModePerm)
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			err = cmd.Execute()
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			res, err := ioutil.ReadFile(tmpOut.Name())
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			if diff := cmp.Diff(tt.want, string(res)); diff != "" {
				t.Errorf("got: %v want: %v diff: %v", string(res), tt.want, diff)
			}
		})
	}
}
