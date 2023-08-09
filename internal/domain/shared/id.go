package shared

import "github.com/oklog/ulid/v2"

type ID string

func NewID() ID {
	newUlid := ulid.Make().String()
	return ID(newUlid)
}

func (i ID) String() string {
	return string(i)
}
