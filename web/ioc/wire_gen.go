// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/identity/tokens"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/mail/identityemail"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/userrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publishersignupunitofwork"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/identitycontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/publishercontrollers"
	"gitlab.com/gear5th/gear5th-app/web/events"
)

// Injectors from dependecyproviders.go:

// API Controllers
func InitManagedUserController() identitycontrollers.UserSignInController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	userSignInController := identitycontrollers.NewUserSignInController(managedUserInteractor, appLogger)
	return userSignInController
}

func InitPublisherSignUpController() publishercontrollers.PublisherSignUpController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	mongoDBPublisherRepository := publisherrepository.NewMongoDBPublisherRepository(mongoDBStoreBootstrap)
	mongoDBPublisherSignUpUnitOfWork := publishersignupunitofwork.NewMongoDBPublisherSignUpUnitOfWork(mongoDBStoreBootstrap, mongoDBUserRepository, mongoDBMangageUserRepository, mongoDBPublisherRepository)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	publisherSignUpInteractor := publisherinteractors.NewPublisherSignUpInteractor(inMemoryEventDispatcher, mongoDBPublisherSignUpUnitOfWork, appLogger)
	publisherSignUpController := publishercontrollers.NewPublisherSignUpController(publisherSignUpInteractor, appLogger)
	return publisherSignUpController
}

func InitRequestPasswordResetController() identitycontrollers.RequestPasswordResetController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	requestPasswordResetController := identitycontrollers.NewRequestPasswordResetController(managedUserInteractor, appLogger)
	return requestPasswordResetController
}

func InitVerifyEmailController() identitycontrollers.VerifyEmailController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	verifyEmailController := identitycontrollers.NewVerifyEmailController(managedUserInteractor, appLogger)
	return verifyEmailController
}

func InitResetPasswordController() identitycontrollers.ResetPasswordController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	resetPasswordController := identitycontrollers.NewResetPasswordController(managedUserInteractor, appLogger)
	return resetPasswordController
}

func InitEventsRegistrar() events.EventHandlerRegistrar {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	verifcationEmailSender := identityemail.NewVerifcationEmailSender(envConfigurationProvider, hs256HMACValidationService, appLogger)
	verificationEmailInteractor := manageduserinteractors.NewVerificationEmailInteractor(verifcationEmailSender, appLogger)
	eventHandlerRegistrar := events.NewEventHandlerRegistrar(inMemoryEventDispatcher, verificationEmailInteractor)
	return eventHandlerRegistrar
}
