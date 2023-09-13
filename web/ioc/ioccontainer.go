//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/identity/googleoauth"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/identity/tokens"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/keyvaluestore/fastcachekeyvaluestore"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/mail/identityemail"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/siteverification"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/advertiserpersistence/adpiecerepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/advertiserpersistence/campaignrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/filestore"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/oauthuserrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/identitypersistence/userrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/adslotrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/publishersignupunitofwork"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/siterepository"
	"gitlab.com/gear5th/gear5th-app/web/controllers/advertiser/adpiececontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/advertiser/campaigncontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/accountcontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/adslotcontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/homecontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/identitycontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/publishercontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/sitecontrollers"
	"gitlab.com/gear5th/gear5th-app/web/events"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var Container wire.ProviderSet = wire.NewSet(

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
	siterepository.NewMongoDBSiteRepository,
	wire.Bind(new(site.SiteRepository), new(siterepository.MongoDBSiteRepository)),
	adslotrepository.NewMongoDBAdSlotRepository,
	wire.Bind(new(adslot.AdSlotRepository), new(adslotrepository.MongoDBAdSlotRepository)),
	oauthuserrepository.NewMongoDBOAuthUserRepository,
	wire.Bind(new(user.OAuthUserRepository), new(oauthuserrepository.MongoDBOAuthUserRepository)),
	filestore.NewMongoDBGridFSFileStore,
	wire.Bind(new(application.FileStore), new(filestore.MongoDBGridFSFileStore)),
	campaignrepository.NewMongoDBCampaignRepository,
	wire.Bind(new(campaign.CampaignRepository), new(campaignrepository.MongoDBCampaignRepository)),
	adpiecerepository.NewMongoDBAdPieceRepository,
	wire.Bind(new(adpiece.AdPieceRepository), new(adpiecerepository.MongoDBAdPieceRepository)),

	//Infrastructures
	infrastructure.NewAppHTTPClient,
	wire.Bind(new(infrastructure.HTTPClient), new(infrastructure.AppHTTPClient)),

	infrastructure.NewEnvConfigurationProvider,
	wire.Bind(new(infrastructure.ConfigurationProvider), new(infrastructure.EnvConfigurationProvider)),
	tokens.NewHS256HMACValidationService,
	wire.Bind(new(identityinteractors.DigitalSignatureService), new(tokens.HS256HMACValidationService)),
	tokens.NewJwtAccessTokenService,
	wire.Bind(new(identityinteractors.AccessTokenService), new(tokens.JwtAccessTokenService)),
	siteverification.NewAdsTxtVerificationService,
	wire.Bind(new(site.AdsTxtVerificationService), new(siteverification.AdsTxtVerificationService)),

	googleoauth.NewGoogleOAuthService,
	wire.Bind(new(user.GoogleOAuthService), new(googleoauth.GoogleOAuthServiceImpl)),

	// Mail
	identityemail.NewVerifcationEmailSender,
	identityemail.NewRequestPassordResetEmailService,
	wire.Bind(new(identityinteractors.RequestPasswordResetEmailService), new(identityemail.RequestPassordResetEmailService)),
	wire.Bind(new(identityinteractors.VerificationEmailService), new(identityemail.VerifcationEmailSender)),

	// Logger
	wire.Bind(new(application.Logger), new(infrastructure.AppLogger)),
	infrastructure.NewAppLogger,

	// Redis
	// rediskeyvaluestore.NewRedisBootstrapper,
	// rediskeyvaluestore.NewRedisKeyValueStore,
	// wire.Bind(new(application.KeyValueStore), new(rediskeyvaluestore.RedisKeyValueStore)),

	// FastCache
	fastcachekeyvaluestore.NewFastCacheKeyValueStore,
	wire.Bind(new(application.KeyValueStore), new(fastcachekeyvaluestore.FastCacheKeyValueStore)),

	//Interactors
	identityinteractors.NewManagedUserInteractor,
	identityinteractors.NewVerificationEmailInteractor,
	identityinteractors.NewOAuthUserInteractor,
	identityinteractors.NewUserAccountInteractor,
	publisherinteractors.NewPublisherSignUpInteractor,
	publisherinteractors.NewSiteInteractor,
	publisherinteractors.NewAdSlotInteractor,
	advertiserinteractors.NewAdPieceInteractor,
	advertiserinteractors.NewCampaignInteractor,

	//Middlewares
	middlewares.NewJwtAuthenticationMiddleware,
	middlewares.NewAdvertiserRefferalMiddleware,

	//Controllers
	homecontrollers.NewHomeController,
	publishercontrollers.NewPublisherSignUpController,
	identitycontrollers.NewUserSignInController,
	identitycontrollers.NewVerifyEmailController,
	identitycontrollers.NewResetPasswordController,
	identitycontrollers.NewRequestPasswordResetController,
	identitycontrollers.NewOAuthSignInController,
	sitecontrollers.NewSiteController,
	sitecontrollers.NewCreateSiteController,
	sitecontrollers.NewVerifySiteController,
	adslotcontrollers.NewAdSlotController,
	adslotcontrollers.NewCreateAdSlotController,
	adslotcontrollers.NewEditAdSlotController,
	adslotcontrollers.NewAdSlotIntegrationSnippetController,
	accountcontrollers.NewAccountController,
	adpiececontrollers.NewAdPieceController,
	campaigncontrollers.NewCampaignController,

	application.NewAppEventDispatcher,
	wire.Bind(new(application.EventDispatcher), new(application.InMemoryEventDispatcher)),
	events.NewEventHandlerRegistrar)
