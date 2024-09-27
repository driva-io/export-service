package services

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

	return source, nil
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
			target[k] = make([]any, len(v))
			for i, prop := range v {
				result, err := Apply(source, prop.(map[string]any), make(map[string]any))
				if err != nil {
					return nil, err
				}
				if result != nil {
					target[k].([]any)[i] = result
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
			if result != nil {
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
		result[key] = tabResult
	}
	return result, nil
}
