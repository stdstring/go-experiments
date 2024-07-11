package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPathNotFound(t *testing.T) {
	test := map[string]interface{}{}
	path, err := NewPath("$.a")
	assert.NoError(t, err)
	_, err = path.Get(test)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "not found")
}

func TestJSONPathGetDefault(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	path, err := NewPath("$")
	assert.NoError(t, err)
	out, err := path.Get(test)
	assert.NoError(t, err)
	assert.Equal(t, out, test)
}

func TestJSONPathGetSimple(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	path, err := NewPath("$.a")
	assert.NoError(t, err)
	out, err := path.Get(test)
	assert.NoError(t, err)
	assert.Equal(t, out, "b")
}

func TestJSONPathGetDeep(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x.a")
	assert.NoError(t, err)
	out, err := path.Get(outer)
	assert.NoError(t, err)
	assert.Equal(t, out, "b")
}

func TestJSONPathGetMap(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x")
	assert.NoError(t, err)
	out, err := path.GetMap(outer)
	assert.NoError(t, err)
	assert.Equal(t, out, test)
}

func TestJSONPathGetMapError(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x.a")
	assert.NoError(t, err)
	_, err = path.GetMap(outer)
	assert.Equal(t, err.Error(), "GetMap Error: must return map")
}

func TestJSONPathGetTime(t *testing.T) {
	const test = "2006-01-02T15:04:05Z"
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x")
	assert.NoError(t, err)
	out, err := path.GetTime(outer)
	assert.NoError(t, err)
	assert.Equal(t, out.Year(), 2006)
}

func TestJSONPathGetBool(t *testing.T) {
	const test = true
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x")
	assert.NoError(t, err)
	out, err := path.GetBool(outer)
	assert.NoError(t, err)
	assert.Equal(t, *out, test)
}

func TestJSONPathGetNumber(t *testing.T) {
	const test = 1.2
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x")
	assert.NoError(t, err)
	out, err := path.GetNumber(outer)
	assert.NoError(t, err)
	assert.Equal(t, *out, test)
}

func TestJSONPathGetString(t *testing.T) {
	const test = "String"
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x")
	assert.NoError(t, err)
	out, err := path.GetString(outer)
	assert.NoError(t, err)
	assert.Equal(t, *out, test)
}

func TestJSONPathGetSlice(t *testing.T) {
	test := []interface{}{1, 2, 3}
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x")
	assert.NoError(t, err)
	out, err := path.GetSlice(outer)
	assert.NoError(t, err)
	assert.Equal(t, out, test)
}

func createObjects(ids ...int) []interface{} {
	dest := make([]interface{}, len(ids))
	for index, id := range ids {
		dest[index] = map[string]interface{}{"id": id}
	}
	return dest
}

func TestJSONPathGetSingleValueFromOuterArray(t *testing.T) {
	simpleValues := []interface{}{666, 667, 668, 669}
	simpleObjects := createObjects(666, 667, 668, 669)
	container := map[string]interface{}{"data": simpleObjects}
	tests := []struct {
		name     string
		source   string
		data     interface{}
		expected int
		hasError bool
	}{
		{
			"Get value by path $.[0]",
			"$.[0]",
			simpleValues,
			666,
			false,
		},
		{
			"Get value by path $.[1]",
			"$.[1]",
			simpleValues,
			667,
			false,
		},
		{
			"Get value by path $.[3]",
			"$.[3]",
			simpleValues,
			669,
			false,
		},
		{
			"Get value by path $.[100]",
			"$.[100]",
			simpleValues,
			0,
			true,
		},
		{
			"Get value by path $.[-1]",
			"$.[-1]",
			simpleValues,
			669,
			false,
		},
		{
			"Get value by path $.[-2]",
			"$.[-2]",
			simpleValues,
			668,
			false,
		},
		{
			"Get value by path $.[-4]",
			"$.[-4]",
			simpleValues,
			666,
			false,
		},
		{
			"Get value by path $.[-100]",
			"$.[-100]",
			simpleValues,
			0,
			true,
		},
		{
			`Get value by path $.[true]`,
			`$.[true]`,
			simpleValues,
			0,
			true,
		},
		{
			"Get value by path $.[0].id",
			"$.[0].id",
			simpleObjects,
			666,
			false,
		},
		{
			"Get value by path $.[1].id",
			"$.[1].id",
			simpleObjects,
			667,
			false,
		},
		{
			"Get value by path $.[3].id",
			"$.[3].id",
			simpleObjects,
			669,
			false,
		},
		{
			"Get value by path $.[100].id",
			"$.[100].id",
			simpleObjects,
			0,
			true,
		},
		{
			"Get value by path $.[-1].id",
			"$.[-1].id",
			simpleObjects,
			669,
			false,
		},
		{
			"Get value by path $.[-2].id",
			"$.[-2].id",
			simpleObjects,
			668,
			false,
		},
		{
			"Get value by path $.[-4].id",
			"$.[-4].id",
			simpleObjects,
			666,
			false,
		},
		{
			"Get value by path $.[-100].id",
			"$.[-100].id",
			simpleObjects,
			0,
			true,
		},
		{
			`Get value by path $.[true].id`,
			`$.[true].id`,
			simpleObjects,
			0,
			true,
		},
		{
			"Get value by path $.data.[0].id",
			"$.data.[0].id",
			container,
			666,
			false,
		},
		{
			"Get value by path $.data.[1].id",
			"$.data.[1].id",
			container,
			667,
			false,
		},
		{
			"Get value by path $.data.[3].id",
			"$.data.[3].id",
			container,
			669,
			false,
		},
		{
			"Get value by path $.data.[100].id",
			"$.data.[100].id",
			container,
			0,
			true,
		},
		{
			"Get value by path $.data.[-1].id",
			"$.data.[-1].id",
			container,
			669,
			false,
		},
		{
			"Get value by path $.data.[-2].id",
			"$.data.[-2].id",
			container,
			668,
			false,
		},
		{
			"Get value by path $.data.[-4].id",
			"$.data.[-4].id",
			container,
			666,
			false,
		},
		{
			"Get value by path $.data.[-100].id",
			"$.data.[-100].id",
			container,
			0,
			true,
		},
		{
			`Get value by path $.data.[true].id`,
			`$.data.[true].id`,
			container,
			0,
			true,
		},
	}
	for _, testCase := range tests {
		currentCase := testCase
		t.Run(currentCase.name, func(t *testing.T) {
			path, err := NewPath(currentCase.source)
			assert.NoError(t, err)
			actual, err := path.Get(currentCase.data)
			if currentCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentCase.expected, actual)
			}
		})
	}
}

