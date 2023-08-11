// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"gitlab.com/gear5th/gear5th-api/cmd/api/controllers/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/cmd/api/controllers/publishercontrollers"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/identity/accesstoken"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/keyvaluestore/rediskeyvaluestore"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/mail/identityemail"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence/userrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/publisherpersistence/publishersignupunitofwork"
)

// Injectors from dependecyproviders.go:

func InitManagedUserController() identitycontrollers.ManagedUserController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := accesstoken.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider)
	redisBootstrapper := rediskeyvaluestore.NewRedisBootstrapper(envConfigurationProvider)
	redisKeyValueStore := rediskeyvaluestore.NewRedisKeyValueStore(redisBootstrapper)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, redisKeyValueStore)
	managedUserController := identitycontrollers.NewManagedUserController(managedUserInteractor)
	return managedUserController
}

func InitPublisherSignUpController() publishercontrollers.PublisherSignUpController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	mongoDBPublisherRepository := publisherrepository.NewMongoDBPublisherRepository(mongoDBStoreBootstrap)
	mongoDBPublisherSignUpUnitOfWork := publishersignupunitofwork.NewMongoDBPublisherSignUpUnitOfWork(mongoDBStoreBootstrap, mongoDBUserRepository, mongoDBMangageUserRepository, mongoDBPublisherRepository)
	publisherSignUpInteractor := publisherinteractors.NewPublisherSignUpInteractor(mongoDBPublisherSignUpUnitOfWork)
	publisherSignUpController := publishercontrollers.NewPublisherSignUpController(publisherSignUpInteractor)
	return publisherSignUpController
}

func InitRequestPasswordResetController() identitycontrollers.RequestPasswordResetController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := accesstoken.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider)
	redisBootstrapper := rediskeyvaluestore.NewRedisBootstrapper(envConfigurationProvider)
	redisKeyValueStore := rediskeyvaluestore.NewRedisKeyValueStore(redisBootstrapper)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, redisKeyValueStore)
	requestPasswordResetController := identitycontrollers.NewRequestPasswordResetController(managedUserInteractor)
	return requestPasswordResetController
}
