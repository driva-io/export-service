package data_presenter

import (
	"errors"
	"time"
)

func handleDate(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	var dateInput string
	if dateStr, isString := result.(string); isString {
		dateInput = dateStr
	} else {
		return nil, errors.New("$date location must be a date formatted YYYY-MM-DD")
	}

	inputLayout := "2006-01-02"
	date, err := time.Parse(inputLayout, dateInput)
	if err != nil {
		date, err = time.Parse("2006-01-02T15:04:05Z", dateInput)
		if err != nil {
			return nil, errors.New("$date location must be a date formatted YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ")
		}
	}

	outputLayout := "02-01-2006"

	formattedDate := date.Format(outputLayout)

	return formattedDate, nil
}
