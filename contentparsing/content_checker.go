package contentparsing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkContent(t *testing.T, expected any, actual any) {
	t.Helper()
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	assert.NotNil(t, actual)
	switch expectedValue := expected.(type) {
	case *XmlData:
		actualValue, typeResult := actual.(*XmlData)
		assert.True(t, typeResult)
		checkXmlData(t, expectedValue, actualValue)
	case *JsonData:
		actualValue, typeResult := actual.(*JsonData)
		assert.True(t, typeResult)
		checkJsonData(t, expectedValue, actualValue)
	default:
		assert.Equal(t, expected, actual)
	}
}

func checkXmlData(t *testing.T, expected *XmlData, actual *XmlData) {
	t.Helper()
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	assert.NotNil(t, actual)
	assert.Equal(t, expected.DoctypeDirectives, actual.DoctypeDirectives)
	assert.Equal(t, expected.EntityDirectives, actual.EntityDirectives)
	assert.Equal(t, expected.OtherDirectives, actual.OtherDirectives)
	checkXmlElement(t, expected.Root, actual.Root)
}

func checkXmlElement(t *testing.T, expected *XmlElement, actual *XmlElement) {
	t.Helper()
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Attributes, actual.Attributes)
	assert.Equal(t, len(expected.Children), len(actual.Children))
	checkContent(t, expected.Value, actual.Value)
	for index := range expected.Children {
		checkXmlElement(t, expected.Children[index], actual.Children[index])
	}
}

func checkJsonData(t *testing.T, expected *JsonData, actual *JsonData) {
	t.Helper()
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	assert.NotNil(t, actual)
	checkJsonValue(t, expected.Value, actual.Value)
}

func checkJsonValue(t *testing.T, expected any, actual any) {
	t.Helper()
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	assert.NotNil(t, actual)
	switch expectedValue := expected.(type) {
	case JsonSimpleValue:
		actualValue, typeResult := actual.(JsonSimpleValue)
		assert.True(t, typeResult)
		checkContent(t, expectedValue, actualValue)
	case JsonArrayValue:
		actualValue, typeResult := actual.(JsonArrayValue)
		assert.True(t, typeResult)
		assert.Equal(t, len(expectedValue), len(actualValue))
		for index := range expectedValue {
			checkJsonValue(t, expectedValue[index], actualValue[index])
		}
	case JsonObjectValue:
		actualValue, typeResult := actual.(JsonObjectValue)
		assert.True(t, typeResult)
		assert.Equal(t, len(expectedValue), len(actualValue))
		for fieldName, expectedFieldValue := range expectedValue {
			actualFieldValue, valueResult := actualValue[fieldName]
			assert.True(t, valueResult)
			checkJsonValue(t, expectedFieldValue, actualFieldValue)
		}
	}
}
