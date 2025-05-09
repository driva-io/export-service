package data_presenter

import (
	"errors"
	"fmt"
)

func handleCnpj(source map[string]any, location any) (any, error) {

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
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		stringResult = fmt.Sprintf("%d", v)
	case float64:
		if v == float64(int64(v)) {
			stringResult = fmt.Sprintf("%.0f", v)
		} else {
			stringResult = fmt.Sprintf("%f", v)
		}
	case float32:
		if v == float32(int64(v)) {
			stringResult = fmt.Sprintf("%.0f", v)
		} else {
			stringResult = fmt.Sprintf("%f", v)
		}
	default:
		return nil, errors.New("all $compositestring key's results must be stringifiable")
	}

	paddedCnpj := fmt.Sprintf("%014s", stringResult)

	return fmt.Sprintf("%s.%s.%s/%s-%s",
		paddedCnpj[0:2],
		paddedCnpj[2:5],
		paddedCnpj[5:8],
		paddedCnpj[8:12],
		paddedCnpj[12:14],
	), nil
}
