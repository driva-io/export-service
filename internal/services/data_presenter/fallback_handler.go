package data_presenter

import "errors"

func handleFallback(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	arrayLocation, isArray := location.([]any)
	if !isArray {
		return nil, errors.New("$fallback requires an array of props")
	}

	for _, value := range arrayLocation {
		result, err := handler.HandleKeywords(source, value)
		if err != nil {
			return nil, err
		}

		if result != nil && result != "" {
			return result, nil
		}
	}

	return nil, nil
}
