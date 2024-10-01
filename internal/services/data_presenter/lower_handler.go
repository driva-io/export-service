package data_presenter

import (
	"errors"
	"strings"
)

func handleLower(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	if stringResult, isString := result.(string); isString {
		return strings.ToLower(stringResult), nil
	}

	return nil, errors.New("$lower requires a string")
}
