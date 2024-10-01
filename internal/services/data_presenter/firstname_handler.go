package data_presenter

import (
	"errors"
	"strings"
)

func handleFirstName(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	stringResult, isString := result.(string)
	if !isString {
		return nil, errors.New("$firstname requires a string")
	}

	return strings.Split(stringResult, " ")[0], nil
}
