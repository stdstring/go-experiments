package contentparsing

import (
	"errors"
	"strings"
)

var maxLevelExceedError = errors.New("max level exceeded")
var unexpectedEOFError = errors.New("unexpected EOF")
var unexpectedValueError = errors.New("unexpected value")

type levelChecker struct {
	currentLevel int
	maxLevel     int
}

func (checker *levelChecker) enter() bool {
	checker.currentLevel++
	return checker.currentLevel <= checker.maxLevel
}

func (checker *levelChecker) exit() {
	checker.currentLevel--
}

func newLevelChecker(maxLevel int) *levelChecker {
	return &levelChecker{currentLevel: 0, maxLevel: maxLevel}
}

type contentParserState struct {
	globalLevelChecker *levelChecker
	maxLevel           int
	maxFieldCount      int
}

func newContentParserState(globalMaxLevel int, maxLevel int, maxFieldCount int) contentParserState {
	return contentParserState{
		globalLevelChecker: newLevelChecker(globalMaxLevel),
		maxLevel:           maxLevel,
		maxFieldCount:      maxFieldCount,
	}
}

func parseContent(source string, state contentParserState) any {
	preparedSource := strings.TrimSpace(source)
	data, contentType := detectProbableContentType(preparedSource)
	switch contentType {
	case ContentTypeUnspecified:
		return source
	case ContentTypeXml:
		xmlParser := newXmlDataParser(state)
		result, parseError := xmlParser.parse(data)
		if parseError != nil {
			return source
		}
		return result
	case ContentTypeJson:
		jsonParser := newJsonDataParser(state)
		result, parseError := jsonParser.parse(data)
		if parseError != nil {
			return source
		}
		return result
	case ContentTypeBase64:
		result := parseContent(data, state)
		return NewBase64Data(result)
	default:
		return source
	}
}

type ParseParams struct {
	MaxGlobalLevel int
	MaxLocalLevel  int
	MaxFieldCount  int
}

func ParseContent(source string, params ParseParams) any {
	state := newContentParserState(params.MaxGlobalLevel, params.MaxLocalLevel, params.MaxFieldCount)
	return parseContent(source, state)
}
