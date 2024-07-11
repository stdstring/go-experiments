// Package jsonpath: simple implementation of JSON Path for state machine
package jsonpath

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
The data must be from JSON Unmarshal, that way we can guarantee the types:

bool, for JSON booleans
float64, for JSON numbers
string, for JSON strings
[]interface{}, for JSON arrays
map[string]interface{}, for JSON objects
nil for JSON null
*/

var errNotFound = errors.New("not found")

type Path struct {
	path []string
}

func NewPath(pathString string) (*Path, error) {
	pathArray, err := parsePathString(pathString)
	path := Path{path: pathArray}
	return &path, err
}

func (path *Path) UnmarshalJSON(b []byte) error {
	var pathString string
	if err := json.Unmarshal(b, &pathString); err != nil {
		return err
	}
	pathArray, err := parsePathString(pathString)
	if err != nil {
		return err
	}
	path.path = pathArray
	return nil
}

func (path *Path) MarshalJSON() ([]byte, error) {
	if len(path.path) == 0 {
		return json.Marshal("$")
	}
	return json.Marshal(path.String())
}

func (path *Path) String() string {
	return fmt.Sprintf("$.%v", strings.Join(path.path, "."))
}

func (path *Path) GetTime(input interface{}) (*time.Time, error) {
	outputValue, err := path.Get(input)
	if err != nil {
		return nil, fmt.Errorf("GetTime Error %w", err)
	}
	var output time.Time
	switch outputValue := outputValue.(type) {
	case string:
		output, err = time.Parse(time.RFC3339, outputValue)
		if err != nil {
			return nil, fmt.Errorf("GetTime Error: time error %w", err)
		}
	default:
		return nil, fmt.Errorf("GetTime Error: time must be string")
	}
	return &output, nil
}

func (path *Path) GetBool(input interface{}) (*bool, error) {
	outputValue, err := path.Get(input)
	if err != nil {
		return nil, fmt.Errorf("GetBool Error %w", err)
	}
	var output bool
	switch outputValue := outputValue.(type) {
	case bool:
		output = outputValue
	default:
		return nil, fmt.Errorf("GetBool Error: must return bool")
	}
	return &output, nil
}

func (path *Path) GetNumber(input interface{}) (*float64, error) {
	outputValue, err := path.Get(input)
	if err != nil {
		return nil, fmt.Errorf("GetFloat Error %w", err)
	}
	var output float64
	switch outputValue := outputValue.(type) {
	case float64:
		output = outputValue
	case int:
		output = float64(outputValue)
	case []interface{}:
		output, ok := outputValue[0].(float64)
		if !ok {
			return &output, errors.New("get number: wrong type")
		}
	default:
		return nil, fmt.Errorf("GetFloat Error: must return float")
	}
	return &output, nil
}

func (path *Path) GetString(input interface{}) (*string, error) {
	outputValue, err := path.Get(input)
	if err != nil {
		return nil, fmt.Errorf("GetString Error %w", err)
	}
	var output string
	switch outputValue := outputValue.(type) {
	case string:
		output = outputValue
	default:
		return nil, fmt.Errorf("GetString Error: must return string")
	}
	return &output, nil
}

func (path *Path) GetMap(input interface{}) (map[string]interface{}, error) {
	outputValue, err := path.Get(input)
	if err != nil {
		return nil, fmt.Errorf("GetMap Error %w", err)
	}
	var output map[string]interface{}
	switch outputValue := outputValue.(type) {
	case map[string]interface{}:
		output = outputValue
	default:
		return nil, fmt.Errorf("GetMap Error: must return map")
	}
	return output, nil
}

func (path *Path) Get(input interface{}) (interface{}, error) {
	if path == nil {
		return input, nil // Default is $
	}
	return recursiveGet(input, path.path)
}

func (path *Path) GetSlice(input interface{}) ([]interface{}, error) {
	outputValue, err := path.Get(input)
	if err != nil {
		return nil, fmt.Errorf("GetSlice Error %w", err)
	}
	var output []interface{}
	switch outputValue := outputValue.(type) {
	case []interface{}:
		output = outputValue
	default:
		return nil, fmt.Errorf("GetSlice Error: must be an array")
	}
	return output, nil
}

func (path *Path) Set(input interface{}, value interface{}) (interface{}, error) {
	var setPath []string
	if path == nil {
		// default "$"
		setPath = []string{}
	} else {
		setPath = path.path
	}
	var output interface{}
	if len(setPath) == 0 {
		// The output is the value
		switch value := value.(type) {
		case map[string]interface{}:
			output = value
			return output, nil
		case []interface{}:
			output = value
			return output, nil
		case float64:
			output = value
			return output, nil
		case string:
			output = value
			return output, nil
		case bool:
			output = value
			return output, nil
		default:
			return nil, fmt.Errorf("cannot Set value %q type %q in root JSON path $", value, reflect.TypeOf(value))
		}
	}
	return recursiveSet(input, value, setPath), nil
}

