package data_presenter

import (
	"errors"
	"fmt"
)

func handleNcnpj(source map[string]any, location any) (any, error) {

	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	var stringResult string
	switch v := result.(type) {
	case string:
		stringResult = v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		stringResult = fmt.Sprintf("%v", v)
	default:
		return nil, errors.New("all $compositestring key's results must be stringifiable")
	}

	return fmt.Sprintf("%014s", stringResult), nil
}