func TestJSONPathGetSeveralValuesFromOuterArray(t *testing.T) {
	simpleValues := []interface{}{666, 667, 668, 669}
	simpleObjects := createObjects(666, 667, 668, 669)
	container := map[string]interface{}{"data": simpleObjects}
	tests := []struct {
		name     string
		source   string
		data     interface{}
		expected []interface{}
		hasError bool
	}{
		{
			"Get values by path $.[0, 1]",
			"$.[0, 1]",
			simpleValues,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[0, 3]",
			"$.[0, 3]",
			simpleValues,
			[]interface{}{666, 669},
			false,
		},
		{
			"Get values by path $.[0, 0]",
			"$.[0, 0]",
			simpleValues,
			[]interface{}{666, 666},
			false,
		},
		{
			"Get values by path $.[0, 100]",
			"$.[0, 100]",
			simpleValues,
			[]interface{}{666},
			false,
		},
		{
			"Get values by path $.[100, 111]",
			"$.[100, 111]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[0, 100, 1]",
			"$.[0, 100, 1]",
			simpleValues,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[-1, -2]",
			"$.[-1, -2]",
			simpleValues,
			[]interface{}{669, 668},
			false,
		},
		{
			"Get values by path $.[-1, -4]",
			"$.[-1, -4]",
			simpleValues,
			[]interface{}{669, 666},
			false,
		},
		{
			"Get values by path $.[-1, -1]",
			"$.[-1, -1]",
			simpleValues,
			[]interface{}{669, 669},
			false,
		},
		{
			"Get values by path $.[-1, -100]",
			"$.[-1, -100]",
			simpleValues,
			[]interface{}{669},
			false,
		},
		{
			"Get values by path $.[-100, -111]",
			"$.[-100, -111]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-1, -100, -2]",
			"$.[-1, -100, -2]",
			simpleValues,
			[]interface{}{669, 668},
			false,
		},
		{
			`Get values by path $.[-1, true]`,
			`$.[1, true]`,
			simpleValues,
			[]interface{}{},
			true,
		},
		{
			"Get values by path $.[0, 1].id",
			"$.[0, 1].id",
			simpleObjects,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[0, 3].id",
			"$.[0, 3].id",
			simpleObjects,
			[]interface{}{666, 669},
			false,
		},
		{
			"Get values by path $.[0, 0].id",
			"$.[0, 0].id",
			simpleObjects,
			[]interface{}{666, 666},
			false,
		},
		{
			"Get values by path $.[0, 100].id",
			"$.[0, 100].id",
			simpleObjects,
			[]interface{}{666},
			false,
		},
		{
			"Get values by path $.[100, 111].id",
			"$.[100, 111].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[0, 100, 1].id",
			"$.[0, 100, 1].id",
			simpleObjects,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[-1, -2].id",
			"$.[-1, -2].id",
			simpleObjects,
			[]interface{}{669, 668},
			false,
		},
		{
			"Get values by path $.[-1, -4].id",
			"$.[-1, -4].id",
			simpleObjects,
			[]interface{}{669, 666},
			false,
		},
		{
			"Get values by path $.[-1, -1].id",
			"$.[-1, -1].id",
			simpleObjects,
			[]interface{}{669, 669},
			false,
		},
		{
			"Get values by path $.[-1, -100].id",
			"$.[-1, -100].id",
			simpleObjects,
			[]interface{}{669},
			false,
		},
		{
			"Get values by path $.[-100, -111].id",
			"$.[-100, -111].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-1, -100, -2].id",
			"$.[-1, -100, -2].id",
			simpleObjects,
			[]interface{}{669, 668},
			false,
		},
		{
			`Get values by path $.[-1, true].id`,
			`$.[-1, true].id`,
			simpleObjects,
			[]interface{}{},
			true,
		},
		{
			"Get values by path $.data.[0, 1].id",
			"$.data.[0, 1].id",
			container,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.data.[0, 3].id",
			"$.data.[0, 3].id",
			container,
			[]interface{}{666, 669},
			false,
		},
		{
			"Get values by path $.data.[0, 0].id",
			"$.data.[0, 0].id",
			container,
			[]interface{}{666, 666},
			false,
		},
		{
			"Get values by path $.data.[0, 100].id",
			"$.data.[0, 100].id",
			container,
			[]interface{}{666},
			false,
		},
		{
			"Get values by path $.data.[100, 111].id",
			"$.data.[100, 111].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[0, 100, 1].id",
			"$.data.[0, 100, 1].id",
			container,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.data.[-1, -2].id",
			"$.data.[-1, -2].id",
			container,
			[]interface{}{669, 668},
			false,
		},
		{
			"Get values by path $.data.[-1, -4].id",
			"$.data.[-1, -4].id",
			container,
			[]interface{}{669, 666},
			false,
		},
		{
			"Get values by path $.data.[-1, -1].id",
			"$.data.[-1, -1].id",
			container,
			[]interface{}{669, 669},
			false,
		},
		{
			"Get values by path $.data.[-1, -100].id",
			"$.data.[-1, -100].id",
			container,
			[]interface{}{669},
			false,
		},
		{
			"Get values by path $.data.[-100, -111].id",
			"$.data.[-100, -111].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[-1, -100, -2].id",
			"$.data.[-1, -100, -2].id",
			container,
			[]interface{}{669, 668},
			false,
		},
		{
			`Get values by path $.data.[-1, true].id`,
			`$.data.[-1, true].id`,
			container,
			[]interface{}{},
			true,
		},
	}
	for _, testCase := range tests {
		currentCase := testCase
		t.Run(currentCase.name, func(t *testing.T) {
			path, err := NewPath(currentCase.source)
			assert.NoError(t, err)
			actual, err := path.Get(currentCase.data)
			if currentCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentCase.expected, actual)
			}
		})
	}
}

