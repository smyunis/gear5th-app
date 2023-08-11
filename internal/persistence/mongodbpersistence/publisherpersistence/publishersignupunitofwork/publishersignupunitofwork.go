package publishersignupunitofwork

import (
	"context"
	"fmt"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoDBPublisherSignUpUnitOfWork struct {
	dbStore               mongodbpersistence.MongoDBStore
	db                    *mongo.Database
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	publisherRepository   publisher.PublisherRepository
}

func NewMongoDBPublisherSignUpUnitOfWork(
	dbStore mongodbpersistence.MongoDBStore,
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository,
	publisherRepository publisher.PublisherRepository) MongoDBPublisherSignUpUnitOfWork {
	return MongoDBPublisherSignUpUnitOfWork{
		dbStore,
		dbStore.Database(),
		userRepository,
		managedUserRepository,
		publisherRepository,
	}
}

func (w MongoDBPublisherSignUpUnitOfWork) Save(ctx context.Context, usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error {

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

		err = w.userRepository.Save(sc, usr)
		if err != nil {
			return fmt.Errorf("save user failed: %w", err)
		}
		err = w.managedUserRepository.Save(sc, managedUser)
		if err != nil {
			return fmt.Errorf("save managed user failed : %w", err)
		}
		err = w.publisherRepository.Save(sc, pub)
		if err != nil {
			return fmt.Errorf("save publisher failed : %w", err)
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

func (uow MongoDBPublisherSignUpUnitOfWork) UserRepository() user.UserRepository {
	return uow.userRepository
}
