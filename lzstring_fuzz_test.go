package lzstring

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/exec"
	"strings"
	"testing"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
)

func FuzzIntegrity(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := Compress(s)
		if !utf8.ValidString(s) {
			if err == nil {
				t.Fatalf("expected not nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		f, err := os.CreateTemp(t.TempDir(), "node")
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		// verify node output
		cmd := exec.Command("node", "testdata/test.js", "invalid-utf16", f.Name())
		cmd.Stdin = bytes.NewBufferString(s)
		se := bytes.Buffer{}
		cmd.Stderr = &se
		err = cmd.Run()
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		content, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		originalCompressed := make([]uint16, len(content)/2)
		err = binary.Read(bytes.NewReader(content), binary.LittleEndian, &originalCompressed)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		if diff := cmp.Diff(originalCompressed, compressed); diff != "" {
			t.Fatalf("got: %v want: %v diff: %v arg: %v", compressed, originalCompressed, diff, s)
		}

		repair, err := Decompress(compressed)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		if diff := cmp.Diff(s, repair); diff != "" {
			t.Errorf("got: %v want: %v diff: %v", repair, s, diff)
		}
	})
}

func FuzzIntegrityBase64(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToBase64(s)
		if !utf8.ValidString(s) {
			if err == nil {
				t.Fatalf("expected not nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		for _, c := range compressed {
			if !strings.Contains(keyStrBase64, string(c)) {
				t.Fatalf("expected only base64 characters, invalid character: %v got: %v", string(c), compressed)
			}
		}

		// verify node output
		cmd := exec.Command("node", "testdata/test.js", "base64")
		var stdout bytes.Buffer
		cmd.Stdin = bytes.NewBufferString(s)
		cmd.Stdout = &stdout
		err = cmd.Run()
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		originalCompressed := stdout.String()
		if diff := cmp.Diff(originalCompressed, compressed); diff != "" {
			t.Errorf("got: %v want: %v diff: %v arg: %v", originalCompressed, compressed, diff, s)
		}

		repair, err := DecompressFromBase64(compressed)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		if diff := cmp.Diff(s, repair); diff != "" {
			t.Errorf("got: %v want: %v diff: %v", repair, s, diff)
		}
	})
}

func FuzzIntegrityUTF16(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToUTF16(s)
		if !utf8.ValidString(s) {
			if err == nil {
				t.Fatalf("expected not nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		isValid := utf8.ValidString(string(utf16.Decode(compressed)))
		if isValid != true {
			t.Fatalf("expected true, got: false arg: %v", compressed)
		}
		repair, err := DecompressFromUTF16(compressed)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		if diff := cmp.Diff(s, repair); diff != "" {
			t.Errorf("got: %v want: %v diff: %v", repair, s, diff)
		}
	})
}

func FuzzIntegrityUint8Array(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToUint8Array(s)
		if !utf8.ValidString(s) {
			if err == nil {
				t.Fatalf("expected not nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		// verify node output
		cmd := exec.Command("node", "testdata/test.js", "uint8array")
		var stdout bytes.Buffer
		cmd.Stdin = bytes.NewBufferString(s)
		cmd.Stdout = &stdout
		err = cmd.Run()
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		originalCompressed := stdout.Bytes()
		if diff := cmp.Diff(originalCompressed, compressed); diff != "" {
			t.Errorf("got: %v want: %v diff: %v arg: %v", originalCompressed, compressed, diff, s)
		}

		repair, err := DecompressFromUint8Array(compressed)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		if diff := cmp.Diff(s, repair); diff != "" {
			t.Errorf("got: %v want: %v diff: %v", repair, s, diff)
		}
	})
}

func FuzzIntegrityEncodedURIComponent(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		compressed, err := CompressToEncodedURIComponent(s)
		if !utf8.ValidString(s) {
			if err == nil {
				t.Fatalf("expected not nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		for _, c := range compressed {
			if !strings.Contains(keyStrUriSafe, string(c)) {
				t.Fatalf("expected only uri safe characters, invalid character: %v got: %v", string(c), compressed)
			}
		}

		// verify node output
		cmd := exec.Command("node", "testdata/test.js", "encodedURIComponent")
		var stdout bytes.Buffer
		cmd.Stdin = bytes.NewBufferString(s)
		cmd.Stdout = &stdout
		err = cmd.Run()
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
		originalCompressed := stdout.String()
		if diff := cmp.Diff(originalCompressed, compressed); diff != "" {
			t.Errorf("got: %v want: %v diff: %v arg: %v", originalCompressed, compressed, diff, s)
		}

		repair, err := DecompressFromEncodedURIComponent(compressed)
		if err != nil {
			t.Fatalf("expected not nil, got: %v", err)
		}
		if diff := cmp.Diff(s, repair); diff != "" {
			t.Errorf("got: %v want: %v diff: %v", repair, s, diff)
		}
	})
}
