package data_presenter

import (
	"errors"
	"reflect"
)

func handleFor(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	mapLocation, isMap := location.(map[string]any)
	if !isMap {
		return nil, errors.New("invalid data type for $for keyword. must be a map[string]any")
	}

	prop, exists := mapLocation["$prop"]
	if !exists {
		return nil, errors.New("$for missing $prop keyword")
	}

	format, exists := mapLocation["$format"]
	if !exists {
		return nil, errors.New("$for missing $format keyword")
	}

	var filter map[string]any
	if filterKey, exists := mapLocation["$filter"]; exists {
		filterKey, isMap := filterKey.(map[string]any)
		if !isMap {
			return nil, errors.New("invalid data type for $filter keyword. must be a map[string]any")
		}
		filter = filterKey
	}

	var index int
	indexKey, indexExists := mapLocation["$index"]
	if indexExists {
		intIndex, isInt := indexKey.(int)
		if !isInt {
			return nil, errors.New("invalid data type for $index keyword. must be an int")
		}
		index = intIndex
	}

	var limit int
	limitKey, limitExists := mapLocation["$limit"]
	if limitExists {
		if indexExists {
			return nil, errors.New("$limit and $index cant exist simultaneously in $for keyword")
		}
		intLimit, isInt := limitKey.(int)
		if !isInt {
			return nil, errors.New("invalid data type for $limit keyword. must be an int")
		}

		if intLimit == 0 {
			return nil, errors.New("$limit must be greater than 0")
		}
		limit = intLimit
	}

	var result any
	var err error
	stringProp, isString := prop.(string)
	if isString {
		result, err = getNestedValue(source, stringProp)
	} else {
		result, err = handler.HandleKeywords(source, prop)
	}

	if result == nil {
		return nil, nil
	}

	arrayResult, isArray := result.([]any)
	if !isArray {
		return nil, errors.New("$for requires an array $prop")
	}

	if err != nil {
		return nil, err
	}

	var forResult []any
	for resultKey, value := range arrayResult {

		mapValue, isMap := value.(map[string]any)
		if !isMap {
			return nil, errors.New("every value in $prop array must be a map[string]any")
		}

		if filter != nil {
			match := true
			for key, condition := range filter {
				if itemVal, exists := mapValue[key]; !exists || !reflect.DeepEqual(itemVal, condition) {
					match = false
					break
				}
			}

			if !match {
				continue
			}
		}

		mapValueCopy := make(map[string]any)
		for key, value := range mapValue {
			mapValueCopy[key] = value
		}
		mapValueCopy["super$"] = source
		res, err := Apply(mapValueCopy, format, make(map[string]any))
		if err != nil {
			return nil, err
		}

		if indexExists && (resultKey == index) {
			return res, nil
		}

		if limitExists && (resultKey == limit) {
			return forResult, nil
		}

		if !indexExists {
			forResult = append(forResult, res)
		}
	}

	if len(forResult) == 0 {
		return nil, nil
	}

	return forResult, nil
}
