package user

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"golang.org/x/exp/slices"
)

type UserRepository interface {
	shared.EntityRepository[User]
	UserWithEmail(email Email) (User, error)
}

type User struct {
	id                   shared.Id
	email                Email
	phoneNumber          PhoneNumber
	roles                []UserRole
	authenticationMethod AuthenticationMethod
}

func NewUser(email Email) User {
	return User{
		id:    shared.NewId(),
		email: email,
	}
}

func (u *User) AsManagedUser(name PersonName, password string) ManagedUser {
	u.authenticationMethod = Managed
	managedUser := ManagedUser{
		userId: u.id,
		name:   name,
	}
	managedUser.SetPassword(password)
	return managedUser
}

func (u *User) AsOAuthUser(userIdentifier string, provider OAuthProvider) OAuthUser {
	u.authenticationMethod = OAuth
	return OAuthUser{
		userId:         u.id,
		userIdentifier: userIdentifier,
		oauthProvider:  provider,
	}
}

func (u *User) SignUpPublisher() publisher.Publisher {
	u.addRoleIfNotExists(Publisher)
	return publisher.NewPublisher(u.id)
}

func (u *User) UserID() shared.Id {
	return u.id
}

func (u *User) AuthenticationMethod() AuthenticationMethod {
	return u.authenticationMethod
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) SetPhoneNumber(phoneNumber PhoneNumber) {
	u.phoneNumber = phoneNumber
}

func (u *User) HasRole(role UserRole) bool {
	return slices.Contains(u.roles, role)
}

func (u *User) addRoleIfNotExists(role UserRole) {
	if !slices.Contains(u.roles, role) {
		u.roles = append(u.roles, role)
	}
}
