package contentparsing

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXmlContentParse(t *testing.T) {
	const globalMaxLevel = 5
	const maxLevel = 4
	const maxFieldCount = 4
	testCases := []struct {
		name     string
		source   string
		success  bool
		expected *XmlData
	}{
		{
			name:    "simple xml with xml declaration",
			source:  `<?xml version="1.0" encoding="UTF-8"?><root><a><b>IDDQD</b><b>IDKFA</b><c>666</c></a></root>`,
			success: true,
			expected: &XmlData{
				EntityDirectives:  make([]string, 0),
				DoctypeDirectives: make([]string, 0),
				OtherDirectives:   make([]string, 0),
				Root: &XmlElement{
					Name:       createSimpleXmlName("root"),
					Attributes: make([]xml.Attr, 0),
					Value:      nil,
					Children: []*XmlElement{
						{
							Name:       createSimpleXmlName("a"),
							Attributes: make([]xml.Attr, 0),
							Value:      nil,
							Children: []*XmlElement{
								{
									Name:       createSimpleXmlName("b"),
									Attributes: make([]xml.Attr, 0),
									Value:      "IDDQD",
									Children:   make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("b"),
									Attributes: make([]xml.Attr, 0),
									Value:      "IDKFA",
									Children:   make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("c"),
									Attributes: make([]xml.Attr, 0),
									Value:      "666",
									Children:   make([]*XmlElement, 0),
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "simple xml",
			source:  `<root><a><b>IDDQD</b><b>IDKFA</b><c>666</c></a></root>`,
			success: true,
			expected: &XmlData{
				EntityDirectives:  make([]string, 0),
				DoctypeDirectives: make([]string, 0),
				OtherDirectives:   make([]string, 0),
				Root: &XmlElement{
					Name:       createSimpleXmlName("root"),
					Attributes: make([]xml.Attr, 0),
					Value:      nil,
					Children: []*XmlElement{
						{
							Name:       createSimpleXmlName("a"),
							Attributes: make([]xml.Attr, 0),
							Value:      nil,
							Children: []*XmlElement{
								{
									Name:       createSimpleXmlName("b"),
									Attributes: make([]xml.Attr, 0),
									Value:      "IDDQD",
									Children:   make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("b"),
									Attributes: make([]xml.Attr, 0),
									Value:      "IDKFA",
									Children:   make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("c"),
									Attributes: make([]xml.Attr, 0),
									Value:      "666",
									Children:   make([]*XmlElement, 0),
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "simple xml with attributes",
			source:  `<root><a><b attr1="1" attr2="2" attr1="111">IDDQD</b><b>IDKFA</b><c>666</c></a></root>`,
			success: true,
			expected: &XmlData{
				EntityDirectives:  make([]string, 0),
				DoctypeDirectives: make([]string, 0),
				OtherDirectives:   make([]string, 0),
				Root: &XmlElement{
					Name:       createSimpleXmlName("root"),
					Attributes: make([]xml.Attr, 0),
					Value:      nil,
					Children: []*XmlElement{
						{
							Name:       createSimpleXmlName("a"),
							Attributes: make([]xml.Attr, 0),
							Value:      nil,
							Children: []*XmlElement{
								{
									Name: createSimpleXmlName("b"),
									Attributes: []xml.Attr{
										createSimpleXmlAttr("attr1", "1"),
										createSimpleXmlAttr("attr2", "2"),
										createSimpleXmlAttr("attr1", "111"),
									},
									Value:    "IDDQD",
									Children: make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("b"),
									Attributes: make([]xml.Attr, 0),
									Value:      "IDKFA",
									Children:   make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("c"),
									Attributes: make([]xml.Attr, 0),
									Value:      "666",
									Children:   make([]*XmlElement, 0),
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "simple xml with declarations",
			source:  `<root><!ENTITY writer \"Duke Nukem\"><!DOCTYPE root SYSTEM \"root.dtd\"><a><b>IDDQD</b><b>IDKFA</b><c>666</c></a></root>`,
			success: true,
			expected: &XmlData{
				EntityDirectives:  []string{`ENTITY writer \"Duke Nukem\"`},
				DoctypeDirectives: []string{`DOCTYPE root SYSTEM \"root.dtd\"`},
				OtherDirectives:   make([]string, 0),
				Root: &XmlElement{
					Name:       createSimpleXmlName("root"),
					Attributes: make([]xml.Attr, 0),
					Value:      nil,
					Children: []*XmlElement{
						{
							Name:       createSimpleXmlName("a"),
							Attributes: make([]xml.Attr, 0),
							Value:      nil,
							Children: []*XmlElement{
								{
									Name:       createSimpleXmlName("b"),
									Attributes: make([]xml.Attr, 0),
									Value:      "IDDQD",
									Children:   make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("b"),
									Attributes: make([]xml.Attr, 0),
									Value:      "IDKFA",
									Children:   make([]*XmlElement, 0),
								},
								{
									Name:       createSimpleXmlName("c"),
									Attributes: make([]xml.Attr, 0),
									Value:      "666",
									Children:   make([]*XmlElement, 0),
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "simple xml with max level exceeded",
			source:   `<root><a><b><c><d>IDDQD+IDKFA+IDCLIP</d></c></b></a></root>`,
			success:  false,
			expected: nil,
		},
		{
			name:     "xml with absent start element",
			source:   `<root><a><b>IDDQD+IDKFA+IDCLIP</c></b></a></root>`,
			success:  false,
			expected: nil,
		},
		{
			name:     "xml with absent end element",
			source:   `<root><a><b><c>IDDQD+IDKFA+IDCLIP</b></a></root>`,
			success:  false,
			expected: nil,
		},
		{
			name:     "two root pieces of xml",
			source:   `<root1><a>some data</a></root1><root2><b>other data</b></root2>`,
			success:  false,
			expected: nil,
		},
	}
	for _, testCase := range testCases {
		currentTestCase := testCase
		t.Run(currentTestCase.name, func(t *testing.T) {
			state := newContentParserState(globalMaxLevel, maxLevel, maxFieldCount)
			xmlParser := newXmlDataParser(state)
			result, err := xmlParser.parse(currentTestCase.source)
			if currentTestCase.success {
				assert.NoError(t, err)
				checkContent(t, testCase.expected, result)
				prettyPrintContent(result, 0)
			} else {
				assert.Error(t, err)
				assert.Nil(t, result)
				fmt.Printf("Error = %v\n", err)
			}
		})
	}
}

func createSimpleXmlName(local string) xml.Name {
	return createXmlName("", local)
}

func createXmlName(space string, local string) xml.Name {
	return xml.Name{Space: space, Local: local}
}

func createXmlAttr(space string, local string, value string) xml.Attr {
	return xml.Attr{Name: createXmlName(space, local), Value: value}
}

func createSimpleXmlAttr(local string, value string) xml.Attr {
	return createXmlAttr("", local, value)
}
