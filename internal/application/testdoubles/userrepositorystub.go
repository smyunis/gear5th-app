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

func (UserRepositoryStub) Get(id shared.Id) (user.User, error) {

	if id != shared.Id("stub-id-xxx") {
		return user.User{}, shared.NewEntityNotFoundError(id.String(), "user")
	}

	u := &user.User{}
	uVal := reflect.ValueOf(u).Elem()

	idField := uVal.FieldByName("id")
	idField = reflect.NewAt(idField.Type(), unsafe.Pointer(idField.UnsafeAddr())).Elem()
	idField.Set(reflect.ValueOf(id))

	emailField := uVal.FieldByName("email")
	emailField = reflect.NewAt(emailField.Type(), unsafe.Pointer(emailField.UnsafeAddr())).Elem()
	mymail, _ := user.NewEmail("mymail@gmail.com")
	emailField.Set(reflect.ValueOf(mymail))

	u.VerifyEmail()

	return *u, nil
}

func (UserRepositoryStub) Save(u user.User) error {
	return nil
}

func (usr UserRepositoryStub) UserWithEmail(email user.Email) (user.User, error) {
	if mymail, _ := user.NewEmail("mymail@gmail.com"); mymail == email {
		stubId := shared.Id("stub-id-xxx")
		usr, err := usr.Get(stubId)
		return usr, err
	}
	return user.User{}, shared.NewEntityNotFoundError(email.Email(), "user")
}

type ManagedUserRepositoryStub struct{}

func (ManagedUserRepositoryStub) Get(id shared.Id) (user.ManagedUser, error) {

	if id != shared.Id("stub-id-xxx") {
		return user.ManagedUser{}, shared.NewEntityNotFoundError(id.String(), "user")
	}

	u := &user.ManagedUser{}
	uVal := reflect.ValueOf(u).Elem()

	userIdField := uVal.FieldByName("userId")
	userIdField = reflect.NewAt(userIdField.Type(), unsafe.Pointer(userIdField.UnsafeAddr())).Elem()
	userIdField.Set(reflect.ValueOf(shared.NewId()))

	u.SetPassword("gokuisking")

	return *u, nil
}

func (ManagedUserRepositoryStub) Save(u user.ManagedUser) error {
	return nil
}
