package expressiontree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExpression(t *testing.T) {
	storage := &parseStorage{
		knownPath: []DataPath{
			CreateDataPathWithMainOnly(HttpDataKey),
			CreateDataPathWithSimpleContent(OptionsKey, "IDDQD"),
			CreateDataPathWithSimpleContent(RequestHeadersKey, "IDKFA"),
			CreateDataPathWithMainOnly(RequestTimeKey),
		},
		checkArguments: []any{
			"",
			"IDCLIP",
		},
	}
	testCases := []struct {
		name          string
		source        string
		expectedError error
	}{
		{
			name:          "EXISTS(http.options.IDDQD)",
			source:        "EXISTS(1)",
			expectedError: nil,
		},
		{
			name:          "MATCH(http.options.IDDQD,666)",
			source:        "MATCH(1,666)",
			expectedError: nil,
		},
		{
			name:          `CHECK(http.options.IDDQD=="IDCLIP")`,
			source:        "CHECK(1,0,1)",
			expectedError: nil,
		},
		{
			name:          "NOT(EXISTS(http.options.IDDQD))",
			source:        "NOT(EXISTS(1))",
			expectedError: nil,
		},
		{
			name:          "AND(EXISTS(http.options.IDDQD),MATCH(http.options.IDDQD,666))",
			source:        "AND(EXISTS(1),MATCH(1,666))",
			expectedError: nil,
		},
		{
			name: "AND(EXISTS(http.options.IDDQD),MATCH(http.options.IDDQD,666)," +
				"EXISTS(http.request.headers.IDKFA),MATCH(http.request.headers.IDKFA,777))",
			source:        "AND(EXISTS(1),MATCH(1,666),EXISTS(2),MATCH(2,777))",
			expectedError: nil,
		},
		{
			name:          "OR(EXISTS(http.options.IDDQD),MATCH(http.options.IDDQD,666))",
			source:        "OR(EXISTS(1),MATCH(1,666))",
			expectedError: nil,
		},
		{
			name: "OR(EXISTS(http.options.IDDQD),MATCH(http.options.IDDQD,666)," +
				"EXISTS(http.request.headers.IDKFA),MATCH(http.request.headers.IDKFA,777))",
			source:        "OR(EXISTS(1),MATCH(1,666),EXISTS(2),MATCH(2,777))",
			expectedError: nil,
		},
		{
			name: "AND(OR(EXISTS(http.options.IDDQD),MATCH(http.options.IDDQD,666))," +
				"OR(NOT(EXISTS(http.request.headers.IDKFA)),MATCH(http.request.headers.IDKFA,777)))",
			source:        "AND(OR(EXISTS(1),MATCH(1,666)),OR(NOT(EXISTS(2)),MATCH(2,777)))",
			expectedError: nil,
		},
		{
			name: "OR(AND(EXISTS(http.options.IDDQD),NOT(MATCH(http.options.IDDQD,666)))," +
				"AND(NOT(EXISTS(http.request.headers.IDKFA)),MATCH(http.request.headers.IDKFA,777)))",
			source:        "OR(AND(EXISTS(1),NOT(MATCH(1,666))),AND(NOT(EXISTS(2)),MATCH(2,777)))",
			expectedError: nil,
		},
		{
			name:          "EXISTS(http.request.time)",
			source:        "EXISTS(3)",
			expectedError: unknownMainPathError,
		},
		{
			name:          "EXISTS()",
			source:        "EXISTS()",
			expectedError: parseError,
		},
		{
			name:          "MATCH()",
			source:        "MATCH()",
			expectedError: badArgsError,
		},
		{
			name:          "MATCH(http.options.IDDQD)",
			source:        "MATCH(1)",
			expectedError: badArgsError,
		},
		{
			name:          "MATCH(http.options.IDDQD,666,777)",
			source:        "MATCH(1,666,777)",
			expectedError: badArgsError,
		},
		{
			name:          `CHECK(http.options.IDDQD<"IDCLIP")`,
			source:        "CHECK(1,2,1)",
			expectedError: unsupportedOperationError,
		},
		{
			name:          `CHECK()`,
			source:        "CHECK()",
			expectedError: badArgsError,
		},
		{
			name:          `CHECK(http.options.IDDQD)`,
			source:        "CHECK(1)",
			expectedError: badArgsError,
		},
		{
			name:          `CHECK(http.options.IDDQD==)`,
			source:        "CHECK(1,0)",
			expectedError: badArgsError,
		},
		{
			name:          `CHECK(http.options.IDDQD=="IDCLIP",666)`,
			source:        "CHECK(1,0,1,666)",
			expectedError: badArgsError,
		},
		{
			name:          "NOT()",
			source:        "NOT()",
			expectedError: parseError,
		},
		{
			name:          "NOT(EXISTS(http.options.IDDQD),EXISTS(http.options.IDDQD))",
			source:        "NOT(EXISTS(1),EXISTS(1))",
			expectedError: parseError,
		},
		{
			name:          "AND()",
			source:        "AND()",
			expectedError: parseError,
		},
		{
			name:          "AND(EXISTS(http.options.IDDQD))",
			source:        "AND(EXISTS(1))",
			expectedError: badArgsError,
		},
		{
			name:          "OR()",
			source:        "OR()",
			expectedError: parseError,
		},
		{
			name:          "OR(EXISTS(http.options.IDDQD))",
			source:        "OR(EXISTS(1))",
			expectedError: badArgsError,
		},
		{
			name:          "EXISTS(http.options.IDDQD,)",
			source:        "EXISTS(1,)",
			expectedError: badArgsError,
		},
		{
			name:          "EXISTS(http.options.IDDQD",
			source:        "EXISTS(1",
			expectedError: parseError,
		},
		{
			name:          "EXISTS(http.options.IDDQD())",
			source:        "EXISTS(1())",
			expectedError: parseError,
		},
		{
			name:          "EXISTS(http.options.IDDQD)EXISTS(http.options.IDDQD)",
			source:        "EXISTS(1)EXISTS(1)",
			expectedError: parseError,
		},
	}
	for _, testCase := range testCases {
		currentTestCase := testCase
		t.Run(currentTestCase.name, func(t *testing.T) {
			checkParse(t, currentTestCase.source, storage, testCase.expectedError)
		})
	}
}

func checkParse(t *testing.T, source string, storage *parseStorage, expectedError error) {
	t.Helper()
	result, actualError := parseExpressionTree(source, storage)
	if expectedError == nil {
		assert.NotNil(t, result)
		assert.NoError(t, actualError)
	} else {
		assert.Nil(t, result)
		assert.Equal(t, expectedError, actualError)
	}
}
