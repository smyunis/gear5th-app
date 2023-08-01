package user

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

type User struct {
	id                   shared.Id
	firstName            string
	lastName             string
	email                shared.Email
	phoneNumber          shared.PhoneNumber
	isEmailVerified      bool
	roles                []UserRole
	authenticationMethod AuthenticationMethod
	
}

func NewUser(id shared.Id) {

}
