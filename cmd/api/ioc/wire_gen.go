// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"gitlab.com/gear5th/gear5th-api/cmd/api/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/internal/application/identity/usersignin"
	"gitlab.com/gear5th/gear5th-api/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/identity/accesstoken"
)

// Injectors from dependecyproviders.go:

func InitManagedUserController() identitycontrollers.ManagedUserController {
	userRepositoryStub := testdoubles.UserRepositoryStub{}
	managedUserRepositoryStub := testdoubles.ManagedUserRepositoryStub{}
	jwtAccessTokenGenenrator := accesstoken.NewJwtAccessTokenGenenrator()
	managedUserInteractor := usersignin.NewManagedUserInteractor(userRepositoryStub, managedUserRepositoryStub, jwtAccessTokenGenenrator)
	managedUserController := identitycontrollers.NewManagedUserController(managedUserInteractor)
	return managedUserController
}
