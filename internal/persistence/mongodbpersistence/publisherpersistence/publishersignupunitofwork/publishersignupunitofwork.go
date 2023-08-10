package publishersignupunitofwork

import (


	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readconcern"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	// "go.mongodb.org/mongo-driver/mongo/writeconcern"
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

func (w MongoDBPublisherSignUpUnitOfWork) Save(usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error {

	// wc := writeconcern.Majority()
	// users := w.db.Collection("users")
	// managedUsers := w.db.Collection("managedUsers")

	// txnOptions := options.Transaction().SetWriteConcern(wc)

	// client := w.db.Client()

	// session, err := client.StartSession()
	// if err != nil {
	// 	return err
	// }

	// defer session.EndSession(context.Background())

	// tnxBody := func(ctx mongo.SessionContext) (interface{}, error) {

	// }

	// result, err := session.WithTransaction(context.Background(), tnxBody, txnOptions)

	// err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
	// 	session.StartTransaction(txnOptions)
	// 	//...
	// 	session.CommitTransaction(sc)
	// 	return nil // returned
	// })
	// if err != nil {
	// 	session.AbortTransaction(context.Background())
	// }








	// // TODO this method needs to be transactional
	// err := uow.userRepository.Save(usr)
	// if err != nil {
	// 	return fmt.Errorf("save user failed: %w", err)
	// }
	// err = uow.managedUserRepository.Save(managedUser)
	// if err != nil {
	// 	return fmt.Errorf("save managed user failed : %w", err)
	// }
	// err = uow.publisherRepository.Save(pub)
	// if err != nil {
	// 	return fmt.Errorf("save publisher failed : %w", err)
	// }
	return nil
}

func (uow MongoDBPublisherSignUpUnitOfWork) UserRepository() user.UserRepository {
	return uow.userRepository
}
