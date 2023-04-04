package lzstring

import (
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

		repair, err := DecompressFromEncodedURIComponent(compressed)
		if err != nil {
			t.Fatalf("expected not nil, got: %v", err)
		}
		if diff := cmp.Diff(s, repair); diff != "" {
			t.Errorf("got: %v want: %v diff: %v", repair, s, diff)
		}
	})
}
