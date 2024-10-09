package data_presenter

import (
	"errors"
)

func handlePhone(source map[string]any, location any) (any, error) {
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
		return nil, errors.New("$phone requires a string")
	}

	return "+55 " + stringResult, nil
}
