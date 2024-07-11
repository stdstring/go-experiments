package contentparsing

import (
	"strings"
)

type ContentType int

// ContentType
const (
	ContentTypeUnspecified ContentType = 0
	ContentTypeXml         ContentType = 1
	ContentTypeJson        ContentType = 2
	ContentTypeBase64      ContentType = 3
)

// we expect, that source is prepared i.e. without leading and trailing spaces
func detectProbableContentType(source string) (string, ContentType) {
	if strings.HasPrefix(source, "<") && strings.HasSuffix(source, ">") {
		return source, ContentTypeXml
	}
	if strings.HasPrefix(source, "{") && strings.HasSuffix(source, "}") {
		return source, ContentTypeJson
	}
	if strings.HasPrefix(source, "[") && strings.HasSuffix(source, "]") {
		return source, ContentTypeJson
	}
	if parseResult, dest := tryParseBase64(source); parseResult {
		return dest, ContentTypeBase64
	}
	return source, ContentTypeUnspecified
}