func TestJSONPathGetSliceValuesFromOuterArray(t *testing.T) {
	simpleValues := []interface{}{666, 667, 668, 669}
	simpleObjects := createObjects(666, 667, 668, 669)
	container := map[string]interface{}{"data": simpleObjects}
	tests := []struct {
		name     string
		source   string
		data     interface{}
		expected []interface{}
		hasError bool
	}{
		{
			"Get values by path $.[0:1]",
			"$.[0:1]",
			simpleValues,
			[]interface{}{666},
			false,
		},
		{
			"Get values by path $.[0:2]",
			"$.[0:2]",
			simpleValues,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[1:2]",
			"$.[1:2]",
			simpleValues,
			[]interface{}{667},
			false,
		},
		{
			"Get values by path $.[1:3]",
			"$.[1:3]",
			simpleValues,
			[]interface{}{667, 668},
			false,
		},
		{
			"Get values by path $.[1:100]",
			"$.[1:100]",
			simpleValues,
			[]interface{}{667, 668, 669},
			false,
		},
		{
			"Get values by path $.[100:200]",
			"$.[100:200]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[0:0]",
			"$.[0:0]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[3:1]",
			"$.[3:1]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[1:]",
			"$.[1:]",
			simpleValues,
			[]interface{}{667, 668, 669},
			false,
		},
		{
			"Get values by path $.[100:]",
			"$.[100:]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[:2]",
			"$.[:2]",
			simpleValues,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[:100]",
			"$.[:100]",
			simpleValues,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[-2:-1]",
			"$.[-2:-1]",
			simpleValues,
			[]interface{}{668},
			false,
		},
		{
			"Get values by path $.[-3:-1]",
			"$.[-3:-1]",
			simpleValues,
			[]interface{}{667, 668},
			false,
		},
		{
			"Get values by path $.[-3:-2]",
			"$.[-3:-2]",
			simpleValues,
			[]interface{}{667},
			false,
		},
		{
			"Get values by path $.[-100:-2]",
			"$.[-100:-2]",
			simpleValues,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[-200:-100]",
			"$.[-200:-100]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-1:-1]",
			"$.[-1:-1]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-1:-2]",
			"$.[-1:-2]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-2:]",
			"$.[-2:]",
			simpleValues,
			[]interface{}{668, 669},
			false,
		},
		{
			"Get values by path $.[-100:]",
			"$.[-100:]",
			simpleValues,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[:-2]",
			"$.[:-2]",
			simpleValues,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[:-100]",
			"$.[:-100]",
			simpleValues,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-2:1]",
			"$.[-2:1]",
			simpleValues,
			[]interface{}{668, 669, 666},
			false,
		},
		{
			"Get values by path $.[-100:1]",
			"$.[-100:1]",
			simpleValues,
			[]interface{}{666, 667, 668, 669, 666},
			false,
		},
		{
			"Get values by path $.[-1:100]",
			"$.[-1:100]",
			simpleValues,
			[]interface{}{669, 666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[-100:100]",
			"$.[-100:100]",
			simpleValues,
			[]interface{}{666, 667, 668, 669, 666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[:]",
			"$.[:]",
			simpleValues,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			`Get values by path $.[1:true]`,
			`$.[1:true]`,
			simpleValues,
			[]interface{}{},
			true,
		},
		{
			`Get values by path $.[true:2]`,
			`$.[true:2]`,
			simpleValues,
			[]interface{}{},
			true,
		},
		{
			`Get values by path $.[true:false]`,
			`$.[true:false]`,
			simpleValues,
			[]interface{}{},
			true,
		},
		{
			"Get values by path $.[0:1].id",
			"$.[0:1].id",
			simpleObjects,
			[]interface{}{666},
			false,
		},
		{
			"Get values by path $.[0:2].id",
			"$.[0:2].id",
			simpleObjects,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[1:2].id",
			"$.[1:2].id",
			simpleObjects,
			[]interface{}{667},
			false,
		},
		{
			"Get values by path $.[1:3].id",
			"$.[1:3].id",
			simpleObjects,
			[]interface{}{667, 668},
			false,
		},
		{
			"Get values by path $.[1:100].id",
			"$.[1:100].id",
			simpleObjects,
			[]interface{}{667, 668, 669},
			false,
		},
		{
			"Get values by path $.[100:200].id",
			"$.[100:200].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[0:0].id",
			"$.[0:0].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[3:1].id",
			"$.[3:1].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[1:].id",
			"$.[1:].id",
			simpleObjects,
			[]interface{}{667, 668, 669},
			false,
		},
		{
			"Get values by path $.[100:].id",
			"$.[100:].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[:2].id",
			"$.[:2].id",
			simpleObjects,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[:100].id",
			"$.[:100].id",
			simpleObjects,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[-2:-1].id",
			"$.[-2:-1].id",
			simpleObjects,
			[]interface{}{668},
			false,
		},
		{
			"Get values by path $.[-3:-1].id",
			"$.[-3:-1].id",
			simpleObjects,
			[]interface{}{667, 668},
			false,
		},
		{
			"Get values by path $.[-3:-2].id",
			"$.[-3:-2].id",
			simpleObjects,
			[]interface{}{667},
			false,
		},
		{
			"Get values by path $.[-100:-2].id",
			"$.[-100:-2].id",
			simpleObjects,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[-200:-100].id",
			"$.[-200:-100].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-1:-1].id",
			"$.[-1:-1].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-1:-2].id",
			"$.[-1:-2].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-2:].id",
			"$.[-2:].id",
			simpleObjects,
			[]interface{}{668, 669},
			false,
		},
		{
			"Get values by path $.[-100:].id",
			"$.[-100:].id",
			simpleObjects,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[:-2].id",
			"$.[:-2].id",
			simpleObjects,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.[:-100].id",
			"$.[:-100].id",
			simpleObjects,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.[-2:1].id",
			"$.[-2:1].id",
			simpleObjects,
			[]interface{}{668, 669, 666},
			false,
		},
		{
			"Get values by path $.[-100:1].id",
			"$.[-100:1].id",
			simpleObjects,
			[]interface{}{666, 667, 668, 669, 666},
			false,
		},
		{
			"Get values by path $.[-1:100].id",
			"$.[-1:100].id",
			simpleObjects,
			[]interface{}{669, 666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[-100:100].id",
			"$.[-100:100].id",
			simpleObjects,
			[]interface{}{666, 667, 668, 669, 666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.[:].id",
			"$.[:].id",
			simpleObjects,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			`Get values by path $.[1:true].id`,
			`$.[1:true].id`,
			simpleObjects,
			[]interface{}{},
			true,
		},
		{
			`Get values by path $.[true:2].id`,
			`$.[true:2].id`,
			simpleObjects,
			[]interface{}{},
			true,
		},
		{
			`Get values by path $.[true:false].id`,
			`$.[true:false].id`,
			simpleObjects,
			[]interface{}{},
			true,
		},
		{
			"Get values by path $.data.[0:1].id",
			"$.data.[0:1].id",
			container,
			[]interface{}{666},
			false,
		},
		{
			"Get values by path $.data.[0:2].id",
			"$.data.[0:2].id",
			container,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.data.[1:2].id",
			"$.data.[1:2].id",
			container,
			[]interface{}{667},
			false,
		},
		{
			"Get values by path $.data.[1:3].id",
			"$.data.[1:3].id",
			container,
			[]interface{}{667, 668},
			false,
		},
		{
			"Get values by path $.data.[1:100].id",
			"$.data.[1:100].id",
			container,
			[]interface{}{667, 668, 669},
			false,
		},
		{
			"Get values by path $.data.[100:200].id",
			"$.data.[100:200].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[0:0].id",
			"$.data.[0:0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[3:1].id",
			"$.data.[3:1].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[1:].id",
			"$.data.[1:].id",
			container,
			[]interface{}{667, 668, 669},
			false,
		},
		{
			"Get values by path $.data.[100:].id",
			"$.data.[100:].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[:2].id",
			"$.data.[:2].id",
			container,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.data.[:100].id",
			"$.data.[:100].id",
			container,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.data.[-2:-1].id",
			"$.data.[-2:-1].id",
			container,
			[]interface{}{668},
			false,
		},
		{
			"Get values by path $.data.[-3:-1].id",
			"$.data.[-3:-1].id",
			container,
			[]interface{}{667, 668},
			false,
		},
		{
			"Get values by path $.data.[-3:-2].id",
			"$.data.[-3:-2].id",
			container,
			[]interface{}{667},
			false,
		},
		{
			"Get values by path $.data.[-100:-2].id",
			"$.data.[-100:-2].id",
			container,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.data.[-200:-100].id",
			"$.data.[-200:-100].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[-1:-1].id",
			"$.data.[-1:-1].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[-1:-2].id",
			"$.data.[-1:-2].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[-2:].id",
			"$.data.[-2:].id",
			container,
			[]interface{}{668, 669},
			false,
		},
		{
			"Get values by path $.data.[-100:].id",
			"$.data.[-100:].id",
			container,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.data.[:-2].id",
			"$.data.[:-2].id",
			container,
			[]interface{}{666, 667},
			false,
		},
		{
			"Get values by path $.data.[:-100].id",
			"$.data.[:-100].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get values by path $.data.[-2:1].id",
			"$.data.[-2:1].id",
			container,
			[]interface{}{668, 669, 666},
			false,
		},
		{
			"Get values by path $.data.[-100:1].id",
			"$.data.[-100:1].id",
			container,
			[]interface{}{666, 667, 668, 669, 666},
			false,
		},
		{
			"Get values by path $.data.[-1:100].id",
			"$.data.[-1:100].id",
			container,
			[]interface{}{669, 666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.data.[-100:100].id",
			"$.data.[-100:100].id",
			container,
			[]interface{}{666, 667, 668, 669, 666, 667, 668, 669},
			false,
		},
		{
			"Get values by path $.data.[:].id",
			"$.data.[:].id",
			container,
			[]interface{}{666, 667, 668, 669},
			false,
		},
		{
			`Get values by path $.data.[1:true].id`,
			`$.data.[1:true].id`,
			container,
			[]interface{}{},
			true,
		},
		{
			`Get values by path $.data.[true:2].id`,
			`$.data.[true:2].id`,
			container,
			[]interface{}{},
			true,
		},
		{
			`Get values by path $.data.[true:false].id`,
			`$.data.[true:false].id`,
			container,
			[]interface{}{},
			true,
		},
	}
	for _, testCase := range tests {
		currentCase := testCase
		t.Run(currentCase.name, func(t *testing.T) {
			path, err := NewPath(currentCase.source)
			assert.NoError(t, err)
			actual, err := path.Get(currentCase.data)
			if currentCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentCase.expected, actual)
			}
		})
	}
}

func TestJSONPathGetSingleValueFromInnerArray(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"id": 666, "values": createObjects(13, 14, 15, 16)},
		map[string]interface{}{"id": 667, "values": createObjects(23, 24, 25, 26)},
		map[string]interface{}{"id": 668, "values": createObjects(33, 34, 35, 36)},
	}
	container := map[string]interface{}{"data": data}
	tests := []struct {
		name     string
		source   string
		data     interface{}
		expected interface{}
		hasError bool
	}{
		{
			"Get value by path $.[0].values.[0].id",
			"$.[0].values.[0].id",
			data,
			13,
			false,
		},
		{
			"Get value by path $.[0].values.[1].id",
			"$.[0].values.[1].id",
			data,
			14,
			false,
		},
		{
			"Get value by path $.[0].values.[3].id",
			"$.[0].values.[3].id",
			data,
			16,
			false,
		},
		{
			"Get value by path $.[0].values.[4].id",
			"$.[0].values.[4].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[0].values.[100].id",
			"$.[0].values.[100].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[0].values.[-1].id",
			"$.[0].values.[-1].id",
			data,
			16,
			false,
		},
		{
			"Get value by path $.[0].values.[-2].id",
			"$.[0].values.[-2].id",
			data,
			15,
			false,
		},
		{
			"Get value by path $.[0].values.[-4].id",
			"$.[0].values.[-4].id",
			data,
			13,
			false,
		},
		{
			"Get value by path $.[0].values.[-5].id",
			"$.[0].values.[-5].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[0].values.[-100].id",
			"$.[0].values.[-100].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[0].id",
			"$.[1].values.[0].id",
			data,
			23,
			false,
		},
		{
			"Get value by path $.[1].values.[1].id",
			"$.[1].values.[1].id",
			data,
			24,
			false,
		},
		{
			"Get value by path $.[1].values.[3].id",
			"$.[1].values.[3].id",
			data,
			26,
			false,
		},
		{
			"Get value by path $.[1].values.[4].id",
			"$.[1].values.[4].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[100].id",
			"$.[1].values.[100].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[-1].id",
			"$.[1].values.[-1].id",
			data,
			26,
			false,
		},
		{
			"Get value by path $.[1].values.[-2].id",
			"$.[1].values.[-2].id",
			data,
			25,
			false,
		},
		{
			"Get value by path $.[1].values.[-4].id",
			"$.[1].values.[-4].id",
			data,
			23,
			false,
		},
		{
			"Get value by path $.[1].values.[-5].id",
			"$.[1].values.[-5].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[-100].id",
			"$.[1].values.[-100].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[0].id",
			"$.[-1].values.[0].id",
			data,
			33,
			false,
		},
		{
			"Get value by path $.[-1].values.[1].id",
			"$.[-1].values.[1].id",
			data,
			34,
			false,
		},
		{
			"Get value by path $.[-1].values.[3].id",
			"$.[-1].values.[3].id",
			data,
			36,
			false,
		},
		{
			"Get value by path $.[-1].values.[4].id",
			"$.[-1].values.[4].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[100].id",
			"$.[-1].values.[100].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[-1].id",
			"$.[-1].values.[-1].id",
			data,
			36,
			false,
		},
		{
			"Get value by path $.[-1].values.[-2].id",
			"$.[-1].values.[-2].id",
			data,
			35,
			false,
		},
		{
			"Get value by path $.[-1].values.[-4].id",
			"$.[-1].values.[-4].id",
			data,
			33,
			false,
		},
		{
			"Get value by path $.[-1].values.[-5].id",
			"$.[-1].values.[-5].id",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[-100].id",
			"$.[-1].values.[-100].id",
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[0].values.[true].id`,
			`$.[0].values.[true].id`,
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[true].values.[0].id`,
			`$.[true].values.[0].id`,
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[true].values.[false].id`,
			`$.[true].values.[false].id`,
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[0].values.[0]",
			"$.[0].values.[0]",
			data,
			map[string]interface{}{"id": 13},
			false,
		},
		{
			"Get value by path $.[0].values.[1]",
			"$.[0].values.[1]",
			data,
			map[string]interface{}{"id": 14},
			false,
		},
		{
			"Get value by path $.[0].values.[3]",
			"$.[0].values.[3]",
			data,
			map[string]interface{}{"id": 16},
			false,
		},
		{
			"Get value by path $.[0].values.[4]",
			"$.[0].values.[4]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[0].values.[100]",
			"$.[0].values.[100]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[0].values.[-1]",
			"$.[0].values.[-1]",
			data,
			map[string]interface{}{"id": 16},
			false,
		},
		{
			"Get value by path $.[0].values.[-2]",
			"$.[0].values.[-2]",
			data,
			map[string]interface{}{"id": 15},
			false,
		},
		{
			"Get value by path $.[0].values.[-4]",
			"$.[0].values.[-4]",
			data,
			map[string]interface{}{"id": 13},
			false,
		},
		{
			"Get value by path $.[0].values.[-5]",
			"$.[0].values.[-5]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[0].values.[-100]",
			"$.[0].values.[-100]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[0]",
			"$.[1].values.[0]",
			data,
			map[string]interface{}{"id": 23},
			false,
		},
		{
			"Get value by path $.[1].values.[1]",
			"$.[1].values.[1]",
			data,
			map[string]interface{}{"id": 24},
			false,
		},
		{
			"Get value by path $.[1].values.[3]",
			"$.[1].values.[3]",
			data,
			map[string]interface{}{"id": 26},
			false,
		},
		{
			"Get value by path $.[1].values.[4]",
			"$.[1].values.[4]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[100]",
			"$.[1].values.[100]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[-1]",
			"$.[1].values.[-1]",
			data,
			map[string]interface{}{"id": 26},
			false,
		},
		{
			"Get value by path $.[1].values.[-2]",
			"$.[1].values.[-2]",
			data,
			map[string]interface{}{"id": 25},
			false,
		},
		{
			"Get value by path $.[1].values.[-4]",
			"$.[1].values.[-4]",
			data,
			map[string]interface{}{"id": 23},
			false,
		},
		{
			"Get value by path $.[1].values.[-5]",
			"$.[1].values.[-5]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[1].values.[-100]",
			"$.[1].values.[-100]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[0]",
			"$.[-1].values.[0]",
			data,
			map[string]interface{}{"id": 33},
			false,
		},
		{
			"Get value by path $.[-1].values.[1]",
			"$.[-1].values.[1]",
			data,
			map[string]interface{}{"id": 34},
			false,
		},
		{
			"Get value by path $.[-1].values.[3]",
			"$.[-1].values.[3]",
			data,
			map[string]interface{}{"id": 36},
			false,
		},
		{
			"Get value by path $.[-1].values.[4]",
			"$.[-1].values.[4]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[100]",
			"$.[-1].values.[100]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[-1]",
			"$.[-1].values.[-1]",
			data,
			map[string]interface{}{"id": 36},
			false,
		},
		{
			"Get value by path $.[-1].values.[-2]",
			"$.[-1].values.[-2]",
			data,
			map[string]interface{}{"id": 35},
			false,
		},
		{
			"Get value by path $.[-1].values.[-4]",
			"$.[-1].values.[-4]",
			data,
			map[string]interface{}{"id": 33},
			false,
		},
		{
			"Get value by path $.[-1].values.[-5]",
			"$.[-1].values.[-5]",
			data,
			nil,
			true,
		},
		{
			"Get value by path $.[-1].values.[-100]",
			"$.[-1].values.[-100]",
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[0].values.[true]`,
			`$.[0].values.[true]`,
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[true].values.[0]`,
			`$.[true].values.[0]`,
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[true].values.[false]`,
			`$.[true].values.[false]`,
			data,
			nil,
			true,
		},
		{
			"Get value by path $.data.[0].values.[0]",
			"$.data.[0].values.[0]",
			container,
			map[string]interface{}{"id": 13},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1]",
			"$.data.[0].values.[1]",
			container,
			map[string]interface{}{"id": 14},
			false,
		},
		{
			"Get value by path $.data.[0].values.[3]",
			"$.data.[0].values.[3]",
			container,
			map[string]interface{}{"id": 16},
			false,
		},
		{
			"Get value by path $.data.[0].values.[4]",
			"$.data.[0].values.[4]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[0].values.[100]",
			"$.data.[0].values.[100]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[0].values.[-1]",
			"$.data.[0].values.[-1]",
			container,
			map[string]interface{}{"id": 16},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-2]",
			"$.data.[0].values.[-2]",
			container,
			map[string]interface{}{"id": 15},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-4]",
			"$.data.[0].values.[-4]",
			container,
			map[string]interface{}{"id": 13},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-5]",
			"$.data.[0].values.[-5]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[0].values.[-100]",
			"$.data.[0].values.[-100]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[1].values.[0]",
			"$.data.[1].values.[0]",
			container,
			map[string]interface{}{"id": 23},
			false,
		},
		{
			"Get value by path $.data.[1].values.[1]",
			"$.data.[1].values.[1]",
			container,
			map[string]interface{}{"id": 24},
			false,
		},
		{
			"Get value by path $.data.[1].values.[3]",
			"$.data.[1].values.[3]",
			container,
			map[string]interface{}{"id": 26},
			false,
		},
		{
			"Get value by path $.data.[1].values.[4]",
			"$.data.[1].values.[4]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[1].values.[100]",
			"$.data.[1].values.[100]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[1].values.[-1]",
			"$.data.[1].values.[-1]",
			container,
			map[string]interface{}{"id": 26},
			false,
		},
		{
			"Get value by path $.data.[1].values.[-2]",
			"$.data.[1].values.[-2]",
			container,
			map[string]interface{}{"id": 25},
			false,
		},
		{
			"Get value by path $.data.[1].values.[-4]",
			"$.data.[1].values.[-4]",
			container,
			map[string]interface{}{"id": 23},
			false,
		},
		{
			"Get value by path $.data.[1].values.[-5]",
			"$.data.[1].values.[-5]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[1].values.[-100]",
			"$.data.[1].values.[-100]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[-1].values.[0]",
			"$.data.[-1].values.[0]",
			container,
			map[string]interface{}{"id": 33},
			false,
		},
		{
			"Get value by path $.data.[-1].values.[1]",
			"$.data.[-1].values.[1]",
			container,
			map[string]interface{}{"id": 34},
			false,
		},
		{
			"Get value by path $.data.[-1].values.[3]",
			"$.data.[-1].values.[3]",
			container,
			map[string]interface{}{"id": 36},
			false,
		},
		{
			"Get value by path $.data.[-1].values.[4]",
			"$.data.[-1].values.[4]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[-1].values.[100]",
			"$.data.[-1].values.[100]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[-1].values.[-1]",
			"$.data.[-1].values.[-1]",
			container,
			map[string]interface{}{"id": 36},
			false,
		},
		{
			"Get value by path $.data.[-1].values.[-2]",
			"$.data.[-1].values.[-2]",
			container,
			map[string]interface{}{"id": 35},
			false,
		},
		{
			"Get value by path $.data.[-1].values.[-4]",
			"$.data.[-1].values.[-4]",
			container,
			map[string]interface{}{"id": 33},
			false,
		},
		{
			"Get value by path $.data.[-1].values.[-5]",
			"$.data.[-1].values.[-5]",
			container,
			nil,
			true,
		},
		{
			"Get value by path $.data.[-1].values.[-100]",
			"$.data.[-1].values.[-100]",
			container,
			nil,
			true,
		},
		{
			`Get value by path $.data.[0].values.[true]`,
			`$.data.[0].values.[true]`,
			container,
			nil,
			true,
		},
		{
			`Get value by path $.data.[true].values.[0]`,
			`$.data.[true].values.[0]`,
			container,
			nil,
			true,
		},
		{
			`Get value by path $.data.[true].values.[false]`,
			`$.data.[true].values.[false]`,
			container,
			nil,
			true,
		},
	}
	for _, testCase := range tests {
		currentCase := testCase
		t.Run(currentCase.name, func(t *testing.T) {
			path, err := NewPath(currentCase.source)
			assert.NoError(t, err)
			actual, err := path.Get(currentCase.data)
			if currentCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentCase.expected, actual)
			}
		})
	}
}

