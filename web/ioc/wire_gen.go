// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/identity/tokens"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/mail/identityemail"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence/userrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/publisherpersistence/publishersignupunitofwork"
	"gitlab.com/gear5th/gear5th-api/web/controllers/publish/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/web/controllers/publish/publishercontrollers"
)

// Injectors from dependecyproviders.go:

// API Controllers
func InitManagedUserController() identitycontrollers.UserSignInController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	userSignInController := identitycontrollers.NewUserSignInController(managedUserInteractor)
	return userSignInController
}

func InitPublisherSignUpController() publishercontrollers.PublisherSignUpController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	mongoDBPublisherRepository := publisherrepository.NewMongoDBPublisherRepository(mongoDBStoreBootstrap)
	mongoDBPublisherSignUpUnitOfWork := publishersignupunitofwork.NewMongoDBPublisherSignUpUnitOfWork(mongoDBStoreBootstrap, mongoDBUserRepository, mongoDBMangageUserRepository, mongoDBPublisherRepository)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	verifcationEmailSender := identityemail.NewVerifcationEmailSender(envConfigurationProvider, hs256HMACValidationService)
	publisherSignUpInteractor := publisherinteractors.NewPublisherSignUpInteractor(mongoDBPublisherSignUpUnitOfWork, verifcationEmailSender)
	publisherSignUpController := publishercontrollers.NewPublisherSignUpController(publisherSignUpInteractor)
	return publisherSignUpController
}

func InitRequestPasswordResetController() identitycontrollers.RequestPasswordResetController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	requestPasswordResetController := identitycontrollers.NewRequestPasswordResetController(managedUserInteractor)
	return requestPasswordResetController
}

func InitVerifyEmailController() identitycontrollers.VerifyEmailController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	verifyEmailController := identitycontrollers.NewVerifyEmailController(managedUserInteractor)
	return verifyEmailController
}

func InitResetPasswordController() identitycontrollers.ResetPasswordController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenGenerator := tokens.NewJwtAccessTokenGenenrator(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService)
	managedUserInteractor := manageduserinteractors.NewManagedUserInteractor(mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenGenerator, requestPassordResetEmailService, hs256HMACValidationService)
	resetPasswordController := identitycontrollers.NewResetPasswordController(managedUserInteractor)
	return resetPasswordController
}

func InitVerifcationEmailSender() identityemail.VerifcationEmailSender {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	verifcationEmailSender := identityemail.NewVerifcationEmailSender(envConfigurationProvider, hs256HMACValidationService)
	return verifcationEmailSender
}
