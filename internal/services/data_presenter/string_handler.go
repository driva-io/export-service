package data_presenter

import (
	"errors"
	"strconv"
)

func handleString(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	switch v := result.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return nil, errors.New("invalid data type for $string conversion")
	}
}
