package data_presenter

import (
	"errors"
	"strings"
)

func handleJoinBy(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	mapLocation, isMap := location.(map[string]any)
	if !isMap {
		return nil, errors.New("$joinby requires a map")
	}

	prop, exists := mapLocation["$prop"]
	if !exists {
		return nil, errors.New("$joinby requires $prop")
	}

	separator := ", "
	if customSeparator, exists := location.(map[string]any)["$separator"]; exists {
		if stringSeparator, isString := customSeparator.(string); isString {
			separator = stringSeparator
		} else {
			return nil, errors.New("invalid data type for $separator")
		}
	}

	var result any
	var err error
	if stringProp, isString := prop.(string); isString {
		result, err = getNestedValue(source, stringProp)
	} else {
		result, err = handler.HandleKeywords(source, mapLocation["$prop"])
	}
	if err != nil {
		return nil, err
	}

	stringValues := make([]string, 0, len(result.([]any)))
	for _, value := range result.([]any) {
		strValue, ok := value.(string)
		if !ok {
			return nil, errors.New("a []string must be provided to joinby")
		}
		stringValues = append(stringValues, strValue)
	}

	result = strings.Join(stringValues, separator)
	if result == "" {
		result = nil
	}

	return result, nil
}
