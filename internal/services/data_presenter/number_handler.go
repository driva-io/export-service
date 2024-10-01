package data_presenter

import (
	"errors"
	"strconv"
)

func handleNumber(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	switch v := result.(type) {
	case int, float32, float64:
		return v, nil
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			return num, nil
		} else if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, nil
		}
		return nil, errors.New("$number requires a string that can be converted to a number")
	default:
		return nil, errors.New("$number requires a string that can be converted to a number")
	}

}
