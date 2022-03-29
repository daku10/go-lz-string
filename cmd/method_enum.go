package cmd

import "errors"

type methodEnum string

const (
	methodInvalidUTF16        methodEnum = "invalid-utf16"
	methodBase64                         = "base64"
	methodUTF16                          = "utf16"
	methodUint8Array                     = "uint8array"
	methodEncodedURIComponent            = "encodedURIComponent"
)

func (e *methodEnum) String() string {
	return string(*e)
}

func (e *methodEnum) Set(v string) error {
	switch v {
	case "invalid-utf16", "base64", "utf16", "uint8array", "encodedURIComponent":
		*e = methodEnum(v)
		return nil
	default:
		return errors.New(`must be one of "invalid-utf16", "base64", "utf16", "uint8array" or "encodedURIComponent"`)
	}
}

func (e *methodEnum) Type() string {
	return "method"
}
