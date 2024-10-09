package data_presenter

import (
	"errors"
	"reflect"
)

// every $case can be a map of conditions, using OR logic. if more than one $case has a true condition, first $use in order is returned
func handleSwitch(source map[string]any, location any) (any, error) {

	mapLocation, isMap := location.(map[string]any)
	if !isMap {
		return nil, errors.New("$switch requires a map")
	}

	cases, exists := mapLocation["$cases"]
	if !exists {
		return nil, errors.New("$switch requires $cases")
	}

	casesArray, isArray := cases.([]any)
	if !isArray {
		return nil, errors.New("$cases must be an array of maps")
	}

	for _, value := range casesArray {
		mapValue, isMap := value.(map[string]any)
		if !isMap {
			return nil, errors.New("$cases must be an array of maps")
		}

		caseKeyword, exists := mapValue["$case"]
		if !exists {
			return nil, errors.New("each $case in $cases must have a $case key")
		}

		_, isMap = caseKeyword.(map[string]any)
		if !isMap {
			return nil, errors.New("every $case in $cases must be a map")
		}

		_, exists = mapValue["$use"]
		if !exists {
			return nil, errors.New("each $case in $cases must have a $use key")
		}
	}

	for _, value := range casesArray {

		mapValue := value.(map[string]any)
		caseKeyword := mapValue["$case"]
		mapCase := caseKeyword.(map[string]any)
		useValue := mapValue["$use"]

		match := false
		handler := NewKeywordHandler()
		for key, condition := range mapCase {
			result, err := handler.HandleKeywords(source, key)
			if err != nil {
				return nil, err
			}
			if reflect.DeepEqual(result, condition) {
				match = true
				break
			}
		}

		if !match {
			continue
		}

		useResult, err := handler.HandleKeywords(source, useValue)
		if err != nil {
			return nil, err
		}

		return useResult, nil
	}

	return nil, nil
}
