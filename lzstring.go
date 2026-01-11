// Package lzstring implements the LZ-String algorithm for string compression
// and decompression. The library features two main sets of functions,
// Compress and Decompress, which are used to compress and decompress strings,
// respectively.
package lzstring

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func f(i int) uint16 {
	return uint16(i)
}

var (
	ErrInputInvalidString = errors.New("input is invalid string")
	ErrInputNotDecodable  = errors.New("input is not decodable")
	ErrInputNil           = errors.New("input should not be nil")
	ErrInputBlank         = errors.New("input should not be blank")
)

// Compress takes an uncompressed string and compresses it into a slice of uint16.
// It returns an error if the input string is not a valid UTF-8 string.
// Note: The resulting uint16 slice may contain invalid UTF-16 characters,
// which is consistent with the original algorithm's behavior.
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

// CompressToBase64 takes an uncompressed string and compresses it into a Base64 string.
// It returns an error if the input string is not a valid UTF-8 string.
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

// CompressToUTF16 takes an uncompressed string and compresses it into a slice of uint16,
// where each element represents a UTF-16 encoded character.
// It returns an error if the input string is not a valid UTF-8 string.
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

// CompressToUint8Array takes an uncompressed string and compresses it into a slice of bytes.
// It returns an error if the input string is not a valid UTF-8 string.
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

// CompressToEncodedURIComponent takes an uncompressed string and compresses it into
// a URL-safe string, where special characters are replaced with safe alternatives.
// It returns an error if the input string is not a valid UTF-8 string.
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

// make consistency with slice of uint16 to be enclosed with bracket.
func uint16ToString(x uint16) string {
	var b bytes.Buffer
	b.WriteByte('[')
	b.WriteString(strconv.Itoa(int(x)))
	b.WriteByte(']')
	return b.String()
}

