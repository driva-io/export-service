package repositories

import (
	"fmt"
)

type RFC7807Error struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e RFC7807Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

type CompanyNotFoundError struct {
	RFC7807Error
	//Aditional fields
}

type SolicitationNotFoundError struct {
	RFC7807Error
	//Aditional fields
}

type CompanyNotUniqueError struct {
	RFC7807Error
	//Aditional fields
}

type PresentationSpecNotFoundError struct {
	RFC7807Error
	//Aditional fields
}

type PresentationSpecNotUniqueError struct {
	RFC7807Error
	//Aditional fields
}

type InvalidJsonBodyError struct {
	RFC7807Error
	//Aditional fields
}

type MissingQueryParametersError struct {
	RFC7807Error
	//Aditional fields
}

type ContactsListNotFoundError struct {
	RFC7807Error
	//Aditional fields
}

type PersonaNotFoundError struct {
	RFC7807Error
}

func NewInternalServerError() RFC7807Error {
	return RFC7807Error{
		Type:   "InternalServerError",
		Title:  "Internal Server Error",
		Detail: "Internal Server Error.",
	}
}

func NewCompanyNotFoundError() CompanyNotFoundError {
	return CompanyNotFoundError{
		RFC7807Error: RFC7807Error{
			Type:   "CompanyNotFoundError",
			Title:  "Company Not Found",
			Detail: "The requested company could not be found.",
		},
	}
}

func NewSolicitationNotFoundError() SolicitationNotFoundError {
	return SolicitationNotFoundError{
		RFC7807Error: RFC7807Error{
			Type:   "SolicitationNotFoundError",
			Title:  "Solicitation Not Found",
			Detail: "The requested solicitation could not be found.",
		},
	}
}

func NewCompanyNotUniqueError() CompanyNotUniqueError {
	return CompanyNotUniqueError{
		RFC7807Error: RFC7807Error{
			Type:   "CompanyNotUniqueError",
			Title:  "Company Not Unique",
			Detail: "More than one company was found for given parameters.",
		},
	}
}

func NewPresentationSpecNotFoundError() PresentationSpecNotFoundError {
	return PresentationSpecNotFoundError{
		RFC7807Error: RFC7807Error{
			Type:   "PresentationSpecNotFoundError",
			Title:  "Presentation Spec Not Found",
			Detail: "The requested presentation specification could not be found.",
		},
	}
}

func NewPresentationSpecNotUniqueError() PresentationSpecNotUniqueError {
	return PresentationSpecNotUniqueError{
		RFC7807Error: RFC7807Error{
			Type:   "PresentationSpecNotUniqueError",
			Title:  "Presentation Spec Not Unique",
			Detail: "More than one presentation specification was found for given parameters.",
		},
	}
}

func NewInvalidJsonBodyError() InvalidJsonBodyError {
	return InvalidJsonBodyError{
		RFC7807Error: RFC7807Error{
			Type:   "InvalidJsonBodyError",
			Title:  "Invalid Request Body",
			Detail: "Invalid json in request body.",
		},
	}
}

func NewMissingQueryParametersError() MissingQueryParametersError {
	return MissingQueryParametersError{
		RFC7807Error: RFC7807Error{
			Type:   "MissingQueryParametersError",
			Title:  "Missing Query Parameters",
			Detail: "Missing query parameters in call url.",
		},
	}
}

func NewContactsListNotFoundError() ContactsListNotFoundError {
	return ContactsListNotFoundError{
		RFC7807Error: RFC7807Error{
			Type:   "ContactsListNotFoundError",
			Title:  "Contacts List Not Found",
			Detail: "The requested contacts list configuration could not be found.",
		},
	}
}

func NewPersonaNotFoundError() PersonaNotFoundError {
	return PersonaNotFoundError{
		RFC7807Error: RFC7807Error{
			Type:   "PersonaNotFoundError",
			Title:  "Persona Not Found",
			Detail: "The requested persona could not be found.",
		},
	}
}
