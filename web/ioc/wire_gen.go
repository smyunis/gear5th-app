// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adslotinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/siteinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/identity/tokens"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/mail/identityemail"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/siteverification"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/adslotpersistence/adslotrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/userrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publishersignupunitofwork"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/sitepersistence/siterepository"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/adslotcontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/homecontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/identitycontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/publishercontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/sitecontrollers"
	"gitlab.com/gear5th/gear5th-app/web/events"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

// Injectors from dependecyproviders.go:

func InitJwtAuthenticationMiddleware() middlewares.JwtAuthenticationMiddleware {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	return jwtAuthenticationMiddleware
}

// Controllers
func InitManagedUserController() identitycontrollers.UserSignInController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := identityinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenService, requestPassordResetEmailService, hs256HMACValidationService)
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
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := identityinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenService, requestPassordResetEmailService, hs256HMACValidationService)
	requestPasswordResetController := identitycontrollers.NewRequestPasswordResetController(managedUserInteractor, appLogger)
	return requestPasswordResetController
}

func InitVerifyEmailController() identitycontrollers.VerifyEmailController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := identityinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenService, requestPassordResetEmailService, hs256HMACValidationService)
	verifyEmailController := identitycontrollers.NewVerifyEmailController(managedUserInteractor, appLogger)
	return verifyEmailController
}

func InitResetPasswordController() identitycontrollers.ResetPasswordController {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBMangageUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(mongoDBStoreBootstrap)
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	requestPassordResetEmailService := identityemail.NewRequestPassordResetEmailService(envConfigurationProvider, hs256HMACValidationService, appLogger)
	managedUserInteractor := identityinteractors.NewManagedUserInteractor(inMemoryEventDispatcher, mongoDBUserRepository, mongoDBMangageUserRepository, jwtAccessTokenService, requestPassordResetEmailService, hs256HMACValidationService)
	resetPasswordController := identitycontrollers.NewResetPasswordController(managedUserInteractor, appLogger)
	return resetPasswordController
}

func InitHomeController() homecontrollers.HomeController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	homeController := homecontrollers.NewHomeController(jwtAuthenticationMiddleware)
	return homeController
}

func InitSiteController() sitecontrollers.SiteController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	mongoDBSiteRepository := siterepository.NewMongoDBSiteRepository(mongoDBStoreBootstrap, appLogger)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	appHTTPClient := infrastructure.NewAppHTTPClient()
	adsTxtVerificationService := siteverification.NewAdsTxtVerificationService(appHTTPClient, appLogger)
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	siteInteractor := siteinteractors.NewSiteInteractor(mongoDBSiteRepository, mongoDBUserRepository, adsTxtVerificationService, inMemoryEventDispatcher, appLogger)
	siteController := sitecontrollers.NewSiteController(jwtAuthenticationMiddleware, siteInteractor, appLogger)
	return siteController
}

func InitCreateSiteController() sitecontrollers.CreateSiteController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	mongoDBSiteRepository := siterepository.NewMongoDBSiteRepository(mongoDBStoreBootstrap, appLogger)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	appHTTPClient := infrastructure.NewAppHTTPClient()
	adsTxtVerificationService := siteverification.NewAdsTxtVerificationService(appHTTPClient, appLogger)
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	siteInteractor := siteinteractors.NewSiteInteractor(mongoDBSiteRepository, mongoDBUserRepository, adsTxtVerificationService, inMemoryEventDispatcher, appLogger)
	createSiteController := sitecontrollers.NewCreateSiteController(jwtAuthenticationMiddleware, siteInteractor, appLogger)
	return createSiteController
}

func InitVerifySiteController() sitecontrollers.VerifySiteController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	mongoDBSiteRepository := siterepository.NewMongoDBSiteRepository(mongoDBStoreBootstrap, appLogger)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	appHTTPClient := infrastructure.NewAppHTTPClient()
	adsTxtVerificationService := siteverification.NewAdsTxtVerificationService(appHTTPClient, appLogger)
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	siteInteractor := siteinteractors.NewSiteInteractor(mongoDBSiteRepository, mongoDBUserRepository, adsTxtVerificationService, inMemoryEventDispatcher, appLogger)
	verifySiteController := sitecontrollers.NewVerifySiteController(jwtAuthenticationMiddleware, siteInteractor, appLogger)
	return verifySiteController
}

