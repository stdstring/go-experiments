package contentparsing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonContentParse(t *testing.T) {
	const globalMaxLevel = 5
	const maxLevel = 4
	const maxFieldCount = 4
	testCases := []struct {
		name     string
		source   string
		success  bool
		expected *JsonData
	}{
		{
			name:    "simple json array",
			source:  `[1,2,3]`,
			success: true,
			expected: &JsonData{
				Value: JsonArrayValue{1.0, 2.0, 3.0},
			},
		},
		{
			name:    "simple json object",
			source:  `{"key": "IDDQD", "value": "666"}`,
			success: true,
			expected: &JsonData{
				Value: JsonObjectValue{"key": "IDDQD", "value": "666"},
			},
		},
		{
			name:    "complex json array",
			source:  `[666,"IDDQD",{"key": "IDDQD", "value": {"param": "IDKFA", "values": [666,777]}},null]`,
			success: true,
			expected: &JsonData{
				Value: JsonArrayValue{
					666.0,
					"IDDQD",
					JsonObjectValue{
						"key": "IDDQD",
						"value": JsonObjectValue{
							"param":  "IDKFA",
							"values": JsonArrayValue{666.0, 777.0},
						},
					},
					nil,
				},
			},
		},
		{
			name:    "complex json object",
			source:  `{"key": ["IDDQD", "IDKFA"], "value": ["666", {"number": 777}, {"number": 888}]}`,
			success: true,
			expected: &JsonData{
				Value: JsonObjectValue{
					"key":   JsonArrayValue{"IDDQD", "IDKFA"},
					"value": JsonArrayValue{"666", JsonObjectValue{"number": 777.0}, JsonObjectValue{"number": 888.0}},
				},
			},
		},
		{
			name:     "simple json array with max level exceeded",
			source:   `[1,[1,[1,[1,[1,7],6],5],4],3]`,
			success:  false,
			expected: nil,
		},
		{
			name:     "simple json object max level exceeded",
			source:   `{"data": {"value": {"key": {"number": {"item": 666}}}}}`,
			success:  false,
			expected: nil,
		},
		{
			name:     "json array with max level exceeded",
			source:   `[1,{"key": [2, {"data": [666]}]}]`,
			success:  false,
			expected: nil,
		},
		{
			name:     "json object max level exceeded",
			source:   `{"data": [1, {"value": [2, {"key": "IDDQD"}]}]}`,
			success:  false,
			expected: nil,
		},
		{
			name:     "simple json array with max field count exceeded",
			source:   `[1,2,3,4,5]`,
			success:  false,
			expected: nil,
		},
		{
			name:     "simple json object max field count exceeded",
			source:   `{"p1": 1, "p2": 2, "p3": 3, "p4": 4, "p5": 5}`,
			success:  false,
			expected: nil,
		},
		{
			name:     "json array with max field count exceeded",
			source:   `[1,[1,2,3,4,5],666]`,
			success:  false,
			expected: nil,
		},
		{
			name:     "simple json object max field count exceeded",
			source:   `{"key": "IDDQD", "value": {"p1": 1, "p2": 2, "p3": 3, "p4": 4, "p5": 5}}`,
			success:  false,
			expected: nil,
		},
		{
			name:     "simple json array without trailing bracket",
			source:   `[1,2,3`,
			success:  false,
			expected: nil,
		},
		{
			name:     "nested json arrays without trailing brackets",
			source:   `[666,[1,2,3`,
			success:  false,
			expected: nil,
		},
		{
			name:     "simple json object without trailing bracket",
			source:   `{"key": "IDDQD", "value": "666"`,
			success:  false,
			expected: nil,
		},
		{
			name:     "nested json object without trailing brackets",
			source:   `{"key": "IDDQD", "value": {"data": 666`,
			success:  false,
			expected: nil,
		},
		{
			name:     "two root object",
			source:   `{"key": "IDDQD"}{"key": "IDKFA"}`,
			success:  false,
			expected: nil,
		},
		{
			name:     "two root arrays",
			source:   `["IDDQD", 666]["IDKFA", 666]`,
			success:  false,
			expected: nil,
		},
		{
			name:     "array without commas",
			source:   `[11 22]`,
			success:  false,
			expected: nil,
		},
		{
			name:     "object without quotes",
			source:   `{key: IDDQD, value: 666}`,
			success:  false,
			expected: nil,
		},
		{
			name:     "object without commas",
			source:   `{"key": "IDDQD" "value": 666}`,
			success:  false,
			expected: nil,
		},
		{
			name:     "object without colons",
			source:   `{"key" "IDDQD", "value" 666}`,
			success:  false,
			expected: nil,
		},
	}
	for _, testCase := range testCases {
		currentTestCase := testCase
		t.Run(currentTestCase.name, func(t *testing.T) {
			state := newContentParserState(globalMaxLevel, maxLevel, maxFieldCount)
			jsonParser := newJsonDataParser(state)
			result, err := jsonParser.parse(currentTestCase.source)
			if currentTestCase.success {
				assert.NoError(t, err)
				checkContent(t, currentTestCase.expected, result)
				prettyPrintContent(result, 0)
			} else {
				assert.Error(t, err)
				assert.Nil(t, result)
				fmt.Printf("Error = %v\n", err)
			}
		})
	}
}
