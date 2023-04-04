package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			tmpOut, err := os.CreateTemp(t.TempDir(), "out")

			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			config := &Config{
				In:  os.Stdin,
				Out: os.Stdout,
				Err: os.Stderr,
			}
			cmd := newDecompressCmd(config)
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
