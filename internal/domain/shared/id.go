package shared

import "github.com/oklog/ulid/v2"

type Id string

func (i Id) New() Id {
	newUlid := ulid.Make().String()
	return Id(newUlid)
}