//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/identity/tokens"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/keyvaluestore/rediskeyvaluestore"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/mail/identityemail"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/userrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publishersignupunitofwork"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/identitycontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/publishercontrollers"
)

var Container wire.ProviderSet = wire.NewSet(
	//Stub Repositories
	wire.Struct(new(testdoubles.UserRepositoryStub), "*"),
	wire.Struct(new(testdoubles.ManagedUserRepositoryStub), "*"),
	wire.Struct(new(testdoubles.PublisherRepositoryStub), "*"),
	wire.Struct(new(testdoubles.PublisherSignUpUnitOfWorkStub), "*"),
	wire.Struct(new(testdoubles.RequestResetPasswordEmailStub), "*"),

	// Repositories with testdouble stubs
	// wire.Bind(new(user.UserRepository), new(testdoubles.UserRepositoryStub)),
	// wire.Bind(new(user.ManagedUserRepository), new(testdoubles.ManagedUserRepositoryStub)),
	// wire.Bind(new(publisher.PublisherRepository), new(testdoubles.PublisherRepositoryStub)),
	// wire.Bind(new(publisherinteractors.PublisherSignUpUnitOfWork), new(testdoubles.PublisherSignUpUnitOfWorkStub)),

	//MongoDB persistence repositores
	mongodbpersistence.NewMongoDBStoreBootstrap,
	wire.Bind(new(mongodbpersistence.MongoDBStore), new(mongodbpersistence.MongoDBStoreBootstrap)),
	userrepository.NewMongoDBUserRepository,
	wire.Bind(new(user.UserRepository), new(userrepository.MongoDBUserRepository)),
	manageduserrepository.NewMongoDBMangageUserRepository,
	wire.Bind(new(user.ManagedUserRepository), new(manageduserrepository.MongoDBMangageUserRepository)),
	publisherrepository.NewMongoDBPublisherRepository,
	wire.Bind(new(publisher.PublisherRepository), new(publisherrepository.MongoDBPublisherRepository)),
	publishersignupunitofwork.NewMongoDBPublisherSignUpUnitOfWork,
	wire.Bind(new(publisherinteractors.PublisherSignUpUnitOfWork), new(publishersignupunitofwork.MongoDBPublisherSignUpUnitOfWork)),

	//Infrastructures
	tokens.NewJwtAccessTokenGenenrator,
	infrastructure.NewEnvConfigurationProvider,
	tokens.NewHS256HMACValidationService,
	identityemail.NewRequestPassordResetEmailService,
	wire.Bind(new(identityinteractors.AccessTokenGenerator), new(tokens.JwtAccessTokenGenerator)),
	wire.Bind(new(identityinteractors.DigitalSignatureService), new(tokens.HS256HMACValidationService)),

	wire.Bind(new(infrastructure.ConfigurationProvider), new(infrastructure.EnvConfigurationProvider)),
	// wire.Bind(new(manageduserinteractors.RequestPasswordResetEmailService), new(testdoubles.RequestResetPasswordEmailStub)),
	wire.Bind(new(manageduserinteractors.RequestPasswordResetEmailService), new(identityemail.RequestPassordResetEmailService)),
	wire.Bind(new(publisherinteractors.VerificationEmailService), new(identityemail.VerifcationEmailSender)),
	
	wire.Bind(new(application.Logger), new(infrastructure.AppLogger)),
	infrastructure.NewAppLogger,
	
	identityemail.NewVerifcationEmailSender,

	rediskeyvaluestore.NewRedisBootstrapper,
	rediskeyvaluestore.NewRedisKeyValueStore,
	wire.Bind(new(application.KeyValueStore), new(rediskeyvaluestore.RedisKeyValueStore)),

	//Interactors
	manageduserinteractors.NewManagedUserInteractor,
	publisherinteractors.NewPublisherSignUpInteractor,

	//Controllers
	publishercontrollers.NewPublisherSignUpController,
	identitycontrollers.NewUserSignInController,
	identitycontrollers.NewVerifyEmailController,
	identitycontrollers.NewResetPasswordController,
	identitycontrollers.NewRequestPasswordResetController)
