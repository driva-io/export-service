package services

import "export-service/internal/core/domain"

func PresentSingle(data map[string]any, spec domain.PresentationSpecSpec) (map[string]any, error) {
	presentedData := make(map[string]any)

	for key, value := range spec {
		d, err := apply(data, value)
		if err != nil {
			return nil, err
		}

		presentedData[key] = d
	}
	return presentedData, nil
}

func apply(data map[string]any, specValue map[string]any) (map[string]any, error) {
	final := make(map[string]any, len(specValue))

	for name, val := range specValue {
		final[name] = data[val.(string)]
	}
	return final, nil
}
