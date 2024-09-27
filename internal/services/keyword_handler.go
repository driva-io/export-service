package services

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type KeywordHandlerFunc func(source map[string]any, value any) (any, error)

func handleUpper(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	if stringResult, isString := result.(string); isString {
		return strings.ToUpper(stringResult), nil
	}

	return nil, errors.New("$upper requires a string")
}

func handleLower(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()
	result, err := handler.HandleKeywords(source, location)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	if stringResult, isString := result.(string); isString {
		return strings.ToLower(stringResult), nil
	}

	return nil, errors.New("$lower requires a string")
}

func handleTemplate(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	formatValue, exists := location.(map[string]any)["$format"]
	if !exists {
		return nil, errors.New("$template requires $format keyword")
	}
	format, isString := formatValue.(string)
	if !isString {
		return nil, errors.New("$format must be a string")
	}

	variablesValue, exists := location.(map[string]any)["$variables"]
	if !exists {
		return nil, errors.New("$template requires $variables keyword")
	}
	variables, isMap := variablesValue.(map[string]any)
	if !isMap || len(variables) == 0 {
		return nil, errors.New("$variables must be a map[string]any")
	}

	re := regexp.MustCompile(`\{(\w+)\}`)

	forProp, exists := location.(map[string]any)["$for"]
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
			for key, variable := range variables {
				value.(map[string]any)["super$"] = source
				varResult, err := handler.HandleKeywords(value.(map[string]any), variable)
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
		return match
	})

	return result, nil
}

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
		return nil, errors.New("invalid data type for $number conversion")
	default:
		return nil, errors.New("invalid data type for $number conversion")
	}

}

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

func handleFallback(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	if arrayValue, isArray := location.([]interface{}); isArray {
		for _, value := range arrayValue {
			if result, err := handler.HandleKeywords(source, value); err == nil && result != nil {
				return result, nil
			}
		}
		return nil, nil
	}

	return nil, errors.New("invalid data type for $fallback")
}

func handleNcnpj(source map[string]any, location any) (any, error) {

	result, err := handleString(source, location)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return fmt.Sprintf("%014s", result), nil
}

func handleCnpj(source map[string]any, location any) (any, error) {

	result, err := handleString(source, location)

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

func handleFor(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	if _, exists := location.(map[string]any)["$prop"]; !exists {
		return nil, errors.New("$for missing $prop keyword")
	}

	if _, exists := location.(map[string]any)["$format"]; !exists {
		return nil, errors.New("$for missing $format keyword")
	}

	var filter map[string]any
	if filterKey, exists := location.(map[string]any)["$filter"]; exists {
		filterKey, isMap := filterKey.(map[string]any)
		if !isMap {
			return nil, errors.New("invalid data type for $filter keyword. must be a map[string]any")
		}
		filter = filterKey
	}

	var result any
	var err error
	if prop, isString := location.(map[string]any)["$prop"].(string); isString {
		result, err = getNestedValue(source, prop)
	} else {
		result, err = handler.HandleKeywords(source, prop)
	}

	if err != nil {
		return nil, err
	}

	var forResult []any
	if arrayResult, isArray := result.([]any); isArray {
		for _, value := range arrayResult {

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

			value.(map[string]any)["super$"] = source
			res, err := Apply(value.(map[string]any), location.(map[string]any)["$format"], make(map[string]any))
			if err != nil {
				return nil, err
			}
			forResult = append(forResult, res)
		}
	} else {
		res, err := Apply(result.(map[string]any), location.(map[string]any)["$format"].(map[string]any), make(map[string]any))
		if err != nil {
			return nil, err
		}
		forResult = append(forResult, res)
	}

	return forResult, nil
}

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
		return nil, errors.New("$date location must be a date formatted YYYY-MM-DD")
	}

	outputLayout := "02-01-2006"

	formattedDate := date.Format(outputLayout)

	return formattedDate, nil
}

func handleJoinBy(source map[string]any, location any) (any, error) {
	handler := NewKeywordHandler()

	if prop, exists := location.(map[string]any)["$prop"]; exists {
		separator := ", "
		if customSeparator, exists := location.(map[string]any)["$separator"]; exists {
			if stringSeparator, isString := customSeparator.(string); isString {
				separator = stringSeparator
			} else {
				return nil, errors.New("invalid data type for $separator")
			}
		}

		var result any
		var err error
		if jbProp, isString := prop.(string); isString {
			result, err = getNestedValue(source, jbProp)
		} else {
			result, err = handler.HandleKeywords(source, location.(map[string]any)["$prop"])
		}
		if err != nil {
			return nil, err
		}

		stringValues := make([]string, 0, len(result.([]any)))
		for _, value := range result.([]any) {
			strValue, ok := value.(string)
			if !ok {
				return nil, errors.New("a []string must be provided to joinby")
			}
			stringValues = append(stringValues, strValue)
		}

		return strings.Join(stringValues, separator), nil
	}

	return nil, errors.New("invalid data type for $joinby")
}

