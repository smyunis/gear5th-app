package manageduserrepository

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBMangageUserRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
}

func NewMongoDBMangageUserRepository(dbStore mongodbpersistence.MongoDBStore) MongoDBMangageUserRepository {
	return MongoDBMangageUserRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
	}
}


func (r MongoDBMangageUserRepository) Get(id shared.ID) (user.ManagedUser, error) {

	managedUsers := r.db.Collection("managedUsers")

	sr := managedUsers.FindOne(context.Background(), bson.D{{"_id", id.String()}})

	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.ManagedUser{}, shared.NewEntityNotFoundError(id.String(), "managed user")
		}
		return user.ManagedUser{}, err
	}

	mu := user.ReconstituteManagedUser(
		shared.ID(res["_id"].(string)),
		user.NewPersonNameWithFullName(res["name"].(string)),
		res["hashedPassword"].(string))

	return mu, nil

}

func (r MongoDBMangageUserRepository) Save(u user.ManagedUser) error {
	managedUsers := r.db.Collection("managedUsers")

	mu := bson.M{
		"_id":            u.UserID().String(),
		"name":           u.Name().FullName(),
		"hashedPassword": u.HashedPassword(),
	}

	updateOptions := options.Update().SetUpsert(true)
	_, err := managedUsers.UpdateByID(context.Background(), u.UserID().String(), bson.D{{"$set", mu}}, updateOptions)

	if err != nil {
		return err
	}

	return nil
}