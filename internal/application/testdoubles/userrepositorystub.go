package testdoubles

import (
	"reflect"
	"unsafe"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type UserRepositoryStub struct{}

func NewUserRepositoryStub() UserRepositoryStub {
	return UserRepositoryStub{}
}

func (UserRepositoryStub) Get(id shared.ID) (user.User, error) {

	if id != shared.ID("stub-id-xxx") {
		return user.User{}, shared.NewEntityNotFoundError(id.String(), "user")
	}

	u := &user.User{}
	setStructField[user.User, shared.ID](u, "id", id)
	mymail, _ := user.NewEmail("mymail@gmail.com")
	setStructField[user.User, user.Email](u, "email", mymail)
	u.VerifyEmail()
	return *u, nil
}

func (UserRepositoryStub) Save(u user.User) error {
	return nil
}

func (usr UserRepositoryStub) UserWithEmail(email user.Email) (user.User, error) {
	if mymail, _ := user.NewEmail("mymail@gmail.com"); mymail == email {
		stubId := shared.ID("stub-id-xxx")
		usr, err := usr.Get(stubId)
		return usr, err
	}
	if somemail, _ := user.NewEmail("somemail@gmail.com"); somemail == email {
		stubId := shared.ID("stub-id-xxx")
		usr, err := usr.Get(stubId)

		setStructField[user.User, user.Email](&usr, "email", somemail)
		setStructField[user.User, bool](&usr, "isEmailVerified", false)

		return usr, err
	}
	return user.User{}, shared.NewEntityNotFoundError(email.Email(), "user")
}

type ManagedUserRepositoryStub struct{}

func (ManagedUserRepositoryStub) Get(id shared.ID) (user.ManagedUser, error) {

	if id != shared.ID("stub-id-xxx") {
		return user.ManagedUser{}, shared.NewEntityNotFoundError(id.String(), "user")
	}

	u := &user.ManagedUser{}
	setStructField[user.ManagedUser, shared.ID](u, "userId", shared.NewID())
	u.SetPassword("gokuisking")

	return *u, nil
}

func (ManagedUserRepositoryStub) Save(u user.ManagedUser) error {
	return nil
}

func setStructField[T, V any](struc *T, field string, value V) {
	uVal := reflect.ValueOf(struc).Elem()
	structField := uVal.FieldByName(field)
	structField = reflect.NewAt(structField.Type(), unsafe.Pointer(structField.UnsafeAddr())).Elem()
	structField.Set(reflect.ValueOf(value))
}
