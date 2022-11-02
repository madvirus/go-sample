package common

import "strings"

func CreateValidationError(errors []ErrorField) *ValidationError {
	return &ValidationError{Errors: errors}
}

type ValidationError struct {
	Errors []ErrorField
}

func (v *ValidationError) Error() string {
	var sb strings.Builder
	sb.WriteString("ValidationError: ")
	for _, e := range v.Errors {
		sb.WriteString(e.Name)
		sb.WriteString("=")
		sb.WriteString(e.Message)
		sb.WriteString(", ")
	}
	return sb.String()
}

type ErrorField struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func CreateNoDataFoundError(message string) *NoDataFoundError {
	return &NoDataFoundError{Message: message}
}

type NoDataFoundError struct {
	Message string `json:"message"`
}

func (n *NoDataFoundError) Error() string {
	return n.Message
}
