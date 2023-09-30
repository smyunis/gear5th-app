package advertiserrepository

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/advertiser"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoDBAdvertiserSignUpUnitOfWork struct {
	dbStore             mongodbpersistence.MongoDBStore
	db                  *mongo.Database
	advertierRepository advertiser.AdvertiserRepository
	userRepository      user.UserRepository
	logger              application.Logger
}

func NewMongoDBAdvertiserSignUpUnitOfWork(
	advertierRepository advertiser.AdvertiserRepository,
	userRepository user.UserRepository,
	dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBAdvertiserSignUpUnitOfWork {
	return MongoDBAdvertiserSignUpUnitOfWork{
		dbStore:             dbStore,
		db:                  dbStore.Database(),
		advertierRepository: advertierRepository,
		userRepository:      userRepository,
		logger:              logger,
	}
}

func (w MongoDBAdvertiserSignUpUnitOfWork) Save(ctx context.Context, u user.User, a advertiser.Advertiser) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)
	client := w.db.Client()
	session, err := client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(context.Background())

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err := session.StartTransaction(txnOptions)
		if err != nil {
			return err
		}

		err = w.userRepository.Save(sc, u)
		if err != nil {
			return err
		}
		err = w.advertierRepository.Save(sc, a)
		if err != nil {
			return err
		}

		err = session.CommitTransaction(sc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		session.AbortTransaction(context.Background())
		return err
	}

	return nil
}
