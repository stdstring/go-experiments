package contentparsing

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

var unexpectedEndElementError = errors.New("unexpected end element")
var badXMLStructureError = errors.New("bad xml structure")

type XmlElement struct {
	Name       xml.Name
	Attributes []xml.Attr
	Children   []*XmlElement
	Value      any // string | *XmlData | *JsonData
}

type XmlData struct {
	EntityDirectives  []string
	DoctypeDirectives []string
	OtherDirectives   []string
	Root              *XmlElement
}

func NewXmlData() *XmlData {
	return &XmlData{
		EntityDirectives:  make([]string, 0),
		DoctypeDirectives: make([]string, 0),
		OtherDirectives:   make([]string, 0),
		Root:              nil,
	}
}

type xmlDataParser struct {
	decoder           *xml.Decoder
	data              *XmlData
	state             contentParserState
	localLevelChecker *levelChecker
}

func (parser *xmlDataParser) parse(source string) (*XmlData, error) {
	parser.data = NewXmlData()
	parser.localLevelChecker = newLevelChecker(parser.state.maxLevel)
	data := bytes.NewBufferString(source)
	parser.decoder = xml.NewDecoder(data)
	for {
		rawToken, tokenError := parser.decoder.Token()
		switch {
		case errors.Is(tokenError, io.EOF):
			return parser.data, nil
		case tokenError != nil:
			return nil, tokenError
		}
		switch token := rawToken.(type) {
		case xml.CharData:
			return nil, badXMLStructureError
		case xml.Comment:
			if parseError := parser.parseComment(token); parseError != nil {
				return nil, parseError
			}
		case xml.Directive:
			if parseError := parser.parseDirective(token); parseError != nil {
				return nil, parseError
			}
		case xml.ProcInst:
			if parseError := parser.parseProcInst(token); parseError != nil {
				return nil, parseError
			}
		case xml.EndElement:
			return nil, badXMLStructureError
		case xml.StartElement:
			if parser.data.Root != nil {
				return nil, badXMLStructureError
			}
			rootElement, parseError := parser.parseElement(token)
			if parseError != nil {
				return nil, parseError
			}
			parser.data.Root = rootElement
		}
	}
}

func (parser *xmlDataParser) parseCharData(token xml.CharData) (any, error) {
	// TODO (std_string) : think about correct encoding
	source := string(token)
	return parseContent(source, parser.state), nil
}

func (parser *xmlDataParser) parseComment(_ xml.Comment) error {
	// do nothing
	return nil
}

func (parser *xmlDataParser) parseDirective(token xml.Directive) error {
	directive := string(token)
	switch {
	case strings.HasPrefix(directive, "ENTITY"):
		parser.data.EntityDirectives = append(parser.data.EntityDirectives, directive)
	case strings.HasPrefix(directive, "DOCTYPE"):
		parser.data.DoctypeDirectives = append(parser.data.DoctypeDirectives, directive)
	default:
		parser.data.OtherDirectives = append(parser.data.OtherDirectives, directive)
	}
	return nil
}

func (parser *xmlDataParser) parseProcInst(_ xml.ProcInst) error {
	// do nothing
	return nil
}

func (parser *xmlDataParser) parseElement(startElement xml.StartElement) (*XmlElement, error) {
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
	currentElement := &XmlElement{
		Name:       startElement.Name,
		Attributes: startElement.Attr,
		Children:   make([]*XmlElement, 0),
		Value:      nil,
	}
	for {
		rawToken, tokenError := parser.decoder.Token()
		switch {
		case errors.Is(tokenError, io.EOF):
			return nil, unexpectedEOFError
		case tokenError != nil:
			return nil, tokenError
		}
		switch token := rawToken.(type) {
		case xml.CharData:
			value, parseError := parser.parseCharData(token)
			if parseError != nil {
				return nil, parseError
			}
			currentElement.Value = value
		case xml.Comment:
			if parseError := parser.parseComment(token); parseError != nil {
				return nil, parseError
			}
		case xml.Directive:
			if parseError := parser.parseDirective(token); parseError != nil {
				return nil, parseError
			}
		case xml.ProcInst:
			if parseError := parser.parseProcInst(token); parseError != nil {
				return nil, parseError
			}
		case xml.EndElement:
			if token.Name != startElement.Name {
				return nil, unexpectedEndElementError
			}
			return currentElement, nil
		case xml.StartElement:
			childElement, parseError := parser.parseElement(token)
			if parseError != nil {
				return nil, parseError
			}
			currentElement.Children = append(currentElement.Children, childElement)
		}
	}
}

func newXmlDataParser(state contentParserState) *xmlDataParser {
	return &xmlDataParser{
		decoder:           nil,
		data:              nil,
		state:             state,
		localLevelChecker: nil,
	}
}

func getXmlFullName(name xml.Name) string {
	if len(name.Space) == 0 {
		return name.Local
	}
	return fmt.Sprintf("%s:%s", name.Space, name.Local)
}