func TestJSONPathGetSeveralValuesFromInnerArray(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"id": 666, "values": createObjects(13, 14, 15, 16)},
		map[string]interface{}{"id": 667, "values": createObjects(23, 24, 25, 26)},
		map[string]interface{}{"id": 668, "values": createObjects(33, 34, 35, 36)},
	}
	container := map[string]interface{}{"data": data}
	tests := []struct {
		name     string
		source   string
		data     interface{}
		expected []interface{}
		hasError bool
	}{
		{
			"Get value by path $.[0].values.[0, 1].id",
			"$.[0].values.[0, 1].id",
			data,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.[0].values.[0, 2].id",
			"$.[0].values.[0, 2].id",
			data,
			[]interface{}{13, 15},
			false,
		},
		{
			"Get value by path $.[0].values.[1, 3].id",
			"$.[0].values.[1, 3].id",
			data,
			[]interface{}{14, 16},
			false,
		},
		{
			"Get value by path $.[0].values.[1, 4].id",
			"$.[0].values.[1, 4].id",
			data,
			[]interface{}{14},
			false,
		},
		{
			"Get value by path $.[0].values.[1, 4].id",
			"$.[0].values.[100, 200].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[0, 0].id",
			"$.[0].values.[0, 0].id",
			data,
			[]interface{}{13, 13},
			false,
		},
		{
			"Get value by path $.[0].values.[0, -1].id",
			"$.[0].values.[0, -1].id",
			data,
			[]interface{}{13, 16},
			false,
		},
		{
			"Get value by path $.[0].values.[1, -3].id",
			"$.[0].values.[1, -3].id",
			data,
			[]interface{}{14, 14},
			false,
		},
		{
			"Get value by path $.[0].values.[100, -100].id",
			"$.[0].values.[100, -100].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[1, -1, 2, -2, 0].id",
			"$.[0].values.[1, -1, 2, -2, 0].id",
			data,
			[]interface{}{14, 16, 15, 15, 13},
			false,
		},
		{
			"Get value by path $.[0, 1].values.[0].id",
			"$.[0, 1].values.[0].id",
			data,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.[0, 2].values.[0].id",
			"$.[0, 2].values.[0].id",
			data,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.[0, 3].values.[0].id",
			"$.[0, 3].values.[0].id",
			data,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.[100, 200].values.[0].id",
			"$.[100, 200].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0, 0].values.[0].id",
			"$.[0, 0].values.[0].id",
			data,
			[]interface{}{13, 13},
			false,
		},
		{
			"Get value by path $.[0, -1].values.[0].id",
			"$.[0, -1].values.[0].id",
			data,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.[0, -2].values.[0].id",
			"$.[0, -2].values.[0].id",
			data,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.[0, -100].values.[0].id",
			"$.[0, -100].values.[0].id",
			data,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.[100, -1].values.[0].id",
			"$.[100, -1].values.[0].id",
			data,
			[]interface{}{33},
			false,
		},
		{
			"Get value by path $.[100, -100].values.[0].id",
			"$.[10, -100].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[1, -1, 2, -2, 100, 0].values.[0].id",
			"$.[1, -1, 2, -2, 100, 0].values.[0].id",
			data,
			[]interface{}{23, 33, 33, 23, 13},
			false,
		},
		{
			"Get value by path $.[0, 1].values.[0, 1].id",
			"$.[0, 1].values.[0, 1].id",
			data,
			[]interface{}{13, 14, 23, 24},
			false,
		},
		{
			"Get value by path $.[0, 1].values.[0, 2].id",
			"$.[0, 1].values.[0, 2].id",
			data,
			[]interface{}{13, 15, 23, 25},
			false,
		},
		{
			"Get value by path $.[0, 2].values.[0, 1].id",
			"$.[0, 2].values.[0, 1].id",
			data,
			[]interface{}{13, 14, 33, 34},
			false,
		},
		{
			"Get value by path $.[0, 2].values.[0, 2].id",
			"$.[0, 2].values.[0, 2].id",
			data,
			[]interface{}{13, 15, 33, 35},
			false,
		},
		{
			"Get value by path $.[0, 1].values.[0, 100].id",
			"$.[0, 1].values.[0, 100].id",
			data,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.[0, 100].values.[0, 1].id",
			"$.[0, 100].values.[0, 1].id",
			data,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.[0, 100].values.[0, 100].id",
			"$.[0, 100].values.[0, 100].id",
			data,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.[100, 200].values.[0, 1].id",
			"$.[100, 200].values.[0, 1].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0, 1].values.[100, 200].id",
			"$.[0, 1].values.[100, 200].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0, 0].values.[0, 0].id",
			"$.[0, 0].values.[0, 0].id",
			data,
			[]interface{}{13, 13, 13, 13},
			false,
		},
		{
			"Get value by path $.[0, -1].values.[0, 1].id",
			"$.[0, -1].values.[0, 1].id",
			data,
			[]interface{}{13, 14, 33, 34},
			false,
		},
		{
			"Get value by path $.[0, 1].values.[0, -1].id",
			"$.[0, 1].values.[0, -1].id",
			data,
			[]interface{}{13, 16, 23, 26},
			false,
		},
		{
			"Get value by path $.[0, -1].values.[0, -1].id",
			"$.[0, -1].values.[0, -1].id",
			data,
			[]interface{}{13, 16, 33, 36},
			false,
		},
		{
			"Get value by path $.[0, -1].values.[0, -100].id",
			"$.[0, -1].values.[0, -100].id",
			data,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.[0, -100].values.[0, -1].id",
			"$.[0, -100].values.[0, -1].id",
			data,
			[]interface{}{13, 16},
			false,
		},
		{
			"Get value by path $.[0, -100].values.[0, -100].id",
			"$.[0, -100].values.[0, -100].id",
			data,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.[0, -1, 100, 1].values.[0, -100, 2].id",
			"$.[0, -1, 100, 1].values.[0, -100, 2].id",
			data,
			[]interface{}{13, 15, 33, 35, 23, 25},
			false,
		},
		{
			`Get value by path $.[0, true].values.[0, 1].id`,
			`$.[0, true].values.[0, 1].id`,
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[0, 1].values.[0, true].id`,
			`$.[0, 1].values.[0, true].id`,
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[0, true].values.[0, false].id`,
			`$.[0, true].values.[0, false].id`,
			data,
			nil,
			true,
		},
		{
			"Get value by path $.data.[0].values.[0, 1].id",
			"$.data.[0].values.[0, 1].id",
			container,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.data.[0].values.[0, 2].id",
			"$.data.[0].values.[0, 2].id",
			container,
			[]interface{}{13, 15},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1, 3].id",
			"$.data.[0].values.[1, 3].id",
			container,
			[]interface{}{14, 16},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1, 4].id",
			"$.data.[0].values.[1, 4].id",
			container,
			[]interface{}{14},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1, 4].id",
			"$.data.[0].values.[100, 200].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[0, 0].id",
			"$.data.[0].values.[0, 0].id",
			container,
			[]interface{}{13, 13},
			false,
		},
		{
			"Get value by path $.data.[0].values.[0, -1].id",
			"$.data.[0].values.[0, -1].id",
			container,
			[]interface{}{13, 16},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1, -3].id",
			"$.data.[0].values.[1, -3].id",
			container,
			[]interface{}{14, 14},
			false,
		},
		{
			"Get value by path $.data.[0].values.[100, -100].id",
			"$.data.[0].values.[100, -100].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1, -1, 2, -2, 0].id",
			"$.data.[0].values.[1, -1, 2, -2, 0].id",
			container,
			[]interface{}{14, 16, 15, 15, 13},
			false,
		},
		{
			"Get value by path $.data.[0, 1].values.[0].id",
			"$.data.[0, 1].values.[0].id",
			container,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.data.[0, 2].values.[0].id",
			"$.data.[0, 2].values.[0].id",
			container,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.data.[0, 3].values.[0].id",
			"$.data.[0, 3].values.[0].id",
			container,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.data.[100, 200].values.[0].id",
			"$.data.[100, 200].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0, 0].values.[0].id",
			"$.data.[0, 0].values.[0].id",
			container,
			[]interface{}{13, 13},
			false,
		},
		{
			"Get value by path $.data.[0, -1].values.[0].id",
			"$.data.[0, -1].values.[0].id",
			container,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.data.[0, -2].values.[0].id",
			"$.data.[0, -2].values.[0].id",
			container,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.data.[0, -100].values.[0].id",
			"$.data.[0, -100].values.[0].id",
			container,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.data.[100, -1].values.[0].id",
			"$.data.[100, -1].values.[0].id",
			container,
			[]interface{}{33},
			false,
		},
		{
			"Get value by path $.data.[100, -100].values.[0].id",
			"$.data.[10, -100].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[1, -1, 2, -2, 100, 0].values.[0].id",
			"$.data.[1, -1, 2, -2, 100, 0].values.[0].id",
			container,
			[]interface{}{23, 33, 33, 23, 13},
			false,
		},
		{
			"Get value by path $.data.[0, 1].values.[0, 1].id",
			"$.data.[0, 1].values.[0, 1].id",
			container,
			[]interface{}{13, 14, 23, 24},
			false,
		},
		{
			"Get value by path $.data.[0, 1].values.[0, 2].id",
			"$.data.[0, 1].values.[0, 2].id",
			container,
			[]interface{}{13, 15, 23, 25},
			false,
		},
		{
			"Get value by path $.data.[0, 2].values.[0, 1].id",
			"$.data.[0, 2].values.[0, 1].id",
			container,
			[]interface{}{13, 14, 33, 34},
			false,
		},
		{
			"Get value by path $.data.[0, 2].values.[0, 2].id",
			"$.data.[0, 2].values.[0, 2].id",
			container,
			[]interface{}{13, 15, 33, 35},
			false,
		},
		{
			"Get value by path $.data.[0, 1].values.[0, 100].id",
			"$.data.[0, 1].values.[0, 100].id",
			container,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.data.[0, 100].values.[0, 1].id",
			"$.data.[0, 100].values.[0, 1].id",
			container,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.data.[0, 100].values.[0, 100].id",
			"$.data.[0, 100].values.[0, 100].id",
			container,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.data.[100, 200].values.[0, 1].id",
			"$.data.[100, 200].values.[0, 1].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0, 1].values.[100, 200].id",
			"$.data.[0, 1].values.[100, 200].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0, 0].values.[0, 0].id",
			"$.data.[0, 0].values.[0, 0].id",
			container,
			[]interface{}{13, 13, 13, 13},
			false,
		},
		{
			"Get value by path $.data.[0, -1].values.[0, 1].id",
			"$.data.[0, -1].values.[0, 1].id",
			container,
			[]interface{}{13, 14, 33, 34},
			false,
		},
		{
			"Get value by path $.data.[0, 1].values.[0, -1].id",
			"$.data.[0, 1].values.[0, -1].id",
			container,
			[]interface{}{13, 16, 23, 26},
			false,
		},
		{
			"Get value by path $.data.[0, -1].values.[0, -1].id",
			"$.data.[0, -1].values.[0, -1].id",
			container,
			[]interface{}{13, 16, 33, 36},
			false,
		},
		{
			"Get value by path $.data.[0, -1].values.[0, -100].id",
			"$.data.[0, -1].values.[0, -100].id",
			container,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.data.[0, -100].values.[0, -1].id",
			"$.data.[0, -100].values.[0, -1].id",
			container,
			[]interface{}{13, 16},
			false,
		},
		{
			"Get value by path $.data.[0, -100].values.[0, -100].id",
			"$.data.[0, -100].values.[0, -100].id",
			container,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.data.[0, -1, 100, 1].values.[0, -100, 2].id",
			"$.data.[0, -1, 100, 1].values.[0, -100, 2].id",
			container,
			[]interface{}{13, 15, 33, 35, 23, 25},
			false,
		},
		{
			`Get value by path $.data.[0, true].values.[0, 1].id`,
			`$.data.[0, true].values.[0, 1].id`,
			container,
			nil,
			true,
		},
		{
			`Get value by path $.data.[0, 1].values.[0, true].id`,
			`$.data.[0, 1].values.[0, true].id`,
			container,
			nil,
			true,
		},
		{
			`Get value by path $.data.[0, true].values.[0, false].id`,
			`$.data.[0, true].values.[0, false].id`,
			container,
			nil,
			true,
		},
	}
	for _, testCase := range tests {
		currentCase := testCase
		t.Run(currentCase.name, func(t *testing.T) {
			path, err := NewPath(currentCase.source)
			assert.NoError(t, err)
			actual, err := path.Get(currentCase.data)

			if currentCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentCase.expected, actual)
			}
		})
	}
}

