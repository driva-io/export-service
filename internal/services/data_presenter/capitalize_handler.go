package data_presenter

import (
	"errors"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func handleCapitalize(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	if stringValue, isString := result.(string); isString {
		caser := cases.Title(language.Und)
		return caser.String(stringValue), nil
	}

	return nil, errors.New("invalid data type for $capitalize conversion")
}
