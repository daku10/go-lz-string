package lzstring

import "io"

func Decompress() {

}

func Compress(reader io.Reader) (io.Reader, error) {
	return reader, nil
}

type GetCharFunc func(i int) string

func compressCore(uncompressed string, bitsPerChar int, getCharFromInt GetCharFunc) (string, error) {
	var i, value int
	contextDictionary := make(map[string]int)
	contextDictionaryToCreate := make(map[string]bool)
	var contextC, contextWC, contextW string
	contextEnLargeIn := 2
	contextDictSize := 3
	contextNumBits := 2
	contextData := make([]string, 0)
	contextDataVal := 0
	contextDataPosition := 0
	var ii int
	uncompressedRune := []rune(uncompressed)
	for ii := 0; ii < len(uncompressedRune); ii++ {
		contextC := string(uncompressedRune[ii])
		if _, ok := contextDictionary[contextC]; !ok {
			contextDictionary[contextC] = contextDictSize
			contextDictSize++
			contextDictionaryToCreate[contextC] = true
		}
		contextWC = contextW + contextC
		if _, ok := contextDictionary[contextWC]; ok {
			contextW = contextWC
		} else {

		}
	}
}