func uint16sToString(xs []uint16) string {
	var b bytes.Buffer
	b.WriteByte('[')
	for i, x := range xs {
		b.WriteString(strconv.Itoa(int(x)))
		if i != len(xs)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')
	return b.String()
}

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
		contextCKey := uint16ToString(contextC)
		if _, ok := contextDictionary[contextCKey]; !ok {
			contextDictionary[contextCKey] = contextDictSize
			contextDictSize++
			contextDictionaryToCreate[contextCKey] = true
		}
		contextWC = make([]uint16, len(contextW))
		copy(contextWC, contextW)
		contextWC = append(contextWC, contextC)
		contextWCKey := uint16sToString(contextWC)
		contextWKey := uint16sToString(contextW)
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
					contextEnLargeIn = 1 << contextNumBits
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
				contextEnLargeIn = 1 << contextNumBits
				contextNumBits++
			}
			contextDictionary[uint16sToString(contextWC)] = contextDictSize
			contextDictSize++
			contextW = []uint16{contextC}
		}
	}
	if len(contextW) != 0 {
		contextWKey := uint16sToString(contextW)
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
				contextEnLargeIn = 1 << contextNumBits
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
			// original algorithm has below expression, but this value is unused probably.
			// contextEnLargeIn = 1 << contextNumBits
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

// Decompress takes a compressed slice of uint16 main contain invalid UTF-16 characters and decompresses it into a string.
// It returns an error if the input is not a valid compressed data.
func Decompress(compressed []uint16) (string, error) {
	if compressed == nil {
		return "", ErrInputNil
	}
	if len(compressed) == 0 {
		return "", ErrInputBlank
	}
	res, err := _decompress(len(compressed), 32768, func(index int) (int, error) {
		if index >= len(compressed) {
			// Match JavaScript behavior: out-of-bounds reads return undefined,
			// which becomes 0 in bitwise operations.
			return 0, nil
		}
		return int(compressed[index]), nil
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

// DecompressFromBase64 takes a compressed Base64 string and decompresses it into a string.
// It returns an error if the input is not a valid compressed data.
func DecompressFromBase64(compressed string) (string, error) {
	if compressed == "" {
		return "", ErrInputBlank
	}
	res, err := _decompress(len(compressed), 32, func(index int) (int, error) {
		if index >= len(compressed) {
			// Match JavaScript behavior: out-of-bounds reads return undefined,
			// which becomes 0 in bitwise operations.
			return 0, nil
		}
		return getBaseValue(keyStrBase64, compressed[index]), nil
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

type getNextValFunc = func(index int) (int, error)

// DecompressFromUTF16 takes a compressed slice of uint16 UTF-16 characters and decompresses it into a string.
// It returns an error if the input is not a valid compressed data.
func DecompressFromUTF16(compressed []uint16) (string, error) {
	if compressed == nil {
		return "", ErrInputNil
	}
	if len(compressed) == 0 {
		return "", ErrInputBlank
	}
	res, err := _decompress(len(compressed), 16384, func(index int) (int, error) {
		if index >= len(compressed) {
			// Match JavaScript behavior: out-of-bounds reads return undefined,
			// which becomes 0 in bitwise operations.
			return 0, nil
		}
		return int(compressed[index] - 32), nil
	})
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(res)), nil
}

// DecompressFromUint8Array takes a compressed slice of bytes and decompresses it into a string.
// It returns an error if the input is not a valid compressed data.
func DecompressFromUint8Array(compressed []byte) (string, error) {
	if compressed == nil {
		return "", ErrInputNil
	}
	if len(compressed) == 0 {
		return "", ErrInputBlank
	}
	length := len(compressed) / 2
	buf := make([]uint16, len(compressed)/2)
	for i := 0; i < length; i++ {
		buf[i] = uint16(compressed[i*2])*256 + uint16(compressed[i*2+1])
	}

	return Decompress(buf)
}

// DecompressFromEncodedURIComponent takes a compressed URL-encoded string and decompresses it into a string.
// It returns an error if the input is not a valid compressed data.
func DecompressFromEncodedURIComponent(compressed string) (string, error) {
	replaced := strings.Replace(compressed, " ", "+", -1)
	if replaced == "" {
		return "", ErrInputBlank
	}
	res, err := _decompress(len(replaced), 32, func(index int) (int, error) {
		if index >= len(replaced) {
			// Match JavaScript behavior: out-of-bounds reads return undefined,
			// which becomes 0 in bitwise operations. This allows the decompression
			// algorithm to properly find the end marker even when it needs to read
			// a few bits past the nominal end of the input.
			return 0, nil
		}
		return getBaseValue(keyStrUriSafe, replaced[index]), nil
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
	dictionary := make(map[int][]uint16)
	var next int
	enlargeIn := 4
	dictSize := 4
	numBits := 3
	var entry []uint16
	result := make([][]uint16, 0)
	var i int
	var bits, resb, maxpower, power int
	var c int
	var w []uint16
	val, err := getNextVal(0)
	if err != nil {
		return nil, err
	}
	data := data{val: val, position: resetValue, index: 1}

	for i = 0; i < 3; i++ {
		dictionary[i] = []uint16{uint16(i)}
	}
	bits = 0
	maxpower = 4 // int(math.Pow(2,2))
	power = 1
	for power != maxpower {
		resb = data.val & data.position
		data.position >>= 1
		if data.position == 0 {
			data.position = resetValue
			data.val, err = getNextVal(data.index)
			if err != nil {
				return nil, err
			}
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
		maxpower = 256 // int(math.Pow(2,8))
		power = 1
		for power != maxpower {
			resb = data.val & data.position
			data.position >>= 1
			if data.position == 0 {
				data.position = resetValue
				data.val, err = getNextVal(data.index)
				if err != nil {
					return nil, err
				}
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
	case 1:
		bits = 0
		maxpower = 65536 // int(math.Pow(2, 16))
		power = 1
		for power != maxpower {
			resb = data.val & data.position
			data.position >>= 1
			if data.position == 0 {
				data.position = resetValue
				data.val, err = getNextVal(data.index)
				if err != nil {
					return nil, err
				}
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
	case 2:
		return nil, nil
	}
	dictionary[3] = []uint16{uint16(c)}
	w = []uint16{uint16(c)}
	result = append(result, []uint16{uint16(c)})
	for {
		if data.index > length {
			return nil, ErrInputNotDecodable
		}
		bits = 0
		maxpower = 1 << numBits
		power = 1
		for power != maxpower {
			resb = data.val & data.position
			data.position >>= 1
			if data.position == 0 {
				data.position = resetValue
				data.val, err = getNextVal(data.index)
				if err != nil {
					return nil, err
				}
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
			maxpower = 256 //int(math.Pow(2, 8))
			power = 1
			for power != maxpower {
				resb = data.val & data.position
				data.position >>= 1
				if data.position == 0 {
					data.position = resetValue
					data.val, err = getNextVal(data.index)
					if err != nil {
						return nil, err
					}
					data.index++
				}
				tmp := 0
				if resb > 0 {
					tmp = 1
				}
				bits |= tmp * power
				power <<= 1
			}

			dictionary[dictSize] = []uint16{uint16(bits)}
			dictSize++
			c = dictSize - 1
			enlargeIn--
		case 1:
			bits = 0
			maxpower = 65536 // int(math.Pow(2, 16))
			power = 1
			for power != maxpower {
				resb = data.val & data.position
				data.position >>= 1
				if data.position == 0 {
					data.position = resetValue
					data.val, err = getNextVal(data.index)
					if err != nil {
						return nil, err
					}
					data.index++
				}
				tmp := 0
				if resb > 0 {
					tmp = 1
				}
				bits |= tmp * power
				power <<= 1
			}
			dictionary[dictSize] = []uint16{uint16(bits)}
			dictSize++
			c = dictSize - 1
			enlargeIn--
		case 2:
			res := make([]uint16, 0)
			for _, r := range result {
				res = append(res, r...)
			}
			return res, nil
		}

		if enlargeIn == 0 {
			enlargeIn = 1 << numBits
			numBits++
		}

		if _, ok := dictionary[c]; ok {
			entry = make([]uint16, len(dictionary[c]))
			copy(entry, dictionary[c])
		} else {
			if c == dictSize {
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
		dictionary[dictSize] = tmp
		dictSize++
		enlargeIn--

		w = entry

		if enlargeIn == 0 {
			enlargeIn = 1 << numBits
			numBits++
		}
	}
}
