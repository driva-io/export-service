package data_presenter

import "errors"

func handleCompositeString(source map[string]any, location any) (any, error) {
	mapLocation, isMap := location.(map[string]any)
	if !isMap {
		return nil, errors.New("$compositestring requires a map")
	}

	handler := NewKeywordHandler()
	var fullResult string
	for key, value := range mapLocation {
		result, err := handler.HandleKeywords(source, value)
		if err != nil {
			return nil, err
		}

		stringResult, isString := result.(string)
		if !isString && result != nil {
			return nil, errors.New("all $compositestring key's results must be a string")
		}

		if result != nil {
			fullResult = fullResult + key + ": " + stringResult + "\n"
		}
	}

	return fullResult, nil
}
