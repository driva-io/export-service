package data_presenter

import "fmt"

func handleCnpj(source map[string]any, location any) (any, error) {

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

	paddedCnpj := fmt.Sprintf("%014s", result)

	return fmt.Sprintf("%s.%s.%s/%s-%s",
		paddedCnpj[0:2],
		paddedCnpj[2:5],
		paddedCnpj[5:8],
		paddedCnpj[8:12],
		paddedCnpj[12:14],
	), nil
}
