package lzstring

import (
	"io"
	"math"
	"strings"
	"unicode/utf16"
)

func Decompress() {

}

func Compress(reader io.Reader) (io.Reader, error) {
	s, _ := io.ReadAll(reader)
	res, err := compressCore(string(s), 16, func(i int) string {
		return string(utf16.Decode([]uint16{uint16(i)}))
	})
	return strings.NewReader(res), err
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
	uncompressedRune := utf16.Encode([]rune(uncompressed))
	for ii = 0; ii < len(uncompressedRune); ii++ {
		contextC = string(uncompressedRune[ii])
		if _, ok := contextDictionary[contextC]; !ok {
			contextDictionary[contextC] = contextDictSize
			contextDictSize++
			contextDictionaryToCreate[contextC] = true
		}
		contextWC = contextW + contextC
		if _, ok := contextDictionary[contextWC]; ok {
			contextW = contextWC
		} else {
			if _, ok := contextDictionaryToCreate[contextW]; ok {
				if charCodeAtZero(contextW) < 256 {
					for i = 0; i < contextNumBits; i++ {
						contextDataVal = contextDataVal << 1
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
					}
					value = int(charCodeAtZero(contextW))
					for i = 0; i < 8; i++ {
						contextDataVal = (contextDataVal << 1) | (value & 1)
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
						value = value >> 1
					}
				} else {
					value = 1
					for i = 0; i < contextNumBits; i++ {
						contextDataVal = (contextDataVal << 1) | value
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
						value = 0
					}
					value = int(charCodeAtZero(contextW))
					for i = 0; i < 16; i++ {
						contextDataVal = (contextDataVal << 1) | (value & 1)
						if contextDataPosition == bitsPerChar-1 {
							contextDataPosition = 0
							contextData = append(contextData, getCharFromInt(contextDataVal))
							contextDataVal = 0
						} else {
							contextDataPosition++
						}
						value = value >> 1
					}
				}
				contextEnLargeIn--
				if contextEnLargeIn == 0 {
					contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
					contextNumBits++
				}
				delete(contextDictionaryToCreate, contextW)
			} else {
				value = contextDictionary[contextW]
				for i = 0; i < contextNumBits; i++ {
					contextDataVal = (contextDataVal << 1) | (value & 1)
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = value >> 1
				}
			}
			contextEnLargeIn--
			if contextEnLargeIn == 0 {
				contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
				contextNumBits++
			}
			contextDictionary[contextWC] = contextDictSize
			contextDictSize++
			contextW = contextC
		}
	}

	if contextW != "" {
		if _, ok := contextDictionaryToCreate[contextW]; ok {
			if charCodeAtZero(contextW) < 256 {
				for i = 0; i < contextNumBits; i++ {
					contextDataVal = contextDataVal << 1
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
				}
				value = int(charCodeAtZero(contextW))
				for i = 0; i < 8; i++ {
					contextDataVal = (contextDataVal << 1) | (value & 1)
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = value >> 1
				}
			} else {
				value = 1
				for i = 0; i < contextNumBits; i++ {
					contextDataVal = (contextDataVal << 1) | value
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = 0
				}
				value = int(charCodeAtZero(contextW))
				for i = 0; i < 16; i++ {
					contextDataVal = (contextDataVal << 1) | (value & 1)
					if contextDataPosition == bitsPerChar-1 {
						contextDataPosition = 0
						contextData = append(contextData, getCharFromInt(contextDataVal))
						contextDataVal = 0
					} else {
						contextDataPosition++
					}
					value = value >> 1
				}
			}
			contextEnLargeIn--
			if contextEnLargeIn == 0 {
				contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
				contextNumBits++
			}
			delete(contextDictionaryToCreate, contextW)
		} else {
			value = contextDictionary[contextW]
			for i = 0; i < contextNumBits; i++ {
				contextDataVal = (contextDataVal << 1) | (value & 1)
				if contextDataPosition == bitsPerChar-1 {
					contextDataPosition = 0
					contextData = append(contextData, getCharFromInt(contextDataVal))
					contextDataVal = 0
				} else {
					contextDataPosition++
				}
				value = value >> 1
			}
		}
		contextEnLargeIn--
		if contextEnLargeIn == 0 {
			contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
			contextNumBits++
		}
	}

	value = 2
	for i = 0; i < contextNumBits; i++ {
		contextDataVal = (contextDataVal << 1) | (value & 1)
		if contextDataPosition == bitsPerChar-1 {
			contextDataPosition = 0
			contextData = append(contextData, getCharFromInt(contextDataVal))
			contextDataVal = 0
		} else {
			contextDataPosition++
		}
		value = value >> 1
	}

	for {
		contextDataVal = contextDataVal << 1
		if contextDataPosition == bitsPerChar-1 {
			contextData = append(contextData, getCharFromInt(contextDataVal))
			break
		} else {
			contextDataPosition++
		}
	}
	return strings.Join(contextData, ""), nil
}

func charCodeAtZero(s string) rune {
	r := ([]rune(s))[0]
	r1, _ := utf16.EncodeRune(r)
	if r1 == '\uFFFD' {
		return r
	}
	return r1
}
