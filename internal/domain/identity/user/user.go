package user

import (
	"context"
	"time"

	"slices"

	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/advertiser"
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
	ID                   shared.ID
	Email                Email
	PhoneNumber          PhoneNumber
	IsEmailVerified      bool
	Roles                []UserRole
	AuthenticationMethod AuthenticationMethod
	SignUpDate           time.Time
	Events               shared.Events
}

func NewUser(email Email) User {
	u := User{
		ID:         shared.NewID(),
		Email:      email,
		Roles:      make([]UserRole, 0),
		SignUpDate: time.Now(),
		Events:     make(shared.Events),
	}
	u.Events.Emit("user/signedup", u)
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
	u.AuthenticationMethod = Managed
	managedUser := ManagedUser{
		userId: u.ID,
		name:   name,
	}
	managedUser.SetPassword(password)
	return managedUser
}

func (u *User) AsOAuthUser(userAccountID string, provider OAuthProvider) OAuthUser {
	u.AuthenticationMethod = OAuth
	return OAuthUser{
		userID:        u.ID,
		userAccountID: userAccountID,
		oauthProvider: provider,
	}
}

func (u *User) SignUpPublisher() publisher.Publisher {
	u.addRoleIfNotExists(Publisher)
	return publisher.NewPublisher(u.ID)
}

func (u *User) SignUpAdvertiser(name string) advertiser.Advertiser {
	u.addRoleIfNotExists(Advertiser)
	return advertiser.NewAdvertiser(u.ID, name)
}

func (u *User) VerifyEmail() {
	u.IsEmailVerified = true
}

func (u *User) HasRole(role UserRole) bool {
	return slices.Contains(u.Roles, role)
}

func (u *User) addRoleIfNotExists(role UserRole) {
	if !slices.Contains(u.Roles, role) {
		u.Roles = append(u.Roles, role)
	}
}
