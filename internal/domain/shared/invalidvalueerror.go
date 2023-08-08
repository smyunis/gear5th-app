package shared

import "fmt"

type ErrInvalidValue struct {
	ValueType  string
	Value      string
	InnerError error
}

func NewInvalidValueError(valueType, value string) ErrInvalidValue {
	return ErrInvalidValue{
		ValueType: valueType,
		Value:     value,
	}
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("%s is not a valid %s", e.Value, e.ValueType)
}
