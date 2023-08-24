package user

import (
	"context"
	"time"

	"slices"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type UserRepository interface {
	shared.EntityRepository[User]
	UserWithEmail(ctx context.Context, email Email) (User, error)
}

type UserCreatedEvent struct {
	UserId          shared.ID
	Email           Email
	IsEmailVerified bool
}

type User struct {
	id                   shared.ID
	email                Email
	phoneNumber          PhoneNumber
	isEmailVerified      bool
	roles                []UserRole
	authenticationMethod AuthenticationMethod
	signUpDate           time.Time
	domainEvents         shared.Events
}

func NewUser(email Email) User {
	u := User{
		id:           shared.NewID(),
		email:        email,
		roles:        make([]UserRole, 0),
		signUpDate:   time.Now(),
		domainEvents: make(shared.Events),
	}
	u.domainEvents.Emit("user/signedup", u)
	return u
}

func ReconstituteUser(
	id shared.ID,
	email Email,
	phoneNumber PhoneNumber,
	isEmailVerified bool,
	roles []UserRole,
	authenticationMethod AuthenticationMethod,
	signUpDate time.Time) User {
	return User{
		id,
		email,
		phoneNumber,
		isEmailVerified,
		roles,
		authenticationMethod,
		signUpDate,
		make(shared.Events),
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

func (u *User) UserID() shared.ID {
	return u.id
}

func (u *User) VerifyEmail() {
	u.isEmailVerified = true
}

func (u *User) IsEmailVerified() bool {
	return u.isEmailVerified
}

func (u *User) DomainEvents() shared.Events {
	return u.domainEvents
}

func (u *User) AuthenticationMethod() AuthenticationMethod {
	return u.authenticationMethod
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) PhoneNumber() PhoneNumber {
	return u.phoneNumber
}

func (u *User) SetPhoneNumber(phoneNumber PhoneNumber) {
	u.phoneNumber = phoneNumber
}

func (u *User) SignUpDate() time.Time {
	return u.signUpDate
}

func (u *User) HasRole(role UserRole) bool {
	return slices.Contains(u.roles, role)
}

func (u *User) Roles() []UserRole {
	roles := make([]UserRole, len(u.roles))
	copy(roles, u.roles)
	return u.roles
}

func (u *User) addRoleIfNotExists(role UserRole) {
	if !slices.Contains(u.roles, role) {
		u.roles = append(u.roles, role)
	}
}
