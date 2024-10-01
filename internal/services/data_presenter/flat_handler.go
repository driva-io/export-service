package data_presenter

import "errors"

func handleFlat(source map[string]any, location any) (any, error) {

	var allResults []any
	arrayLocation, isArray := location.([]any)
	if !isArray {
		return nil, errors.New("$flat requires an array")
	}

	for _, value := range arrayLocation {
		result, err := Apply(source, value, make(map[string]any))
		if err != nil {
			return nil, err
		}
		if result != nil {
			allResults = append(allResults, result)
		}
	}

	if (allResults != nil) && len(allResults) > 0 {
		return flatMap(allResults), nil
	}

	return nil, nil
}
