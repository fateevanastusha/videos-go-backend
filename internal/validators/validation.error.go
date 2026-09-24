package validators

import "strings"

type ErrorMessage struct {
	Message string `json:"message"`
	Field   string `json:"field"`
}

type ValidationError struct {
	Errors []ErrorMessage
}

func (e ValidationError) Error() string {
	values := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		values[i] = err.Field + ": " + err.Message
	}

	return strings.Join(values, ", ")
}
