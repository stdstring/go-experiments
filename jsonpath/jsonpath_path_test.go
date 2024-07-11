package jsonpath

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPathParsePath(t *testing.T) {
	out, err := parsePathString("$")
	assert.NoError(t, err)
	assert.Equal(t, len(out), 0)
}

func TestJSONPathParsePathLong(t *testing.T) {
	out, err := parsePathString("$.a.b.c")
	assert.NoError(t, err)
	assert.Equal(t, len(out), 3)
	assert.Equal(t, out[0], "a")
	assert.Equal(t, out[1], "b")
	assert.Equal(t, out[2], "c")
}

func TestJSONPathNewPath(t *testing.T) {
	path, err := NewPath("$.a.b.c")
	assert.NoError(t, err)
	assert.Equal(t, len(path.path), 3)
	assert.Equal(t, path.path[0], "a")
	assert.Equal(t, path.path[1], "b")
	assert.Equal(t, path.path[2], "c")
}

func TestJSONPathParsing(t *testing.T) {
	raw := []byte(`"$.a.b.c"`)
	var pathstr Path
	err := json.Unmarshal(raw, &pathstr)
	assert.NoError(t, err)
	assert.Equal(t, len(pathstr.path), 3)
	assert.Equal(t, pathstr.path[0], "a")
	assert.Equal(t, pathstr.path[1], "b")
	assert.Equal(t, pathstr.path[2], "c")
}

func TestParsePathString(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []string
		hasError bool
	}{
		{
			`Parse $`,
			`$`,
			[]string{},
			false,
		},
		{
			"Parse $.a.b.c",
			"$.a.b.c",
			[]string{"a", "b", "c"},
			false,
		},
		{
			`Parse $.['a'].['b'].['c']`,
			`$.['a'].['b'].['c']`,
			[]string{"a", "b", "c"},
			false,
		},
		{
			`Parse $['a']['b']['c']`,
			`$['a']['b']['c']`,
			[]string{"a", "b", "c"},
			false,
		},
		{
			`Parse $.["a"].["b"].["c"]`,
			`$.["a"].["b"].["c"]`,
			[]string{"a", "b", "c"},
			false,
		},
		{
			`Parse $["a"]["b"]["c"]`,
			`$["a"]["b"]["c"]`,
			[]string{"a", "b", "c"},
			false,
		},
		{
			"Parse $.a.[0]",
			"$.a.[0]",
			[]string{"a", "[0]"},
			false,
		},
		{
			"Parse $.a[0]",
			"$.a[0]",
			[]string{"a", "[0]"},
			false,
		},
		{
			`Parse $.['a'].[0]`,
			`$.['a'].[0]`,
			[]string{"a", "[0]"},
			false,
		},
		{
			`Parse $['a'][0]`,
			`$['a'][0]`,
			[]string{"a", "[0]"},
			false,
		},
		{
			`Parse $.["a"].[0]`,
			`$.["a"].[0]`,
			[]string{"a", "[0]"},
			false,
		},
		{
			`Parse $["a"][0]`,
			`$["a"][0]`,
			[]string{"a", "[0]"},
			false,
		},
		{
			`Parse $['a"b"c'][0]`,
			`$['a"b"c'][0]`,
			[]string{`a"b"c`, "[0]"},
			false,
		},
		{
			`Parse $["a'b'c"][0]`,
			`$["a'b'c"][0]`,
			[]string{`a'b'c`, "[0]"},
			false,
		},
		{
			"Parse $.data.[1]",
			"$.data.[1]",
			[]string{"data", "[1]"},
			false,
		},
		{
			"Parse $.data.[1,2,3]",
			"$.data.[1,2,3]",
			[]string{"data", "[1,2,3]"},
			false,
		},
		{
			"Parse $.data.[ 1   , 2  ,  3  ]",
			"$.data.[ 1   , 2  ,  3  ]",
			[]string{"data", "[1   , 2  ,  3  ]"},
			false,
		},
		{
			"Parse $.data.[1:3]",
			"$.data.[1:3]",
			[]string{"data", "[1:3]"},
			false,
		},
		{
			"Parse $.data.[1:2].values.[0].id",
			"$.data.[1:2].values.[0].id",
			[]string{"data", "[1:2]", "values", "[0]", "id"},
			false,
		},
		{
			`Parse $.data.[1:"IDDQD"].id`,
			`$.data.[1:"IDDQD"].id`,
			nil,
			true,
		},
		{
			`Parse $.data.[1:'IDDQD'].id`,
			`$.data.[1:'IDDQD'].id`,
			nil,
			true,
		},
		{
			`Parse $.data.[$.value].id`,
			`$.data.[$.value].id`,
			nil,
			true,
		},
		{
			`Parse $.`,
			`$.`,
			nil,
			true,
		},
		{
			`Parse $.value.`,
			`$.value.`,
			nil,
			true,
		},
		{
			`Parse .value`,
			`.value`,
			nil,
			true,
		},
		{
			`Parse value`,
			`value`,
			nil,
			true,
		},
		{
			`Parse ['value']`,
			`['value']`,
			nil,
			true,
		},
		{
			`Parse $.value..data`,
			`$.value..data`,
			nil,
			true,
		},
		{
			`Parse $.data.[1:true].id`,
			`$.data.[1:true].id`,
			[]string{"data", "[1:true]", "id"},
			false,
		},
		{
			`Parse $.data.[1 2; 3].id`,
			`$.data.[1 2; 3].id`,
			[]string{"data", "[1 2; 3]", "id"},
			false,
		},
	}
	for _, testCase := range tests {
		currentCase := testCase
		t.Run(currentCase.name, func(t *testing.T) {
			actual, err := parsePathString(currentCase.source)
			if currentCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentCase.expected, actual)
			}
		})
	}
}
