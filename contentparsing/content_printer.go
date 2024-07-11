package contentparsing

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const indentationDelta = 2

func prettyPrintContent(rawContent any, indentationLevel int) {
	indentation := strings.Repeat(" ", indentationLevel)
	switch content := rawContent.(type) {
	case *JsonData:
		fmt.Printf("%sJSON:\n", indentation)
		prettyPrintJsonData(content, indentationLevel+indentationDelta)
	case *XmlData:
		fmt.Printf("%sXML:\n", indentation)
		prettyPrintXmlData(content, indentationLevel+indentationDelta)
	case *Base64Data:
		fmt.Printf("%sBASE64:\n", indentation)
		prettyPrintContent(content.Value, indentationLevel+indentationDelta)
	default:
		fmt.Printf("%s\"%v\"\n", indentation, content)
	}
}

func prettyPrintXmlData(xmlData *XmlData, indentationLevel int) {
	prettyPrintXmlDirectives(xmlData.EntityDirectives, "ENTITIES", indentationLevel)
	prettyPrintXmlDirectives(xmlData.DoctypeDirectives, "DOCTYPES", indentationLevel)
	prettyPrintXmlDirectives(xmlData.OtherDirectives, "OTHER", indentationLevel)
	prettyPrintXmlElement(xmlData.Root, indentationLevel)
}

func prettyPrintXmlDirectives(directives []string, title string, indentationLevel int) {
	indentation := strings.Repeat(" ", indentationLevel)
	if len(directives) == 0 {
		fmt.Printf("%s%s: []\n", indentation, title)
		return
	}
	indentationForData := strings.Repeat(" ", indentationLevel+indentationDelta)
	fmt.Printf("%s%s:\n", indentation, title)
	for _, directive := range directives {
		fmt.Printf("%s%s\n", indentationForData, directive)
	}
}

func prettyPrintXmlElement(element *XmlElement, indentationLevel int) {
	indentation := strings.Repeat(" ", indentationLevel)
	fmt.Printf("%sName = \"%s\"\n", indentation, getXmlFullName(element.Name))
	prettyPrintAttributes(element.Attributes, indentation)
	if element.Value != nil {
		fmt.Printf("%sValue:\n", indentation)
		prettyPrintContent(element.Value, indentationLevel+indentationDelta)
	}
	for _, child := range element.Children {
		prettyPrintXmlElement(child, indentationLevel+indentationDelta)
	}
}

func prettyPrintAttributes(attributes []xml.Attr, indentation string) {
	if len(attributes) == 0 {
		return
	}
	fmt.Printf("%sAttributes: ", indentation)
	for index, attribute := range attributes {
		if index > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s = \"%s\"", getXmlFullName(attribute.Name), attribute.Value)
	}
	fmt.Println()
}

func prettyPrintJsonData(jsonData *JsonData, indentationLevel int) {
	prettyPrintJsonValue(jsonData.Value, indentationLevel)
}

func prettyPrintJsonValue(rawValue any, indentationLevel int) {
	indentation := strings.Repeat(" ", indentationLevel)
	switch value := rawValue.(type) {
	case JsonArrayValue:
		prettyPrintJsonArray(value, indentationLevel)
	case JsonObjectValue:
		prettyPrintJsonObject(value, indentationLevel)
	case nil:
		fmt.Printf("%snull\n", indentation)
	case JsonSimpleValue:
		prettyPrintContent(rawValue, indentationLevel)
	}
}

func prettyPrintJsonObject(jsonObject JsonObjectValue, indentationLevel int) {
	indentation := strings.Repeat(" ", indentationLevel)
	indentationForData := strings.Repeat(" ", indentationLevel+indentationDelta)
	if len(jsonObject) == 0 {
		fmt.Printf("%sOBJECT: {}\n", indentation)
		return
	}
	fmt.Printf("%sOBJECT: {\n", indentation)
	for fieldName, jsonValue := range jsonObject {
		fmt.Printf("%s%s:\n", indentationForData, fieldName)
		prettyPrintJsonValue(jsonValue, indentationLevel+2*indentationDelta)
	}
	fmt.Printf("%s}\n", indentation)
}

func prettyPrintJsonArray(jsonArray JsonArrayValue, indentationLevel int) {
	indentation := strings.Repeat(" ", indentationLevel)
	if len(jsonArray) == 0 {
		fmt.Printf("%sARRAY: []\n", indentation)
		return
	}
	fmt.Printf("%sARRAY: [\n", indentation)
	for _, jsonValue := range jsonArray {
		prettyPrintJsonValue(jsonValue, indentationLevel+indentationDelta)
	}
	fmt.Printf("%s]\n", indentation)
}
