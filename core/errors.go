package core

import (
	"encoding/json"
	"net/http"
)

// TYPoc represent the errors used in the Your Minute api
type TYPoc struct {
	Code    int    `json:"code"`
	From    string `json:"error-from, omitempty"`
	Element string `json:"error-on, omitempty"`
	Message string `json:"error-message"`
	Method  string `json:"method"`
}

// NewModelError is used to manage error from Models definitions
func NewModelError(from string, element string, message string) *TYPoc {
	return &TYPoc{
		Code:    http.StatusUnprocessableEntity,
		From:    from,
		Element: element,
		Message: message,
	}
}

// NewDatastoreError is used to manage error from data store part
func NewDatastoreError(from string, element string, message string) *TYPoc {
	return &TYPoc{
		Code:    http.StatusInternalServerError,
		From:    from,
		Element: element,
		Message: message,
	}
}

// NewAPIError is used to manage error from api part
func NewAPIError(from string, element string, message string, code int, method string) *TYPoc {
	return &TYPoc{
		Code:    code,
		From:    from,
		Element: element,
		Message: message,
		Method:  method,
	}
}

// NewFatalError is used to manage Fatal error from dependencies. Those error should not be returned except
// when depedencie broke
func NewFatalError(from string, element string, message string, method string) *TYPoc {
	return &TYPoc{
		Code:    http.StatusInternalServerError,
		From:    from,
		Element: element,
		Message: "FATAL - " + message,
		Method:  method,
	}
}

// ToJSON function take an avatar and tranform it into Json object
func (ym *TYPoc) ToJSON() string {
	// Error cannot be reached in normal used case as the function can only be call on a valid avatar object who does not contain invalid types nor recursive reference
	b, _ := json.Marshal(ym)
	return string(b)
}
