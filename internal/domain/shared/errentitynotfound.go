package shared

import "fmt"

type ErrEntityNotFound struct {
	ID    string
	Scope string
}

func NewEntityNotFoundError(id, scope string) ErrEntityNotFound {
	return ErrEntityNotFound{
		ID:    id,
		Scope: scope,
	}
}

func (e ErrEntityNotFound) Error() string {
	return fmt.Sprintf("entity %s not found in %s", e.ID, e.Scope)
}