func handleLiteral(source map[string]any, location any) (any, error) {
	return location, nil
}

func handleStringify(source map[string]any, location any) (any, error) {
	stringify := map[string]any{
		"$prop": location,
	}

	return handleJoinBy(source, stringify)
}

func handleFlat(source map[string]any, location any) (any, error) {

	var allResults []any
	arrayLocation, isArray := location.([]any)
	if !isArray {
		return nil, errors.New("$flat requires an array")
	}

	for _, value := range arrayLocation {
		result, err := Apply(source, value, make(map[string]any))
		if err != nil {
			return nil, err
		}
		if result != nil {
			allResults = append(allResults, result)
		}
	}

	if (allResults != nil) && len(allResults) > 0 {
		return flatMap(allResults), nil
	}

	return nil, nil
}

func handleCompositeString(source map[string]any, location any) (any, error) {
	mapLocation, isMap := location.(map[string]any)
	if !isMap {
		return nil, errors.New("$compositestring requires a map")
	}

	handler := NewKeywordHandler()
	var fullResult string
	for key, value := range mapLocation {
		result, err := handler.HandleKeywords(source, value)
		if err != nil {
			return nil, err
		}

		stringResult, isString := result.(string)
		if !isString {
			return nil, errors.New("all $compositestring key's results must be a string")
		}

		if result != nil {
			fullResult = fullResult + key + ": " + stringResult + "\n"
		}
	}

	return fullResult, nil
}

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

func handleSwitch(source map[string]any, location any) (any, error) {

	cases, exists := location.(map[string]any)["$cases"]
	if !exists {
		return nil, errors.New("$switch requires $cases")
	}

	casesArray, isArray := cases.([]any)
	if !isArray {
		return nil, errors.New("$cases must be an array of maps")
	}

	for _, value := range casesArray {
		caseValue, exists := value.(map[string]any)["$case"]
		if !exists {
			return nil, errors.New("each $case in $cases must have a $case key")
		}

		useValue, exists := value.(map[string]any)["$use"]
		if !exists {
			return nil, errors.New("each $case in $cases must have a $use key")
		}

		mapCase, isMap := caseValue.(map[string]any)
		if !isMap {
			return nil, errors.New("every $case in $cases must be a map")
		}

		match := false
		handler := NewKeywordHandler()
		for key, condition := range mapCase {
			result, err := handler.HandleKeywords(source, key)
			if err != nil {
				return nil, err
			}
			if reflect.DeepEqual(result, condition) {
				match = true
				break
			}
		}

		if !match {
			continue
		}

		useResult, err := handler.HandleKeywords(source, useValue)
		if err != nil {
			return nil, err
		}

		return useResult, nil
	}

	return nil, nil
}

type KeywordHandler struct {
	handlers map[string]KeywordHandlerFunc
}

func NewKeywordHandler() *KeywordHandler {
	return &KeywordHandler{
		handlers: map[string]KeywordHandlerFunc{
			"$upper":           handleUpper,
			"$lower":           handleLower,
			"$string":          handleString,
			"$number":          handleNumber,
			"$capitalize":      handleCapitalize,
			"$literal":         handleLiteral,
			"$fallback":        handleFallback,
			"$joinby":          handleJoinBy,
			"$for":             handleFor,
			"$ncnpj":           handleNcnpj,
			"$cnpj":            handleCnpj,
			"$date":            handleDate,
			"$stringify":       handleStringify,
			"$template":        handleTemplate,
			"$flat":            handleFlat,
			"$compositestring": handleCompositeString,
			"$phone":           handlePhone,
			"$switch":          handleSwitch,
		},
	}
}

func (h *KeywordHandler) HandleKeywords(source map[string]any, spec any) (any, error) {
	if stringSpec, isString := spec.(string); isString {
		return getNestedValue(source, stringSpec)
	}

	var nestedKeyword string
	var nestedSpec any
	if mapSpec, isMap := spec.(map[string]any); isMap {
		for key, value := range mapSpec {
			nestedKeyword = key
			nestedSpec = value
			break
		}
	}

	if handler, exists := h.handlers[nestedKeyword]; exists {
		return handler(source, nestedSpec)
	}

	return nil, errors.New("no handler found for the keyword " + nestedKeyword)
}

func (h *KeywordHandler) HandleKeyword(source map[string]any, spec map[string]any) (any, error) {
	for key, value := range spec {
		if handler, exists := h.handlers[key]; exists {
			return handler(source, value)
		}

		return nil, errors.New("no handler found for the keyword " + key)
	}

	return nil, nil
}
