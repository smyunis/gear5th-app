package shared

import "fmt"

type InvalidValueError struct {
	ValueType  string
	Value      string
	InnerError error
}

func NewInvalidValueError(valueType, value string) InvalidValueError {
	return InvalidValueError{
		ValueType: valueType,
		Value:     value,
	}
}

func (e InvalidValueError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.ValueType, e.Value)
}