func InitAdSlotController() adslotcontrollers.AdSlotController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	mongoDBSiteRepository := siterepository.NewMongoDBSiteRepository(mongoDBStoreBootstrap, appLogger)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBAdSlotRepository := adslotrepository.NewMongoDBAdSlotRepository(mongoDBStoreBootstrap, appLogger)
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	adSlotInteractor := adslotinteractors.NewAdSlotInteractor(mongoDBSiteRepository, mongoDBUserRepository, mongoDBAdSlotRepository, inMemoryEventDispatcher)
	adSlotController := adslotcontrollers.NewAdSlotController(jwtAuthenticationMiddleware, adSlotInteractor, appLogger)
	return adSlotController
}

func InitCreateAdSlotController() adslotcontrollers.CreateAdSlotController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	mongoDBSiteRepository := siterepository.NewMongoDBSiteRepository(mongoDBStoreBootstrap, appLogger)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBAdSlotRepository := adslotrepository.NewMongoDBAdSlotRepository(mongoDBStoreBootstrap, appLogger)
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	adSlotInteractor := adslotinteractors.NewAdSlotInteractor(mongoDBSiteRepository, mongoDBUserRepository, mongoDBAdSlotRepository, inMemoryEventDispatcher)
	appHTTPClient := infrastructure.NewAppHTTPClient()
	adsTxtVerificationService := siteverification.NewAdsTxtVerificationService(appHTTPClient, appLogger)
	siteInteractor := siteinteractors.NewSiteInteractor(mongoDBSiteRepository, mongoDBUserRepository, adsTxtVerificationService, inMemoryEventDispatcher, appLogger)
	createAdSlotController := adslotcontrollers.NewCreateAdSlotController(jwtAuthenticationMiddleware, adSlotInteractor, siteInteractor, appLogger)
	return createAdSlotController
}

func InitEditAdSlotController() adslotcontrollers.EditAdSlotController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	mongoDBSiteRepository := siterepository.NewMongoDBSiteRepository(mongoDBStoreBootstrap, appLogger)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBAdSlotRepository := adslotrepository.NewMongoDBAdSlotRepository(mongoDBStoreBootstrap, appLogger)
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	adSlotInteractor := adslotinteractors.NewAdSlotInteractor(mongoDBSiteRepository, mongoDBUserRepository, mongoDBAdSlotRepository, inMemoryEventDispatcher)
	editAdSlotController := adslotcontrollers.NewEditAdSlotController(jwtAuthenticationMiddleware, adSlotInteractor, appLogger)
	return editAdSlotController
}

func InitAdSlotIntegrationSnippetController() adslotcontrollers.AdSlotIntegrationSnippetController {
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	jwtAccessTokenService := tokens.NewJwtAccessTokenService(envConfigurationProvider)
	jwtAuthenticationMiddleware := middlewares.NewJwtAuthenticationMiddleware(jwtAccessTokenService)
	mongoDBStoreBootstrap := mongodbpersistence.NewMongoDBStoreBootstrap(envConfigurationProvider)
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	mongoDBSiteRepository := siterepository.NewMongoDBSiteRepository(mongoDBStoreBootstrap, appLogger)
	mongoDBUserRepository := userrepository.NewMongoDBUserRepository(mongoDBStoreBootstrap)
	mongoDBAdSlotRepository := adslotrepository.NewMongoDBAdSlotRepository(mongoDBStoreBootstrap, appLogger)
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	adSlotInteractor := adslotinteractors.NewAdSlotInteractor(mongoDBSiteRepository, mongoDBUserRepository, mongoDBAdSlotRepository, inMemoryEventDispatcher)
	adSlotIntegrationSnippetController := adslotcontrollers.NewAdSlotIntegrationSnippetController(jwtAuthenticationMiddleware, adSlotInteractor, appLogger)
	return adSlotIntegrationSnippetController
}

func InitEventsRegistrar() events.EventHandlerRegistrar {
	inMemoryEventDispatcher := application.NewAppEventDispatcher()
	envConfigurationProvider := infrastructure.NewEnvConfigurationProvider()
	hs256HMACValidationService := tokens.NewHS256HMACValidationService()
	appLogger := infrastructure.NewAppLogger(envConfigurationProvider)
	verifcationEmailSender := identityemail.NewVerifcationEmailSender(envConfigurationProvider, hs256HMACValidationService, appLogger)
	verificationEmailInteractor := identityinteractors.NewVerificationEmailInteractor(verifcationEmailSender, appLogger)
	eventHandlerRegistrar := events.NewEventHandlerRegistrar(inMemoryEventDispatcher, verificationEmailInteractor)
	return eventHandlerRegistrar
}
