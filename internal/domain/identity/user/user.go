package user

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type User struct {
	id                   shared.Id
	email                shared.Email
	roles                []UserRole
	authenticationMethod AuthenticationMethod
}

func NewUser(email shared.Email, authNMethod AuthenticationMethod) User {
	return User{
		id:                   shared.NewId(),
		email:                email,
		authenticationMethod: authNMethod,
	}
}

func (u *User) AsManagedUser(fullName string, phoneNumber shared.PhoneNumber) ManagedUser {
	u.authenticationMethod = Managed
	return ManagedUser{
		User:        *u,
		fullName:    fullName,
		phoneNumber: phoneNumber,
	}
}

func (u *User) AsOAuthUser(userIdentifier any, provider OAuthProvider) OAuthUser {
	u.authenticationMethod = OAuth
	return OAuthUser{
		User:           *u,
		userIdentifier: userIdentifier,
		oauthProvider:  provider,
	}
}

func (u *User) SignUpPublisher() publisher.Publisher {
	u.addRoleIfNotExists(u.roles, Publisher)
	return publisher.NewPublisher(u.id)
}

func (u *User) AuthenticationMethod() AuthenticationMethod {
	return u.authenticationMethod
}

func (u *User) addRoleIfNotExists(roles []UserRole, item UserRole) []UserRole {
	for _, v := range roles {
		if v == item {
			return roles
		}
	}
	return append(roles, item)
}
