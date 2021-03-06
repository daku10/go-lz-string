package lzstring

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func f(i int) uint16 {
	return uint16(i)
}

var (
	ErrInputInvalidString = errors.New("Input is invalid string")
	ErrInputNotDecodable  = errors.New("Input is not decodable")
	ErrInputNil           = errors.New("Input should not be nil")
	ErrInputBlank         = errors.New("Input should not be blank")
)

func Compress(uncompressed string) ([]uint16, error) {
	if !utf8.ValidString(uncompressed) {
		return nil, ErrInputInvalidString
	}
	res, err := _compress(uncompressed, 16, func(i int) []uint16 {
		return []uint16{uint16(i)}
	})
	return res, err
}

const keyStrBase64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="

func CompressToBase64(uncompressed string) (string, error) {
	if !utf8.ValidString(uncompressed) {
		return "", ErrInputInvalidString
	}
	res, err := _compress(uncompressed, 6, func(i int) []uint16 {
		return []uint16{uint16(keyStrBase64[i])}
	})
	if err != nil {
		return "", err
	}
	resStr := string(utf16.Decode(res))
	switch len(resStr) % 4 {
	case 0:
		return resStr, nil
	case 1:
		return resStr + "===", nil
	case 2:
		return resStr + "==", nil
	case 3:
		return resStr + "=", nil
	default:
		return resStr, nil
	}
}

func CompressToUTF16(uncompressed string) ([]uint16, error) {
	if !utf8.ValidString(uncompressed) {
		return nil, ErrInputInvalidString
	}
	res, err := _compress(uncompressed, 15, func(i int) []uint16 {
		return []uint16{f(i + 32)}
	})
	if err != nil {
		return nil, err
	}
	// 32 means " "(space) character
	res = append(res, 32)
	return res, nil
}

func CompressToUint8Array(uncompressed string) ([]byte, error) {
	if !utf8.ValidString(uncompressed) {
		return nil, ErrInputInvalidString
	}
	res, err := Compress(uncompressed)
	if err != nil {
		return nil, err
	}
	length := len(res)
	buf := make([]byte, length*2)
	for i := 0; i < length; i++ {
		currentValue := res[i]
		buf[i*2] = byte(currentValue >> 8)
		buf[i*2+1] = byte(currentValue % 256)
	}
	return buf, nil
}

const keyStrUriSafe = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+-$"

