package data_presenter

import (
	"errors"
	"strings"
)

func handleUpper(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	if stringResult, isString := result.(string); isString {
		return strings.ToUpper(stringResult), nil
	}

	return nil, errors.New("$upper requires a string")
}
