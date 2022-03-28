package lzstring

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf16"
	"unicode/utf8"
)

func f(i int) rune {
	return rune(i)
}

var (
	ErrInvalidString = errors.New("Invalid string")
)

func Compress(uncompressed string) ([]rune, error) {
	if !utf8.ValidString(uncompressed) {
		return nil, ErrInvalidString
	}
	if len(uncompressed) == 0 {
		return []rune{}, nil
	}
	res, err := _compress(uncompressed, 16, func(i int) []rune {
		return []rune{rune(i)}
	})
	return res, err
}

type GetCharFunc func(i int) []rune

func _compress(uncompressed string, bitsPerChar int, getCharFromInt GetCharFunc) ([]rune, error) {
	var i, value int
	contextDictionary := make(map[string]int)
	contextDictionaryToCreate := make(map[string]bool)
	var contextC rune
	var contextWC, contextW []rune
	contextEnLargeIn := 2
	contextDictSize := 3
	contextNumBits := 2
	contextData := make([][]rune, 0)
	contextDataVal := 0
	contextDataPosition := 0
	var ii int
	uncompressedRune := utf16.Encode([]rune(uncompressed))
	for ii = 0; ii < len(uncompressedRune); ii++ {
		contextC = rune(uncompressedRune[ii])
		// contextW, contextWC are slice of runes, keys should be enclosed in brackets
		contextCKey := fmt.Sprintf("[%d]", contextC)
		if _, ok := contextDictionary[contextCKey]; !ok {
			contextDictionary[contextCKey] = contextDictSize
			contextDictSize++
			contextDictionaryToCreate[contextCKey] = true
		}
		contextWC = append(contextW, contextC)
		contextWCKey := fmt.Sprint(contextWC)
		contextWKey := fmt.Sprint(contextW)
		if _, ok := contextDictionary[contextWCKey]; ok {
			contextW = contextWC
		} else {
			if _, ok := contextDictionaryToCreate[contextWKey]; ok {
				if len(contextW) > 0 && contextW[0] < 256 {
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
					value = int(contextW[0])
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
					value = int(contextW[0])
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
				delete(contextDictionaryToCreate, contextWKey)
			} else {
				value = contextDictionary[contextWKey]
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
			contextDictionary[fmt.Sprint(contextWC)] = contextDictSize
			contextDictSize++
			contextW = []rune{contextC}
		}
	}

	if len(contextW) != 0 {
		contextWKey := fmt.Sprint(contextW)
		if _, ok := contextDictionaryToCreate[contextWKey]; ok {
			if contextW[0] < 256 {
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
				value = int(contextW[0])
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
				value = int(contextW[0])
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
			delete(contextDictionaryToCreate, contextWKey)
		} else {
			value = contextDictionary[contextWKey]
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
	result := make([]rune, 0)
	for _, cd := range contextData {
		result = append(result, cd...)
	}
	return result, nil
}

func Decompress(compressed []rune) string {
	if len(compressed) == 0 {
		return ""
	}
	compressedRune := compressed
	res := _decompress(len(compressedRune), 32768, func(index int) int {
		return int(compressedRune[index])
	})
	result := make([]uint16, 0)
	for _, r := range res {
		result = append(result, uint16(r))
	}
	return string(utf16.Decode(result))
}

type GetNextValFunc = func(index int) int

type Data struct {
	val      int
	position int
	index    int
}

func _decompress(length int, resetValue int, getNextVal GetNextValFunc) []rune {
	// for init
	dictionary := make(map[rune][]rune)
	var next int
	enlargeIn := 4
	dictSize := 4
	numBits := 3
	var entry []rune
	result := make([][]rune, 0)
	var i, bits, resb, maxpower, power int
	var c rune
	var w []rune
	data := Data{val: getNextVal(0), position: resetValue, index: 1}

	for i = 0; i < 3; i++ {
		dictionary[rune(i)] = []rune{rune(i)}
	}
	bits = 0
	maxpower = int(math.Pow(2, 2))
	power = 1
	for power != maxpower {
		resb = data.val & data.position
		data.position >>= 1
		if data.position == 0 {
			data.position = resetValue
			data.val = getNextVal(data.index)
			data.index += 1
		}
		tmp := 0
		if resb > 0 {
			tmp = 1
		}
		bits |= tmp * power
		power <<= 1
	}
	next = bits
	switch next {
	case 0:
		bits = 0
		maxpower = int(math.Pow(2, 8))
		power = 1
		for power != maxpower {
			resb = data.val & data.position
			data.position >>= 1
			if data.position == 0 {
				data.position = resetValue
				data.val = getNextVal(data.index)
				data.index += 1
			}
			tmp := 0
			if resb > 0 {
				tmp = 1
			}
			bits |= tmp * power
			power <<= 1
		}
		c = f(bits)
	case 1:
		bits = 0
		maxpower = int(math.Pow(2, 16))
		power = 1
		for power != maxpower {
			resb = data.val & data.position
			data.position >>= 1
			if data.position == 0 {
				data.position = resetValue
				data.val = getNextVal(data.index)
				data.index += 1
			}
			tmp := 0
			if resb > 0 {
				tmp = 1
			}
			bits |= tmp * power
			power <<= 1
		}
		c = f(bits)
		break
	case 2:
		return nil
	}
	dictionary[3] = []rune{c}
	w = []rune{c}
	result = append(result, []rune{c})
	for {
		if data.index > length {
			return []rune{}
		}

		bits = 0
		maxpower = int(math.Pow(2, float64(numBits)))
		power = 1
		for power != maxpower {
			resb = data.val & data.position
			data.position >>= 1
			if data.position == 0 {
				data.position = resetValue
				data.val = getNextVal(data.index)
				data.index += 1
			}
			tmp := 0
			if resb > 0 {
				tmp = 1
			}
			bits |= tmp * power
			power <<= 1
		}

		c = rune(bits)
		switch c {
		case 0:
			bits = 0
			maxpower = int(math.Pow(2, 8))
			power = 1
			for power != maxpower {
				resb = data.val & data.position
				data.position >>= 1
				if data.position == 0 {
					data.position = resetValue
					data.val = getNextVal(data.index)
					data.index++
				}
				tmp := 0
				if resb > 0 {
					tmp = 1
				}
				bits |= tmp * power
				power <<= 1
			}

			dictionary[rune(dictSize)] = []rune{f(bits)}
			dictSize++
			c = rune(dictSize - 1)
			enlargeIn--
		case 1:
			bits = 0
			maxpower = int(math.Pow(2, 16))
			power = 1
			for power != maxpower {
				resb = data.val & data.position
				data.position >>= 1
				if data.position == 0 {
					data.position = resetValue
					data.val = getNextVal(data.index)
					data.index++
				}
				tmp := 0
				if resb > 0 {
					tmp = 1
				}
				bits |= tmp * power
				power <<= 1
			}
			dictionary[rune(dictSize)] = []rune{f(bits)}
			dictSize++
			c = rune(dictSize - 1)
			enlargeIn--
		case 2:
			res := make([]rune, 0)
			for _, r := range result {
				res = append(res, r...)
			}
			return res
		}

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}

		if _, ok := dictionary[c]; ok {
			entry = append([]rune{}, dictionary[c]...)
		} else {
			if c == rune(dictSize) {
				entry = append(w[:0:0], w...)
				entry = append(entry, w[0])
			} else {
				return []rune{}
			}
		}
		result = append(result, []rune(entry))

		tmp := append(w[:0:0], w...)
		tmp = append(tmp, entry[0])
		dictionary[rune(dictSize)] = tmp
		dictSize++
		enlargeIn--

		w = entry

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}
	}
}
