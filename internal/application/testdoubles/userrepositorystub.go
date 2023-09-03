package testdoubles

import (
	"context"
	"reflect"
	"unsafe"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type UserRepositoryStub struct{}

func NewUserRepositoryStub() UserRepositoryStub {
	return UserRepositoryStub{}
}

func (UserRepositoryStub) Get(ctx context.Context, id shared.ID) (user.User, error) {

	if id != shared.ID("stub-id-xxx") {
		return user.User{}, application.ErrEntityNotFound
	}

	u := &user.User{}
	setStructField[user.User, shared.ID](u, "id", id)
	mymail, _ := user.NewEmail("mymail@gmail.com")
	setStructField[user.User, user.Email](u, "email", mymail)
	setStructField[user.User, []user.UserRole](u, "roles", []user.UserRole{user.Publisher})
	u.VerifyEmail()
	return *u, nil
}

func (UserRepositoryStub) Save(ctx context.Context, u user.User) error {
	return nil
}

func (usr UserRepositoryStub) UserWithEmail(ctx context.Context, email user.Email) (user.User, error) {
	if mymail, _ := user.NewEmail("mymail@gmail.com"); mymail == email {
		stubId := shared.ID("stub-id-xxx")
		usr, err := usr.Get(context.Background(), stubId)
		return usr, err
	}
	if somemail, _ := user.NewEmail("somemail@gmail.com"); somemail == email {
		stubId := shared.ID("stub-id-xxx")
		usr, err := usr.Get(context.Background(), stubId)

		setStructField[user.User, user.Email](&usr, "email", somemail)
		setStructField[user.User, bool](&usr, "isEmailVerified", false)

		return usr, err
	}
	return user.User{}, application.ErrEntityNotFound
}

type ManagedUserRepositoryStub struct{}

func (ManagedUserRepositoryStub) Get(ctx context.Context, id shared.ID) (user.ManagedUser, error) {

	if id != shared.ID("stub-id-xxx") {
		return user.ManagedUser{}, application.ErrEntityNotFound
	}

	u := &user.ManagedUser{}
	setStructField[user.ManagedUser, shared.ID](u, "userId", shared.NewID())
	u.SetPassword("gokuisking")

	return *u, nil
}

func (ManagedUserRepositoryStub) Save(ctx context.Context, u user.ManagedUser) error {
	return nil
}

func setStructField[T, V any](struc *T, field string, value V) {
	uVal := reflect.ValueOf(struc).Elem()
	structField := uVal.FieldByName(field)
	structField = reflect.NewAt(structField.Type(), unsafe.Pointer(structField.UnsafeAddr())).Elem()
	structField.Set(reflect.ValueOf(value))
}

type OAuthUserRepositoryStub struct{}

// UserWithAccountIdentifier implements user.OAuthUserRepository.
func (OAuthUserRepositoryStub) UserWithAccountIdentifier(accountID string) (user.OAuthUser, error) {
	return user.ReconstituteOAuthUser(shared.NewID(), accountID, user.GoogleOAuth), nil
}

// Get implements user.OAuthUserRepository.
func (OAuthUserRepositoryStub) Get(ctx context.Context, id shared.ID) (user.OAuthUser, error) {
	return user.ReconstituteOAuthUser(id, "xxx-yyy", user.GoogleOAuth), nil
}

// Save implements user.OAuthUserRepository.
func (OAuthUserRepositoryStub) Save(ctx context.Context, e user.OAuthUser) error {
	return nil
}