func recursiveSet(data interface{}, value interface{}, path []string) map[string]interface{} {
	var dataMap map[string]interface{}
	switch data := data.(type) {
	case map[string]interface{}:
		dataMap = data
	default:
		// Overwrite current data with new map
		// this will work for nil as well
		dataMap = make(map[string]interface{})
	}
	if len(path) == 1 {
		dataMap[path[0]] = value
	} else {
		dataMap[path[0]] = recursiveSet(dataMap[path[0]], value, path[1:])
	}
	return dataMap
}

func recursiveGet(data interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return data, nil
	}
	if data == nil {
		return nil, errNotFound
	}
	switch data := data.(type) {
	case []interface{}:
		dest, isSingle, err := filterInterfaces(data, &path[0])
		if err != nil {
			return nil, err
		}
		if isSingle {
			value := dest[0]
			if len(path) == 1 {
				return value, nil
			}
			return recursiveGet(value, path[1:])
		}
		if len(path) == 1 {
			return dest, nil
		}
		combinedResult := make([]interface{}, 0)
		for _, fd := range dest {
			result, err := recursiveGet(fd, path[1:])
			if err != nil {
				return nil, err
			}
			combinedResult = appendResult(combinedResult, result)
		}
		return combinedResult, nil
	case map[string]interface{}:
		value, ok := data[path[0]]
		if !ok {
			return data, errNotFound
		}
		return recursiveGet(value, path[1:])
	case string:
		jsonStr := data
		var dataMap interface{}
		err := json.Unmarshal([]byte(jsonStr), &dataMap)
		if err != nil {
			return nil, errNotFound
		}
		if dataMap == nil {
			return nil, errNotFound
		}
		mapStruct, ok := dataMap.(map[string]interface{})
		if !ok {
			return mapStruct, errNotFound
		}
		value, ok := mapStruct[path[0]]
		if !ok {
			return mapStruct, errNotFound
		}
		return recursiveGet(value, path[1:])
	default:
		return data, errNotFound
	}
}

func extractSliceIndices(source string, dataSize int) (int, int, error) {
	sliceParts := strings.Split(source, ":")
	if parts := 2; len(sliceParts) != parts {
		return 0, 0, errors.New("bad slice definition")
	}
	sliceFrom := strings.Trim(sliceParts[0], " ")
	sliceTo := strings.Trim(sliceParts[1], " ")
	if (len(sliceFrom) == 0) && (len(sliceTo) == 0) {
		return 0, dataSize, nil
	}
	var indexFrom, indexTo int
	var err error
	if len(sliceFrom) > 0 {
		indexFrom, err = strconv.Atoi(strings.Trim(sliceFrom, " "))
		if err != nil {
			return 0, 0, err
		}
	}
	if len(sliceTo) > 0 {
		indexTo, err = strconv.Atoi(strings.Trim(sliceTo, " "))
		if err != nil {
			return 0, 0, err
		}
	}
	if len(sliceFrom) == 0 {
		if indexTo >= 0 {
			indexFrom = 0
		} else {
			indexFrom = -dataSize
		}
	}
	if len(sliceTo) == 0 {
		if indexFrom >= 0 {
			indexTo = dataSize
		} else {
			indexTo = 0
		}
	}
	if indexFrom > 0 {
		indexFrom = min(indexFrom, dataSize)
	} else {
		indexFrom = max(indexFrom, -dataSize)
	}
	if indexTo > 0 {
		indexTo = min(indexTo, dataSize)
	} else {
		indexTo = max(indexTo, -dataSize)
	}
	return indexFrom, indexTo, nil
}

func filterInterfaces(source []interface{}, filter *string) ([]interface{}, bool, error) {
	trimmedFilter := strings.Trim(*filter, "[]")
	filterParts := strings.Split(trimmedFilter, ",")
	if len(filterParts) == 1 {
		part := filterParts[0]
		if !strings.Contains(part, ":") {
			index, err := strconv.Atoi(strings.Trim(part, " "))
			if err != nil {
				return nil, false, err
			}
			dest := getValueByIndex(source, index)
			if len(dest) == 0 {
				return nil, false, errors.New("bad index")
			}
			return dest, true, nil
		}
		dest := make([]interface{}, 0)
		indexFrom, indexTo, err := extractSliceIndices(part, len(source))
		if err != nil {
			return nil, false, err
		}
		for index := indexFrom; index < indexTo; index++ {
			dest = append(dest, getValueByIndex(source, index)...)
		}
		return dest, false, nil
	}
	dest := make([]interface{}, 0)
	for _, part := range filterParts {
		index, err := strconv.Atoi(strings.Trim(part, " "))
		if err != nil {
			return nil, false, err
		}
		dest = append(dest, getValueByIndex(source, index)...)
	}
	return dest, false, nil
}

func getValueByIndex[T any](source []T, index int) []T {
	if (index >= len(source)) || (index < -len(source)) {
		return []T{}
	}
	if index < 0 {
		return []T{source[len(source)+index]}
	}
	return []T{source[index]}
}

func appendResult(dest []interface{}, result interface{}) []interface{} {
	switch result := result.(type) {
	case []interface{}:
		return append(dest, result...)
	default:
		return append(dest, result)
	}
}
