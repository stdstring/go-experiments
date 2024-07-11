package contentparsing

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var maxFieldCountExceedError = errors.New("max field count exceeded")
var badFieldNameError = errors.New("bad field name")
var unexpectedEndOfArrayError = errors.New("unexpected end of array")
var unexpectedEndOfObjectError = errors.New("unexpected end of object")
var badJsonStructureError = errors.New("bad json structure")

// json = null | bool | number | string | array | object

type JsonSimpleValue = any

type JsonArrayValue = []any

type JsonObjectValue = map[string]any

type JsonData struct {
	Value any // JsonSimpleValue | *XmlData | JsonArrayValue | JsonObjectValue
}

func NewJsonData() *JsonData {
	return &JsonData{
		Value: nil,
	}
}

type jsonDataParser struct {
	decoder           *json.Decoder
	data              *JsonData
	state             contentParserState
	localLevelChecker *levelChecker
}

func (parser *jsonDataParser) parse(source string) (*JsonData, error) {
	parser.data = NewJsonData()
	parser.localLevelChecker = newLevelChecker(parser.state.maxLevel)
	data := bytes.NewBufferString(source)
	parser.decoder = json.NewDecoder(data)
	value, valueError := parser.parseValue()
	if valueError != nil {
		return nil, valueError
	}
	_, tokenError := parser.decoder.Token()
	if !errors.Is(tokenError, io.EOF) {
		return nil, badJsonStructureError
	}
	parser.data.Value = value
	return parser.data, nil
}

func (parser *jsonDataParser) parseValue() (any, error) {
	rawValue, valueError := parser.decoder.Token()
	if errors.Is(valueError, io.EOF) {
		return nil, unexpectedEOFError
	}
	if valueError != nil {
		return nil, valueError
	}
	switch value := rawValue.(type) {
	case json.Delim:
		complexEntity, complexEntityError := parser.parseComplexEntity(value, unexpectedEndOfArrayError)
		if complexEntityError != nil {
			return nil, complexEntityError
		}
		return complexEntity, nil
	case bool:
		return JsonSimpleValue(value), nil
	case float64:
		return JsonSimpleValue(value), nil
	case nil:
		return JsonSimpleValue(value), nil
	case string:
		return parseContent(value, parser.state), nil
	}
	return nil, unexpectedValueError
}

func (parser *jsonDataParser) parseObject() (any, error) {
	if !parser.localLevelChecker.enter() {
		return nil, maxLevelExceedError
	}
	if !parser.state.globalLevelChecker.enter() {
		return nil, maxLevelExceedError
	}
	defer func() {
		parser.state.globalLevelChecker.exit()
		parser.localLevelChecker.exit()
	}()
	destObject := make(JsonObjectValue)
	for parser.decoder.More() {
		if len(destObject) == parser.state.maxFieldCount {
			return nil, maxFieldCountExceedError
		}
		fieldName, nameError := parser.parseFieldName()
		if nameError != nil {
			return nil, nameError
		}
		value, valueError := parser.parseValue()
		if valueError != nil {
			return nil, valueError
		}
		destObject[fieldName] = value
	}
	endDelimiterError := parser.parseEndDelimiter("}", unexpectedEndOfObjectError)
	if endDelimiterError != nil {
		return nil, endDelimiterError
	}
	return destObject, nil
}

func (parser *jsonDataParser) parseArray() (any, error) {
	if !parser.localLevelChecker.enter() {
		return nil, maxLevelExceedError
	}
	if !parser.state.globalLevelChecker.enter() {
		return nil, maxLevelExceedError
	}
	defer func() {
		parser.state.globalLevelChecker.exit()
		parser.localLevelChecker.exit()
	}()
	destArray := make(JsonArrayValue, 0)
	for parser.decoder.More() {
		if len(destArray) == parser.state.maxFieldCount {
			return nil, maxFieldCountExceedError
		}
		value, valueError := parser.parseValue()
		if valueError != nil {
			return nil, valueError
		}
		destArray = append(destArray, value)
	}
	endDelimiterError := parser.parseEndDelimiter("]", unexpectedEndOfArrayError)
	if endDelimiterError != nil {
		return nil, endDelimiterError
	}
	return destArray, nil
}

func (parser *jsonDataParser) parseComplexEntity(delimiter json.Delim, unexpectedError error) (any, error) {
	switch delimiter.String() {
	case "{":
		objectValue, parseError := parser.parseObject()
		if parseError != nil {
			return nil, parseError
		}
		return objectValue, nil
	case "[":
		arrayValue, parseError := parser.parseArray()
		if parseError != nil {
			return nil, parseError
		}
		return arrayValue, nil
	default:
		return nil, unexpectedError
	}
}

func (parser *jsonDataParser) parseFieldName() (string, error) {
	rawName, nameError := parser.decoder.Token()
	if errors.Is(nameError, io.EOF) {
		return "", unexpectedEOFError
	}
	if nameError != nil {
		return "", nameError
	}
	switch name := rawName.(type) {
	case string:
		return name, nil
	default:
		return "", badFieldNameError
	}
}

func (parser *jsonDataParser) parseEndDelimiter(expectedDelimiter string, unexpectedError error) error {
	rawToken, tokenError := parser.decoder.Token()
	if errors.Is(tokenError, io.EOF) {
		return unexpectedEOFError
	}
	if tokenError != nil {
		return tokenError
	}
	switch token := rawToken.(type) {
	case json.Delim:
		if token.String() == expectedDelimiter {
			return nil
		}
		return unexpectedError
	default:
		return unexpectedError
	}
}

func newJsonDataParser(state contentParserState) *jsonDataParser {
	return &jsonDataParser{
		decoder:           nil,
		data:              nil,
		state:             state,
		localLevelChecker: nil,
	}
}
