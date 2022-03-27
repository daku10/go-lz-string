package lzstring

import (
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf16"
)

func f(i int) string {
	return string([]rune{rune(i)})
}

func Compress(reader io.Reader) (string, error) {
	s, _ := io.ReadAll(reader)
	res, err := _compress(string(s), 16, func(i int) string {
		return string([]rune{rune(i)})
	})
	return res, err
}

type GetCharFunc func(i int) string

func _compress(uncompressed string, bitsPerChar int, getCharFromInt GetCharFunc) (string, error) {
	var i, value int
	contextDictionary := make(map[string]int)
	contextDictionaryToCreate := make(map[string]bool)
	var contextC rune
	var contextWC, contextW []rune
	contextEnLargeIn := 2
	contextDictSize := 3
	contextNumBits := 2
	contextData := make([]string, 0)
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
	return strings.Join(contextData, ""), nil
}

func Decompress(compressed string) string {
	if compressed == "" {
		return ""
	}
	compressedRune := utf16.Encode([]rune(compressed))
	return _decompress(len(compressedRune), 32768, func(index int) int {
		return int(compressedRune[index])
	})
}

type GetNextValFunc = func(index int) int

type Data struct {
	val      int
	position int
	index    int
}

func _decompress(length int, resetValue int, getNextVal GetNextValFunc) string {
	// for init
	dictionary := make(map[any]any)
	var next int
	enlargeIn := 4
	dictSize := 4
	numBits := 3
	entry := ""
	result := make([]string, 0)
	var i, bits, resb, maxpower, power int
	var w, c interface{}
	data := Data{val: getNextVal(0), position: resetValue, index: 1}

	for i = 0; i < 3; i++ {
		dictionary[i] = i
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
		return ""
	}
	dictionary[3] = c
	w = c
	result = append(result, c.(string))
	for {
		if data.index > length {
			return ""
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

		c = bits
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

			dictionary[dictSize] = f(bits)
			dictSize++
			c = dictSize - 1
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
			dictionary[dictSize] = f(bits)
			dictSize++
			c = dictSize - 1
			enlargeIn--
		case 2:
			return strings.Join(result, "")
		}

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}

		if _, ok := dictionary[c]; ok {
			entry = dictionary[c].(string)
		} else {
			if c == dictSize {
				entry = w.(string) + string([]rune(w.(string))[0])
			} else {
				return ""
			}
		}
		result = append(result, entry)

		dictionary[dictSize] = w.(string) + string([]rune(entry)[0])
		dictSize++
		enlargeIn--

		w = entry

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}
	}
}