func TestJSONPathGetSliceValuesFromInnerArray(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"id": 666, "values": createObjects(13, 14, 15, 16)},
		map[string]interface{}{"id": 667, "values": createObjects(23, 24, 25, 26)},
		map[string]interface{}{"id": 668, "values": createObjects(33, 34, 35, 36)},
	}
	container := map[string]interface{}{"data": data}
	tests := []struct {
		name     string
		source   string
		data     interface{}
		expected []interface{}
		hasError bool
	}{
		{
			"Get value by path $.[0].values.[0:1].id",
			"$.[0].values.[0:1].id",
			data,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.[0].values.[0:2].id",
			"$.[0].values.[0:2].id",
			data,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.[0].values.[1:3].id",
			"$.[0].values.[1:3].id",
			data,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.[0].values.[1:4].id",
			"$.[0].values.[1:4].id",
			data,
			[]interface{}{14, 15, 16},
			false,
		},
		{
			"Get value by path $.[0].values.[1:100].id",
			"$.[0].values.[1:100].id",
			data,
			[]interface{}{14, 15, 16},
			false,
		},
		{
			"Get value by path $.[0].values.[100:200].id",
			"$.[0].values.[100:200].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[0:0].id",
			"$.[0].values.[0:0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[2:1].id",
			"$.[0].values.[2:1].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[-2:-1].id",
			"$.[0].values.[-2:-1].id",
			data,
			[]interface{}{15},
			false,
		},
		{
			"Get value by path $.[0].values.[-3:-1].id",
			"$.[0].values.[-3:-1].id",
			data,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.[0].values.[-100:-1].id",
			"$.[0].values.[-100:-1].id",
			data,
			[]interface{}{13, 14, 15},
			false,
		},
		{
			"Get value by path $.[0].values.[-200:-100].id",
			"$.[0].values.[-200:-100].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[-1:-1].id",
			"$.[0].values.[-1:-1].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[-1:-2].id",
			"$.[0].values.[-1:-2].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0].values.[-2:1].id",
			"$.[0].values.[-2:1].id",
			data,
			[]interface{}{15, 16, 13},
			false,
		},
		{
			"Get value by path $.[0].values.[-100:100].id",
			"$.[0].values.[-100:100].id",
			data,
			[]interface{}{13, 14, 15, 16, 13, 14, 15, 16},
			false,
		},
		{
			"Get value by path $.[0:1].values.[0].id",
			"$.[0:1].values.[0].id",
			data,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.[0:2].values.[0].id",
			"$.[0:2].values.[0].id",
			data,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.[1:2].values.[0].id",
			"$.[1:2].values.[0].id",
			data,
			[]interface{}{23},
			false,
		},
		{
			"Get value by path $.[1:100].values.[0].id",
			"$.[1:100].values.[0].id",
			data,
			[]interface{}{23, 33},
			false,
		},
		{
			"Get value by path $.[100:200].values.[0].id",
			"$.[100:200].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[0:0].values.[0].id",
			"$.[0:0].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[2:1].values.[0].id",
			"$.[2:1].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[-2:-1].values.[0].id",
			"$.[-2:-1].values.[0].id",
			data,
			[]interface{}{23},
			false,
		},
		{
			"Get value by path $.[-3:-1].values.[0].id",
			"$.[-3:-1].values.[0].id",
			data,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.[-100:-1].values.[0].id",
			"$.[-100:-1].values.[0].id",
			data,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.[-200:-100].values.[0].id",
			"$.[-200:-100].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[-1:-1].values.[0].id",
			"$.[-1:-1].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[-1:-2].values.[0].id",
			"$.[-1:-2].values.[0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[-2:1].values.[0].id",
			"$.[-2:1].values.[0].id",
			data,
			[]interface{}{23, 33, 13},
			false,
		},
		{
			"Get value by path $.[-100:100].values.[0].id",
			"$.[-100:100].values.[0].id",
			data,
			[]interface{}{13, 23, 33, 13, 23, 33},
			false,
		},
		{
			"Get value by path $.[0, 2].values.[0:1].id",
			"$.[0, 2].values.[0:1].id",
			data,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.[0, 2].values.[1:3].id",
			"$.[0, 2].values.[1:3].id",
			data,
			[]interface{}{14, 15, 34, 35},
			false,
		},
		{
			"Get value by path $.[0, 0].values.[1:3].id",
			"$.[0, 0].values.[1:3].id",
			data,
			[]interface{}{14, 15, 14, 15},
			false,
		},
		{
			"Get value by path $.[0, 100].values.[1:3].id",
			"$.[0, 100].values.[1:3].id",
			data,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.[0, -1].values.[1:3].id",
			"$.[0, -1].values.[1:3].id",
			data,
			[]interface{}{14, 15, 34, 35},
			false,
		},
		{
			"Get value by path $.[0, -100].values.[1:3].id",
			"$.[0, -100].values.[1:3].id",
			data,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.[0, -100, -1, -100, 1].values.[1:3].id",
			"$.[0, -100, -1, -100, 1].values.[1:3].id",
			data,
			[]interface{}{14, 15, 34, 35, 24, 25},
			false,
		},
		{
			"Get value by path $.[0:1].values.[0, 2].id",
			"$.[0:1].values.[0, 2].id",
			data,
			[]interface{}{13, 15},
			false,
		},
		{
			"Get value by path $.[1:3].values.[0, 2].id",
			"$.[1:3].values.[0, 2].id",
			data,
			[]interface{}{23, 25, 33, 35},
			false,
		},
		{
			"Get value by path $.[1:3].values.[0, 0].id",
			"$.[1:3].values.[0, 0].id",
			data,
			[]interface{}{23, 23, 33, 33},
			false,
		},
		{
			"Get value by path $.[1:3].values.[0, 100].id",
			"$.[1:3].values.[0, 100].id",
			data,
			[]interface{}{23, 33},
			false,
		},
		{
			"Get value by path $.[1:3].values.[0, -1].id",
			"$.[1:3].values.[0, -1].id",
			data,
			[]interface{}{23, 26, 33, 36},
			false,
		},
		{
			"Get value by path $.[1:3].values.[0, -100].id",
			"$.[1:3].values.[0, -100].id",
			data,
			[]interface{}{23, 33},
			false,
		},
		{
			"Get value by path $.[1:3].values.[0, -100, -1, 100, 1].id",
			"$.[1:3].values.[0, -100, -1, 100, 1].id",
			data,
			[]interface{}{23, 26, 24, 33, 36, 34},
			false,
		},
		{
			"Get value by path $.[0:1].values.[0:2].id",
			"$.[0:1].values.[0:2].id",
			data,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.[1:3].values.[0:2].id",
			"$.[1:3].values.[0:2].id",
			data,
			[]interface{}{23, 24, 33, 34},
			false,
		},
		{
			"Get value by path $.[1:3].values.[2:4].id",
			"$.[1:3].values.[2:4].id",
			data,
			[]interface{}{25, 26, 35, 36},
			false,
		},
		{
			"Get value by path $.[1:3].values.[-4:-2].id",
			"$.[1:3].values.[-4:-2].id",
			data,
			[]interface{}{23, 24, 33, 34},
			false,
		},
		{
			"Get value by path $.[-3:-1].values.[2:4].id",
			"$.[-3:-1].values.[2:4].id",
			data,
			[]interface{}{15, 16, 25, 26},
			false,
		},
		{
			"Get value by path $.[-3:-1].values.[-4:-2].id",
			"$.[-3:-1].values.[-4:-2].id",
			data,
			[]interface{}{13, 14, 23, 24},
			false,
		},
		{
			"Get value by path $.[0:0].values.[1:4].id",
			"$.[0:0].values.[1:4].id",
			data,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.[1:4].values.[0:0].id",
			"$.[1:4].values.[0:0].id",
			data,
			[]interface{}{},
			false,
		},
		{
			`Get value by path $.[1:4].values.[0:true].id`,
			`$.[1:4].values.[0:true].id`,
			data,
			nil,
			true,
		},
		{
			`Get value by path $.[1:4].values.[true: 0].id`,
			`$.[1:4].values.[true: 0].id`,
			data,
			nil,
			true,
		},
		{
			"Get value by path $.data.[0].values.[0:1].id",
			"$.data.[0].values.[0:1].id",
			container,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.data.[0].values.[0:2].id",
			"$.data.[0].values.[0:2].id",
			container,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1:3].id",
			"$.data.[0].values.[1:3].id",
			container,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.data.[0].values.[1:4].id",
			"$.data.[0].values.[1:4].id",
			container,
			[]interface{}{14, 15, 16},
			false,
		},
		{
			"Get value by path $.[0].values.[1:100].id",
			"$.[0].values.[1:100].id",
			data,
			[]interface{}{14, 15, 16},
			false,
		},
		{
			"Get value by path $.data.[0].values.[100:200].id",
			"$.data.[0].values.[100:200].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[0:0].id",
			"$.data.[0].values.[0:0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[2:1].id",
			"$.data.[0].values.[2:1].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-2:-1].id",
			"$.data.[0].values.[-2:-1].id",
			container,
			[]interface{}{15},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-3:-1].id",
			"$.data.[0].values.[-3:-1].id",
			container,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-100:-1].id",
			"$.data.[0].values.[-100:-1].id",
			container,
			[]interface{}{13, 14, 15},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-200:-100].id",
			"$.data.[0].values.[-200:-100].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-1:-1].id",
			"$.data.[0].values.[-1:-1].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-1:-2].id",
			"$.data.[0].values.[-1:-2].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-2:1].id",
			"$.data.[0].values.[-2:1].id",
			container,
			[]interface{}{15, 16, 13},
			false,
		},
		{
			"Get value by path $.data.[0].values.[-100:100].id",
			"$.data.[0].values.[-100:100].id",
			container,
			[]interface{}{13, 14, 15, 16, 13, 14, 15, 16},
			false,
		},
		{
			"Get value by path $.data.[0:1].values.[0].id",
			"$.data.[0:1].values.[0].id",
			container,
			[]interface{}{13},
			false,
		},
		{
			"Get value by path $.data.[0:2].values.[0].id",
			"$.data.[0:2].values.[0].id",
			container,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.data.[1:2].values.[0].id",
			"$.data.[1:2].values.[0].id",
			container,
			[]interface{}{23},
			false,
		},
		{
			"Get value by path $.data.[1:100].values.[0].id",
			"$.data.[1:100].values.[0].id",
			container,
			[]interface{}{23, 33},
			false,
		},
		{
			"Get value by path $.data.[100:200].values.[0].id",
			"$.data.[100:200].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[0:0].values.[0].id",
			"$.data.[0:0].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[2:1].values.[0].id",
			"$.data.[2:1].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[-2:-1].values.[0].id",
			"$.data.[-2:-1].values.[0].id",
			container,
			[]interface{}{23},
			false,
		},
		{
			"Get value by path $.data.[-3:-1].values.[0].id",
			"$.data.[-3:-1].values.[0].id",
			container,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.data.[-100:-1].values.[0].id",
			"$.data.[-100:-1].values.[0].id",
			container,
			[]interface{}{13, 23},
			false,
		},
		{
			"Get value by path $.data.[-200:-100].values.[0].id",
			"$.data.[-200:-100].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[-1:-1].values.[0].id",
			"$.data.[-1:-1].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[-1:-2].values.[0].id",
			"$.data.[-1:-2].values.[0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[-2:1].values.[0].id",
			"$.data.[-2:1].values.[0].id",
			container,
			[]interface{}{23, 33, 13},
			false,
		},
		{
			"Get value by path $.data.[-100:100].values.[0].id",
			"$.data.[-100:100].values.[0].id",
			container,
			[]interface{}{13, 23, 33, 13, 23, 33},
			false,
		},
		{
			"Get value by path $.data.[0, 2].values.[0:1].id",
			"$.data.[0, 2].values.[0:1].id",
			container,
			[]interface{}{13, 33},
			false,
		},
		{
			"Get value by path $.data.[0, 2].values.[1:3].id",
			"$.data.[0, 2].values.[1:3].id",
			container,
			[]interface{}{14, 15, 34, 35},
			false,
		},
		{
			"Get value by path $.data.[0, 0].values.[1:3].id",
			"$.data.[0, 0].values.[1:3].id",
			container,
			[]interface{}{14, 15, 14, 15},
			false,
		},
		{
			"Get value by path $.data.[0, 100].values.[1:3].id",
			"$.data.[0, 100].values.[1:3].id",
			container,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.data.[0, -1].values.[1:3].id",
			"$.data.[0, -1].values.[1:3].id",
			container,
			[]interface{}{14, 15, 34, 35},
			false,
		},
		{
			"Get value by path $.data.[0, -100].values.[1:3].id",
			"$.data.[0, -100].values.[1:3].id",
			container,
			[]interface{}{14, 15},
			false,
		},
		{
			"Get value by path $.data.[0, -100, -1, -100, 1].values.[1:3].id",
			"$.data.[0, -100, -1, -100, 1].values.[1:3].id",
			container,
			[]interface{}{14, 15, 34, 35, 24, 25},
			false,
		},
		{
			"Get value by path $.data.[0:1].values.[0, 2].id",
			"$.data.[0:1].values.[0, 2].id",
			container,
			[]interface{}{13, 15},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[0, 2].id",
			"$.data.[1:3].values.[0, 2].id",
			container,
			[]interface{}{23, 25, 33, 35},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[0, 0].id",
			"$.data.[1:3].values.[0, 0].id",
			container,
			[]interface{}{23, 23, 33, 33},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[0, 100].id",
			"$.data.[1:3].values.[0, 100].id",
			container,
			[]interface{}{23, 33},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[0, -1].id",
			"$.data.[1:3].values.[0, -1].id",
			container,
			[]interface{}{23, 26, 33, 36},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[0, -100].id",
			"$.data.[1:3].values.[0, -100].id",
			container,
			[]interface{}{23, 33},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[0, -100, -1, 100, 1].id",
			"$.data.[1:3].values.[0, -100, -1, 100, 1].id",
			container,
			[]interface{}{23, 26, 24, 33, 36, 34},
			false,
		},
		{
			"Get value by path $.data.[0:1].values.[0:2].id",
			"$.data.[0:1].values.[0:2].id",
			container,
			[]interface{}{13, 14},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[0:2].id",
			"$.data.[1:3].values.[0:2].id",
			container,
			[]interface{}{23, 24, 33, 34},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[2:4].id",
			"$.data.[1:3].values.[2:4].id",
			container,
			[]interface{}{25, 26, 35, 36},
			false,
		},
		{
			"Get value by path $.data.[1:3].values.[-4:-2].id",
			"$.data.[1:3].values.[-4:-2].id",
			container,
			[]interface{}{23, 24, 33, 34},
			false,
		},
		{
			"Get value by path $.data.[-3:-1].values.[2:4].id",
			"$.data.[-3:-1].values.[2:4].id",
			container,
			[]interface{}{15, 16, 25, 26},
			false,
		},
		{
			"Get value by path $.data.[-3:-1].values.[-4:-2].id",
			"$.data.[-3:-1].values.[-4:-2].id",
			container,
			[]interface{}{13, 14, 23, 24},
			false,
		},
		{
			"Get value by path $.data.[0:0].values.[1:4].id",
			"$.data.[0:0].values.[1:4].id",
			container,
			[]interface{}{},
			false,
		},
		{
			"Get value by path $.data.[1:4].values.[0:0].id",
			"$.data.[1:4].values.[0:0].id",
			container,
			[]interface{}{},
			false,
		},
		{
			`Get value by path $.data.[1:4].values.[0:true].id`,
			`$.data.[1:4].values.[0:true].id`,
			container,
			nil,
			true,
		},
		{
			`Get value by path $.data.[1:4].values.[true: 0].id`,
			`$.data.[1:4].values.[true: 0].id`,
			container,
			nil,
			true,
		},
	}
	for _, testCase := range tests {
		currentCase := testCase
		t.Run(currentCase.name, func(t *testing.T) {
			path, err := NewPath(currentCase.source)
			assert.NoError(t, err)
			actual, err := path.Get(currentCase.data)
			if currentCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentCase.expected, actual)
			}
		})
	}
}
