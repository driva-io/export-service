package data_presenter

import (
	"errors"
	"fmt"
)

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

		resultStr := fmt.Sprintf("%v", result)

		if result != nil {
			fullResult = fullResult + key + ": " + resultStr + "\n"
		}
	}

	return fullResult, nil
}
