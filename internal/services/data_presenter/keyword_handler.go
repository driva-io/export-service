package data_presenter

import (
	"errors"
)

type KeywordHandler struct {
	handlers map[string]KeywordHandlerFunc
}

type KeywordHandlerFunc func(source map[string]any, value any) (any, error)

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
			"$firstname":       handleFirstName,
			"$lastname":        handleLastName,
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
