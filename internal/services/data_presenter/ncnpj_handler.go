package data_presenter

import "fmt"

func handleNcnpj(source map[string]any, location any) (any, error) {

	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	result, err = handleString(source, location)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return fmt.Sprintf("%014s", result), nil
}
