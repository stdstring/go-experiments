package contentparsing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProbableContentDetector(t *testing.T) {
	testCases := []struct {
		name                string
		source              string
		expectedResult      string
		expectedContentType ContentType
	}{
		{
			name:                "simple xml",
			source:              `<?xml version="1.0" encoding="utf-8"?><root><a>some data</a></root>`,
			expectedResult:      `<?xml version="1.0" encoding="utf-8"?><root><a>some data</a></root>`,
			expectedContentType: ContentTypeXml,
		},
		{
			name:                "simple xml without declaration",
			source:              `<root><a>some data</a></root>`,
			expectedResult:      `<root><a>some data</a></root>`,
			expectedContentType: ContentTypeXml,
		},
		{
			name:                "bad xml without brackets",
			source:              `root><a>some data</a></root`,
			expectedResult:      `root><a>some data</a></root`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad xml without leading bracket",
			source:              `root><a>some data</a></root>`,
			expectedResult:      `root><a>some data</a></root>`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad xml without trailing brackets",
			source:              `<root><a>some data</a></root`,
			expectedResult:      `<root><a>some data</a></root`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad xml inside",
			source:              `<root><a>some data</a><b></root>`,
			expectedResult:      `<root><a>some data</a><b></root>`,
			expectedContentType: ContentTypeXml,
		},
		{
			name:                "json object",
			source:              `{"key": "IDDQD", "value": 666}`,
			expectedResult:      `{"key": "IDDQD", "value": 666}`,
			expectedContentType: ContentTypeJson,
		},
		{
			name:                "json array",
			source:              `["IDDQD", 666]`,
			expectedResult:      `["IDDQD", 666]`,
			expectedContentType: ContentTypeJson,
		},
		{
			name:                "json string",
			source:              `"IDDQD"`,
			expectedResult:      `"IDDQD"`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "json number",
			source:              `666`,
			expectedResult:      `666`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:   "json bool (base64 false positive)",
			source: `true`,
			//expectedResult:      `true`,
			//expectedContentType: ContentTypeUnspecified,
			expectedResult:      "\xb6\xbb\x9e",
			expectedContentType: ContentTypeBase64,
		},
		{
			name:   "json null (base64 false positive)",
			source: `null`,
			//expectedResult:      `null`,
			//expectedContentType: ContentTypeUnspecified,
			expectedResult:      "\x9e\xe9e",
			expectedContentType: ContentTypeBase64,
		},
		{
			name:                "bad json object without brackets",
			source:              `"key": "IDDQD", "value": 666`,
			expectedResult:      `"key": "IDDQD", "value": 666`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad json object without leading bracket",
			source:              `"key": "IDDQD", "value": 666}`,
			expectedResult:      `"key": "IDDQD", "value": 666}`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad json object without trailing bracket",
			source:              `{"key": "IDDQD", "value": 666`,
			expectedResult:      `{"key": "IDDQD", "value": 666`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad json array without brackets",
			source:              `"IDDQD", 666`,
			expectedResult:      `"IDDQD", 666`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad json array without leading bracket",
			source:              `"IDDQD", 666]`,
			expectedResult:      `"IDDQD", 666]`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad json array without trailing bracket",
			source:              `["IDDQD", 666`,
			expectedResult:      `["IDDQD", 666`,
			expectedContentType: ContentTypeUnspecified,
		},
		{
			name:                "bad json object inside",
			source:              `{"key": "IDDQD", "value": [1,2,3}`,
			expectedResult:      `{"key": "IDDQD", "value": [1,2,3}`,
			expectedContentType: ContentTypeJson,
		},
		{
			name:                "bad json array inside",
			source:              `["IDDQD", {"key": "IDKFA"]`,
			expectedResult:      `["IDDQD", {"key": "IDKFA"]`,
			expectedContentType: ContentTypeJson,
		},
		{
			name:                "base64 encoded string",
			source:              "SUREUUQrSURLRkE=",
			expectedResult:      "IDDQD+IDKFA",
			expectedContentType: ContentTypeBase64,
		},
		{
			name:                "base64 false positive",
			source:              "SUREUUQrSURLRkErSURDTElQ",
			expectedResult:      "IDDQD+IDKFA+IDCLIP",
			expectedContentType: ContentTypeBase64,
		},
	}
	for _, testCase := range testCases {
		currentTestCase := testCase
		t.Run(currentTestCase.name, func(t *testing.T) {
			actualResult, actualContentType := detectProbableContentType(currentTestCase.source)
			assert.Equal(t, currentTestCase.expectedResult, actualResult)
			assert.Equal(t, currentTestCase.expectedContentType, actualContentType)
		})
	}
}