func CompressToEncodedURIComponent(uncompressed string) (string, error) {
	if !utf8.ValidString(uncompressed) {
		return "", ErrInputInvalidString
	}
	res, err := _compress(uncompressed, 6, func(i int) []uint16 {
		return []uint16{uint16(keyStrUriSafe[i])}
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

type getCharFunc func(i int) []uint16

func _compress(uncompressed string, bitsPerChar int, getCharFromInt getCharFunc) ([]uint16, error) {
	var i, value int
	contextDictionary := make(map[string]int)
	contextDictionaryToCreate := make(map[string]bool)
	var contextC uint16
	var contextWC, contextW []uint16
	contextEnLargeIn := 2
	contextDictSize := 3
	contextNumBits := 2
	contextData := make([][]uint16, 0)
	contextDataVal := 0
	contextDataPosition := 0
	var ii int
	uncompressedRune := utf16.Encode([]rune(uncompressed))
	for ii = 0; ii < len(uncompressedRune); ii++ {
		contextC = uncompressedRune[ii]
		// contextW, contextWC are slice of runes, keys should be enclosed in brackets
		contextCKey := fmt.Sprintf("[%d]", contextC)
		if _, ok := contextDictionary[contextCKey]; !ok {
			contextDictionary[contextCKey] = contextDictSize
			contextDictSize++
			contextDictionaryToCreate[contextCKey] = true
		}
		contextWC = make([]uint16, len(contextW))
		copy(contextWC, contextW)
		contextWC = append(contextWC, contextC)
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
			contextW = []uint16{contextC}
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
	result := make([]uint16, 0)
	for _, cd := range contextData {
		result = append(result, cd...)
	}
	return result, nil
}

func Decompress(compressed []uint16) (string, error) {
	if compressed == nil {
		return "", ErrInputNil
	}
	if len(compressed) == 0 {
		return "", nil
	}
	res, err := _decompress(len(compressed), 32768, func(index int) int {
		return int(compressed[index])
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

func DecompressFromBase64(compressed string) (string, error) {
	if compressed == "" {
		return "", ErrInputBlank
	}
	res, err := _decompress(len(compressed), 32, func(index int) int {
		return getBaseValue(keyStrBase64, compressed[index])
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

var baseReverseDic map[string]map[byte]int = make(map[string]map[byte]int)

func getBaseValue(alphabet string, character byte) int {
	if _, ok := baseReverseDic[alphabet]; !ok {
		baseReverseDic[alphabet] = make(map[byte]int)
		for i := 0; i < len(alphabet); i++ {
			baseReverseDic[alphabet][alphabet[i]] = i
		}
	}
	return baseReverseDic[alphabet][character]
}

type getNextValFunc = func(index int) int

func DecompressFromUTF16(compressed []uint16) (string, error) {
	if compressed == nil {
		return "", ErrInputNil
	}
	if len(compressed) == 0 {
		return "", ErrInputBlank
	}
	res, err := _decompress(len(compressed), 16384, func(index int) int {
		return int(compressed[index] - 32)
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

func DecompressFromUint8Array(compressed []byte) (string, error) {
	if compressed == nil {
		return "", ErrInputNil
	}
	length := len(compressed) / 2
	buf := make([]uint16, len(compressed)/2)
	for i := 0; i < length; i++ {
		buf[i] = uint16(compressed[i*2])*256 + uint16(compressed[i*2+1])
	}

	return Decompress(buf)
}

func DecompressFromEncodedURIComponent(compressed string) (string, error) {
	replaced := strings.Replace(compressed, " ", "+", -1)
	res, err := _decompress(len(replaced), 32, func(index int) int {
		return getBaseValue(keyStrUriSafe, replaced[index])
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

type data struct {
	val      int
	position int
	index    int
}

func _decompress(length int, resetValue int, getNextVal getNextValFunc) ([]uint16, error) {
	// for init
	dictionary := make(map[uint16][]uint16)
	var next int
	enlargeIn := 4
	dictSize := 4
	numBits := 3
	var entry []uint16
	result := make([][]uint16, 0)
	var i uint16
	var bits, resb, maxpower, power int
	var c uint16
	var w []uint16
	data := data{val: getNextVal(0), position: resetValue, index: 1}

	for i = 0; i < 3; i++ {
		dictionary[i] = []uint16{i}
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
		return nil, nil
	}
	dictionary[3] = []uint16{c}
	w = []uint16{c}
	result = append(result, []uint16{c})
	for {
		if data.index > length {
			return []uint16{}, nil
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

		c = f(bits)
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

			dictionary[uint16(dictSize)] = []uint16{f(bits)}
			dictSize++
			c = uint16(dictSize - 1)
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
			dictionary[uint16(dictSize)] = []uint16{f(bits)}
			dictSize++
			c = uint16(dictSize - 1)
			enlargeIn--
		case 2:
			res := make([]uint16, 0)
			for _, r := range result {
				res = append(res, r...)
			}
			return res, nil
		}

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}

		if _, ok := dictionary[c]; ok {
			entry = make([]uint16, len(dictionary[c]))
			copy(entry, dictionary[c])
		} else {
			if c == uint16(dictSize) {
				entry = make([]uint16, len(w))
				copy(entry, w)
				entry = append(entry, w[0])
			} else {
				return nil, ErrInputNotDecodable
			}
		}
		result = append(result, entry)

		tmp := make([]uint16, len(w))
		copy(tmp, w)
		tmp = append(tmp, entry[0])
		dictionary[uint16(dictSize)] = tmp
		dictSize++
		enlargeIn--

		w = entry

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}
	}
}
