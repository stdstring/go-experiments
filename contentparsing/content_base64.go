package contentparsing

import "encoding/base64"

type Base64Data struct {
	Value any
}

func NewBase64Data(value any) *Base64Data {
	return &Base64Data{Value: value}
}

func tryParseBase64(source string) (bool, string) {
	resultData, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		return false, ""
	}
	// TODO (std_string) : think about correct encoding
	result := string(resultData)
	return true, result
}
