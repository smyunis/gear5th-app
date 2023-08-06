package user

import (
	"fmt"
	"strings"
)

type PersonName struct {
	firstName string
	lastName  string
}

func NewPersonName(firstName, lastName string) PersonName {
	return PersonName{
		firstName: firstName,
		lastName:  lastName,
	}
}

func NewPersonNameWithFullName(fullName string) PersonName {
	names := strings.Split(fullName, " ")
	return PersonName{
		firstName: names[0],
		lastName:  names[1],
	}
}

func (p *PersonName) FullName() string {
	return fmt.Sprintf("%s %s", p.firstName, p.lastName)
}
