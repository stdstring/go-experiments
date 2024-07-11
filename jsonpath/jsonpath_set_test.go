package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPathSetDefault(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	value := map[string]interface{}{"c": "d"}
	path, err := NewPath("$")
	assert.NoError(t, err)
	result, err := path.Set(test, value)
	assert.NoError(t, err)
	assert.Equal(t, result, value)
}

func TestJSONPathSetSimple(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	path, err := NewPath("$.a")
	assert.NoError(t, err)
	result, err := path.Set(test, "s")
	assert.NoError(t, err)
	out, err := path.Get(result)
	assert.NoError(t, err)
	assert.Equal(t, "s", out)
}

func TestJSONPathSetDeep(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	outer := map[string]interface{}{"x": test}
	path, err := NewPath("$.x.a")
	assert.NoError(t, err)
	result, err := path.Set(outer, "s")
	assert.NoError(t, err)
	out, err := path.Get(result)
	assert.NoError(t, err)
	assert.Equal(t, "s", out)
}

func TestJSONPathSetCreate(t *testing.T) {
	test := map[string]interface{}{}
	path, err := NewPath("$.a")
	assert.NoError(t, err)
	result, err := path.Set(test, "s")
	assert.NoError(t, err)
	out, err := path.Get(result)
	assert.NoError(t, err)
	assert.Equal(t, "s", out)
}

func TestJSONPathSetOverwrite(t *testing.T) {
	test := map[string]interface{}{"a": "b"}
	path, err := NewPath("$.a.b")
	assert.NoError(t, err)
	result, err := path.Set(test, "s")
	assert.NoError(t, err)
	out, err := path.Get(result)
	assert.NoError(t, err)
	assert.Equal(t, "s", out)
}
