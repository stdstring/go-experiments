package contentparsing

import (
	"encoding/xml"
	"testing"
)

func TestContentParse(t *testing.T) {
	const maxGlobalLevel = 6
	const maxLocalLevel = 4
	const maxFieldCount = 4
	testCases := []struct {
		name     string
		source   string
		expected any
	}{
		{
			name:   "simple xml",
			source: `<?xml version="1.0" encoding="UTF-8"?><root><a><b>IDDQD</b><b>IDKFA</b><c>666</c></a></root>`,
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
			name:     "simple json array",
			source:   `[1,2,3]`,
			expected: &JsonData{Value: JsonArrayValue{1.0, 2.0, 3.0}},
		},
		{
			name:     "simple json object",
			source:   `{"key": "IDDQD", "value": "666"}`,
			expected: &JsonData{Value: JsonObjectValue{"key": "IDDQD", "value": "666"}},
		},
		{
			name:   "json array in xml",
			source: `<root><a><b>IDDQD</b><c>[1,2,3]</c></a></root>`,
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
									Name:       createSimpleXmlName("c"),
									Attributes: make([]xml.Attr, 0),
									Value:      &JsonData{Value: JsonArrayValue{1.0, 2.0, 3.0}},
									Children:   make([]*XmlElement, 0),
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "json object in xml",
			source: `<root><a><b>IDDQD</b><c>{"key": "IDDQD", "value": 666}</c></a></root>`,
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
									Name:       createSimpleXmlName("c"),
									Attributes: make([]xml.Attr, 0),
									Value:      &JsonData{Value: JsonObjectValue{"key": "IDDQD", "value": 666.0}},
									Children:   make([]*XmlElement, 0),
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "xml in json array",
			source: `[666,"<root><a><b>IDKFA</b></a></root>","IDDQD",null]`,
			expected: &JsonData{
				Value: JsonArrayValue{
					666.0,
					&XmlData{
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
											Value:      "IDKFA",
											Children:   make([]*XmlElement, 0),
										},
									},
								},
							},
						},
					},
					"IDDQD",
					nil,
				},
			},
		},
		{
			name:   "xml in json object",
			source: `{"key": "IDDQD", "value": "<root><a><b>IDKFA</b></a></root>"}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"key": "IDDQD",
					"value": &XmlData{
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
											Value:      "IDKFA",
											Children:   make([]*XmlElement, 0),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "xml in json object in xml",
			source: `<root><a>{"key": "IDDQD", "value": "&lt;root1&gt;&lt;b&gt;IDKFA&lt;/b&gt;&lt;/root1&gt;"}</a></root>`,
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
							Value: &JsonData{
								Value: JsonObjectValue{
									"key": "IDDQD",
									"value": &XmlData{
										EntityDirectives:  make([]string, 0),
										DoctypeDirectives: make([]string, 0),
										OtherDirectives:   make([]string, 0),
										Root: &XmlElement{
											Name:       createSimpleXmlName("root1"),
											Attributes: make([]xml.Attr, 0),
											Value:      nil,
											Children: []*XmlElement{
												{
													Name:       createSimpleXmlName("b"),
													Attributes: make([]xml.Attr, 0),
													Value:      "IDKFA",
													Children:   make([]*XmlElement, 0),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "json object in xml in json object",
			source: `{"key": "IDDQD", "value": "<root><a>{\"name\": \"IDKFA\", \"data\": 666}</a></root>"}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"key": "IDDQD",
					"value": &XmlData{
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
									Children:   make([]*XmlElement, 0),
									Value:      &JsonData{Value: JsonObjectValue{"name": "IDKFA", "data": 666.0}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "json object in xml with local maxLevel exceed",
			source: `<root>{"entry": {"data": {"object": {"record": {"item": 666}}}}}</root>`,
			expected: &XmlData{
				EntityDirectives:  make([]string, 0),
				DoctypeDirectives: make([]string, 0),
				OtherDirectives:   make([]string, 0),
				Root: &XmlElement{
					Name:       createSimpleXmlName("root"),
					Attributes: make([]xml.Attr, 0),
					Children:   make([]*XmlElement, 0),
					Value:      `{"entry": {"data": {"object": {"record": {"item": 666}}}}}`,
				},
			},
		},
		{
			name:   "json object in xml with global maxLevel exceed",
			source: `<root><a><b>{"entry": {"data": {"record": {"item": 666}}}}</b></a></root>`,
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
									Children:   make([]*XmlElement, 0),
									Value:      `{"entry": {"data": {"record": {"item": 666}}}}`,
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "json array in xml with local maxLevel exceed",
			source: `<root>[1, [2, [3, [4, [5]]]]]</root>`,
			expected: &XmlData{
				EntityDirectives:  make([]string, 0),
				DoctypeDirectives: make([]string, 0),
				OtherDirectives:   make([]string, 0),
				Root: &XmlElement{
					Name:       createSimpleXmlName("root"),
					Attributes: make([]xml.Attr, 0),
					Children:   make([]*XmlElement, 0),
					Value:      `[1, [2, [3, [4, [5]]]]]`,
				},
			},
		},
		{
			name:   "json array in xml with global maxLevel exceed",
			source: `<root><a><b>[1, [2, [3, [4]]]]</b></a></root>`,
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
									Children:   make([]*XmlElement, 0),
									Value:      `[1, [2, [3, [4]]]]`,
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "xml in json object with local maxLevel exceed",
			source: `{"key": "IDDQD", "value": "<root><a><b><c><d>IDKFA</d></c></b></a></root>"}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"key":   "IDDQD",
					"value": "<root><a><b><c><d>IDKFA</d></c></b></a></root>",
				},
			},
		},
		{
			name:   "xml in json object with global maxLevel exceed",
			source: `{"data": {"record": {"entry": "<root><a><b><c>IDKFA</c></b></a></root>"}}}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"data": JsonObjectValue{
						"record": JsonObjectValue{
							"entry": "<root><a><b><c>IDKFA</c></b></a></root>",
						},
					},
				},
			},
		},
		{
			name:   "xml in json array with local maxLevel exceed",
			source: `[1, "<root><a><b><c><d>IDKFA</d></c></b></a></root>", "IDDQD", null]`,
			expected: &JsonData{
				Value: JsonArrayValue{1.0, "<root><a><b><c><d>IDKFA</d></c></b></a></root>", "IDDQD", nil},
			},
		},
		{
			name:   "xml in json array with global maxLevel exceed",
			source: `[1, [2, [3, "<root><a><b><c>IDKFA</c></b></a></root>"]]]`,
			expected: &JsonData{
				Value: JsonArrayValue{1.0, JsonArrayValue{2.0, JsonArrayValue{3.0, "<root><a><b><c>IDKFA</c></b></a></root>"}}}},
		},
		{
			name:   "bad json object in xml",
			source: `<root><a>{"data": 666</a></root>`,
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
							Children:   make([]*XmlElement, 0),
							Value:      `{"data": 666`,
						},
					},
				},
			},
		},
		{
			name:   "bad json array in xml",
			source: `<root><a>["data", 666</a></root>`,
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
							Children:   make([]*XmlElement, 0),
							Value:      `["data", 666`,
						},
					},
				},
			},
		},
		{
			name:     "bad xml in json object",
			source:   `{"key": "IDDQD", "value": "<root><a>IDKFA</a>"}`,
			expected: &JsonData{Value: JsonObjectValue{"key": "IDDQD", "value": "<root><a>IDKFA</a>"}},
		},
		{
			name:     "bad xml in json array",
			source:   `[1,"<root><a>IDKFA</a>"]`,
			expected: &JsonData{Value: JsonArrayValue{1.0, "<root><a>IDKFA</a>"}},
		},
		{
			name:   "xml in bad json object in xml",
			source: `<root><a>{"data": "&lt;root1&gt;&lt;b&gt;IDKFA&lt;/b&gt;&lt;/root1&gt;"</a></root>`,
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
							Children:   make([]*XmlElement, 0),
							Value:      `{"data": "<root1><b>IDKFA</b></root1>"`,
						},
					},
				},
			},
		},
		{
			name:     "json object in bad xml in json object",
			source:   `{"key": "IDDQD", "value": "<root><a>{\"data\": 666}</a>"}`,
			expected: &JsonData{Value: JsonObjectValue{"key": "IDDQD", "value": `<root><a>{"data": 666}</a>`}},
		},
		{
			name: "base64 string",
			// IDDQD+IDKFA+IDCLIP
			source:   `SUREUUQrSURLRkErSURDTElQ`,
			expected: &Base64Data{Value: "IDDQD+IDKFA+IDCLIP"},
		},
		{
			name: "base64 simple xml",
			// <root><a><b>IDDQD</b><b>IDKFA</b><c>666</c></a></root>
			source: `PHJvb3Q+PGE+PGI+SUREUUQ8L2I+PGI+SURLRkE8L2I+PGM+NjY2PC9jPjwvYT48L3Jvb3Q+`,
			expected: &Base64Data{
				Value: &XmlData{
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
		},
		{
			name: "base64 simple json array",
			// [1,2,3]
			source: `WzEsMiwzXQ==`,
			expected: &Base64Data{
				Value: &JsonData{Value: JsonArrayValue{1.0, 2.0, 3.0}},
			},
		},
		{
			name: "base64 simple json object",
			// {"key": "IDDQD", "value": "666"}
			source: `eyJrZXkiOiAiSUREUUQiLCAidmFsdWUiOiAiNjY2In0=`,
			expected: &Base64Data{
				Value: &JsonData{Value: JsonObjectValue{"key": "IDDQD", "value": "666"}},
			},
		},
		{
			name: "base64 bad xml",
			// <root><a>IDDQD</a>
			source:   `PHJvb3Q+PGE+SUREUUQ8L2E+`,
			expected: &Base64Data{Value: "<root><a>IDDQD</a>"},
		},
		{
			name: "base64 bad json",
			// {"key": "IDDQD", "value": "666"
			source:   `eyJrZXkiOiAiSUREUUQiLCAidmFsdWUiOiAiNjY2Ig==`,
			expected: &Base64Data{Value: `{"key": "IDDQD", "value": "666"`},
		},
		{
			name: "base64 string in xml",
			// IDDQD+IDKFA+IDCLIP
			source: `<root><a>SUREUUQrSURLRkErSURDTElQ</a></root>`,
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
							Children:   make([]*XmlElement, 0),
							Value:      &Base64Data{Value: "IDDQD+IDKFA+IDCLIP"},
						},
					},
				},
			},
		},
		{
			name: "base64 xml in xml",
			// <root1><b>IDDQD</b></root1>
			source: `<root><a>PHJvb3QxPjxiPklERFFEPC9iPjwvcm9vdDE+</a></root>`,
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
							Children:   make([]*XmlElement, 0),
							Value: &Base64Data{
								Value: &XmlData{
									EntityDirectives:  make([]string, 0),
									DoctypeDirectives: make([]string, 0),
									OtherDirectives:   make([]string, 0),
									Root: &XmlElement{
										Name:       createSimpleXmlName("root1"),
										Attributes: make([]xml.Attr, 0),
										Value:      nil,
										Children: []*XmlElement{
											{
												Name:       createSimpleXmlName("b"),
												Attributes: make([]xml.Attr, 0),
												Children:   make([]*XmlElement, 0),
												Value:      "IDDQD",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "base64 json in xml",
			// {"key": "IDDQD", "value": "666"}
			source: `<root><a>eyJrZXkiOiAiSUREUUQiLCAidmFsdWUiOiAiNjY2In0=</a></root>`,
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
							Children:   make([]*XmlElement, 0),
							Value: &Base64Data{
								Value: &JsonData{Value: JsonObjectValue{"key": "IDDQD", "value": "666"}},
							},
						},
					},
				},
			},
		},
		{
			name: "base64 json in json",
			// {"key": "IDDQD", "value": "666"}
			source: `{"data": "eyJrZXkiOiAiSUREUUQiLCAidmFsdWUiOiAiNjY2In0="}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"data": &Base64Data{
						Value: &JsonData{Value: JsonObjectValue{"key": "IDDQD", "value": "666"}},
					},
				},
			},
		},
		{
			name: "base64 xml in json",
			// <root1><b>IDDQD</b></root1>
			source: `{"data": "PHJvb3QxPjxiPklERFFEPC9iPjwvcm9vdDE+"}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"data": &Base64Data{
						Value: &XmlData{
							EntityDirectives:  make([]string, 0),
							DoctypeDirectives: make([]string, 0),
							OtherDirectives:   make([]string, 0),
							Root: &XmlElement{
								Name:       createSimpleXmlName("root1"),
								Attributes: make([]xml.Attr, 0),
								Value:      nil,
								Children: []*XmlElement{
									{
										Name:       createSimpleXmlName("b"),
										Attributes: make([]xml.Attr, 0),
										Children:   make([]*XmlElement, 0),
										Value:      "IDDQD",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "base64 xml in json with total level == globalMaxLevel",
			// <root><a><b><c>IDDQD</c></b></a></root>
			source: `{"data": {"item": "PHJvb3Q+PGE+PGI+PGM+SUREUUQ8L2M+PC9iPjwvYT48L3Jvb3Q+"}}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"data": JsonObjectValue{
						"item": &Base64Data{
							Value: &XmlData{
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
													Value:      nil,
													Children: []*XmlElement{
														{
															Name:       createSimpleXmlName("c"),
															Attributes: make([]xml.Attr, 0),
															Children:   make([]*XmlElement, 0),
															Value:      "IDDQD",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "base64 xml maxLevel in json with global maxLevel exceed",
			// <root><a><b><c>IDDQD</c></b></a></root>
			source: `{"data": {"item": {"value": "PHJvb3Q+PGE+PGI+PGM+SUREUUQ8L2M+PC9iPjwvYT48L3Jvb3Q+"}}}`,
			expected: &JsonData{
				Value: JsonObjectValue{
					"data": JsonObjectValue{
						"item": JsonObjectValue{
							"value": &Base64Data{
								Value: "<root><a><b><c>IDDQD</c></b></a></root>",
							},
						},
					},
				},
			},
		},
		{
			name: "base64 json in xml with total level == globalMaxLevel",
			// {"data": {"record": {"object": {"value": "IDDQD"}}}}
			source: "<root><a>eyJkYXRhIjogeyJyZWNvcmQiOiB7Im9iamVjdCI6IHsidmFsdWUiOiAiSUREUUQifX19fQ==</a></root>",
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
							Children:   make([]*XmlElement, 0),
							Value: &Base64Data{
								Value: &JsonData{
									Value: JsonObjectValue{
										"data": JsonObjectValue{
											"record": JsonObjectValue{
												"object": JsonObjectValue{
													"value": "IDDQD",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "base64 json in xml with global maxLevel exceed",
			// {"data": {"record": {"object": {"value": "IDDQD"}}}}
			source: "<root><a><b>eyJkYXRhIjogeyJyZWNvcmQiOiB7Im9iamVjdCI6IHsidmFsdWUiOiAiSUREUUQifX19fQ==</b></a></root>",
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
									Children:   make([]*XmlElement, 0),
									Value: &Base64Data{
										Value: `{"data": {"record": {"object": {"value": "IDDQD"}}}}`,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "base64 json object in xml with maxFieldCount exceed",
			// {"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
			source: "<root><a>eyJhIjogMSwgImIiOiAyLCAiYyI6IDMsICJkIjogNCwgImUiOiA1fQ==</a></root>",
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
							Children:   make([]*XmlElement, 0),
							Value: &Base64Data{
								Value: `{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}`,
							},
						},
					},
				},
			},
		},
		{
			name: "base64 json array in xml with maxFieldCount exceed",
			// [1,2,3,4,5]
			source: "<root><a>WzEsMiwzLDQsNV0=</a></root>",
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
							Children:   make([]*XmlElement, 0),
							Value: &Base64Data{
								Value: `[1,2,3,4,5]`,
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		currentTestCase := testCase
		t.Run(currentTestCase.name, func(t *testing.T) {
			params := ParseParams{
				MaxGlobalLevel: maxGlobalLevel,
				MaxLocalLevel:  maxLocalLevel,
				MaxFieldCount:  maxFieldCount,
			}
			result := ParseContent(currentTestCase.source, params)
			checkContent(t, testCase.expected, result)
			prettyPrintContent(result, 0)
		})
	}
}
