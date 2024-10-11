package ports

import "fmt"

type RFC7807Error struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e RFC7807Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

type InvalidQueryParamsError struct {
	RFC7807Error
	//Aditional fields
}

type InvalidBodyError struct {
	RFC7807Error
	//Aditional fields
}

func NewInvalidBodyError() InvalidBodyError {
	return InvalidBodyError{
		RFC7807Error: RFC7807Error{
			Type:   "InvalidBody",
			Title:  "Invalid Body Params",
			Detail: "Invalid Body Params.",
		},
	}
}

func NewInvalidQueryParamsError() InvalidQueryParamsError {
	return InvalidQueryParamsError{
		RFC7807Error: RFC7807Error{
			Type:   "InvalidQueryParams",
			Title:  "Invalid Query Params",
			Detail: "Invalid Query Params.",
		},
	}
}
