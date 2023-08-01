package shared

import "fmt"

type InvalidValueError struct {
	ValueType string
	Value     string
}

func (e InvalidValueError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.ValueType, e.Value)
}
