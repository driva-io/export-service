package data_presenter

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

func handleTemplate(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	mapLocation, isMap := location.(map[string]any)
	if !isMap {
		return nil, errors.New("invalid data type for $template keyword. must be a map[string]any")
	}

	formatValue, exists := mapLocation["$format"]
	if !exists {
		return nil, errors.New("$template requires $format keyword")
	}
	format, isString := formatValue.(string)
	if !isString {
		return nil, errors.New("$format must be a string")
	}

	variablesValue, exists := mapLocation["$variables"]
	if !exists {
		return nil, errors.New("$template requires $variables keyword")
	}
	variables, isMap := variablesValue.(map[string]any)
	if !isMap || len(variables) == 0 {
		return nil, errors.New("$variables must be a map[string]any")
	}

	re := regexp.MustCompile(`\{(\w+)\}`)

	forProp, exists := mapLocation["$for"]
	if exists {
		var filter map[string]any
		if filterKey, exists := location.(map[string]any)["$filter"]; exists {
			filterKey, isMap := filterKey.(map[string]any)
			if !isMap {
				return nil, errors.New("invalid data type for $filter keyword. must be a map[string]any")
			}
			filter = filterKey
		}

		prop, isString := forProp.(string)
		if !isString {
			return nil, errors.New("$for prop must be a string")
		}

		results, err := getNestedValue(source, prop)
		if err != nil {
			return nil, err
		}

		if results == nil {
			return nil, nil
		}

		arrayResult, isArray := results.([]any)
		if !isArray {
			return nil, errors.New("$for prop must result in an array of maps")
		}

		var allResults []string
		for _, value := range arrayResult {
			variableValues := make(map[string]string, len(variables))

			valueCopy := make(map[string]any)
			for key, val := range value.(map[string]any) {
				valueCopy[key] = val
			}
			valueCopy["super$"] = source

			for key, variable := range variables {
				varResult, err := handler.HandleKeywords(valueCopy, variable)
				if err != nil {
					return nil, err
				}

				if varResult == nil {
					continue
				}

				result, isString := varResult.(string)
				if !isString {
					return nil, errors.New("every $template variable must result in a string")
				}
				variableValues[key] = result
			}

			if filter != nil {
				match := true
				for key, condition := range filter {
					if itemVal, exists := value.(map[string]any)[key]; !exists || !reflect.DeepEqual(itemVal, condition) {
						match = false
						break
					}
				}

				if !match {
					continue
				}
			}

			result := re.ReplaceAllStringFunc(format, func(match string) string {
				varName := re.FindStringSubmatch(match)[1]
				if value, exists := variableValues[varName]; exists {
					return value
				}
				return match
			})

			allResults = append(allResults, result)
		}

		return strings.Join(allResults, ""), nil
	}

	variableValues := make(map[string]string, len(variables))
	for key, variable := range variables {
		varResult, err := handler.HandleKeywords(source, variable)
		if err != nil {
			return nil, err
		}
		if varResult == nil {
			continue
		}

		result, isString := varResult.(string)
		if !isString {
			return nil, errors.New("every $template variable must result in a string")
		}
		variableValues[key] = result
	}
	result := re.ReplaceAllStringFunc(format, func(match string) string {
		varName := re.FindStringSubmatch(match)[1]
		if value, exists := variableValues[varName]; exists {
			return value
		}
		return ""
	})

	return result, nil
}
