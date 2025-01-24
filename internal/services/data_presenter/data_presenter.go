package data_presenter

import (
	"errors"
	"export-service/internal/core/domain"
	"strings"
)

func flatMap(nestedArray []any) []any {
	var flatResult []any

	for _, value := range nestedArray {
		// Check if the value is also an array
		if innerArray, isInnerArray := value.([]any); isInnerArray {
			flatResult = append(flatResult, innerArray...)
		} else {
			flatResult = append(flatResult, value)
		}
	}

	return flatResult
}

func copyValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		copyMap := make(map[string]any)
		for key, val := range v {
			copyMap[key] = copyValue(val)
		}
		return copyMap
	case []any:
		copySlice := make([]any, len(v))
		for i, val := range v {
			copySlice[i] = copyValue(val)
		}
		return copySlice
	default:
		return v
	}
}

// source is an object or array of objects and location is a string or array of strings
func getNestedValue(source any, location any) (any, error) {
	var locs []string
	if stringLocation, isString := location.(string); isString {
		locs = strings.Split(stringLocation, ".")
	} else {
		locs = location.([]string)
	}

	if len(locs) > 0 && source != nil {
		nested := source.(map[string]any)[locs[0]]
		if nestedArray, isArray := nested.([]any); isArray {
			var result []any
			for _, value := range nestedArray {
				nestedValue, err := getNestedValue(value, locs[1:])
				if err != nil {
					return nil, err
				}
				if nestedValue != nil {
					result = append(result, nestedValue)
				}
			}

			flatResult := flatMap(result)
			return flatResult, nil
		} else {
			return getNestedValue(source.(map[string]any)[locs[0]], locs[1:])
		}
	}

	return copyValue(source), nil
}

func Apply(source map[string]any, spec any, target map[string]any) (any, error) {
	if _, isString := spec.(string); isString {
		return getNestedValue(source, spec)
	}

	var mapSpec map[string]any
	if mapStringSpec, isMap := spec.(map[string]any); isMap {
		mapSpec = mapStringSpec
	} else {
		return nil, errors.New("spec should be a string or map[string]any")
	}

	keys := make([]string, 0, len(mapSpec))
	for k := range mapSpec {
		keys = append(keys, k)
	}
	return buildStructure(keys, source, mapSpec, target)
}

func buildStructure(keys []string, source, spec map[string]any, target map[string]any) (any, error) {

	for _, k := range keys {
		if strings.HasPrefix(k, "$") {
			handler := NewKeywordHandler()
			return handler.HandleKeyword(source, spec)
		}

		//Calls recursion or get value depending on received mapping
		switch v := spec[k].(type) {
		case []any:
			target[k] = []any{}
			for _, prop := range v {
				result, err := Apply(source, prop, make(map[string]any))
				if err != nil {
					return nil, err
				}
				if result != nil {
					target[k] = append(target[k].([]any), result)
				}
			}
		case map[string]any:
			result, err := Apply(source, v, make(map[string]any))
			if err != nil {
				return nil, err
			}
			if result != nil {
				target[k] = result
			}
		case string:
			result, err := getNestedValue(source, v)
			if err != nil {
				return nil, err
			}
			if result != nil && result != "" {
				target[k] = result
			}
		}
	}

	return target, nil
}

func PresentSingle(data map[string]any, spec domain.PresentationSpecSpec) (map[string]any, error) {
	result := make(map[string]any)
	for key, mapping := range spec {
		tabResult, err := Apply(data, mapping, make(map[string]any))
		if err != nil {
			return nil, err
		}

		if tabResult == nil {
			continue
		}

		arrayResult, isArray := tabResult.([]any)
		if isArray && len(arrayResult) == 0 {
			continue
		}

		result[key] = tabResult

	}
	return result, nil
}

func PresentMultiple(data []map[string]any, spec domain.PresentationSpecSpec) (map[string]any, error) {
	result := make(map[string]any)
	for _, values := range data {
		for specKey, mapping := range spec {
			tabResult, err := Apply(values, mapping, make(map[string]any))
			if err != nil {
				return nil, err
			}
			result[specKey] = tabResult
		}
	}
	return result, nil
}
