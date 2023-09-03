package oauthuserrepository

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBOAuthUserRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
}



func NewMongoDBOAuthUserRepository(dbStore mongodbpersistence.MongoDBStore) MongoDBOAuthUserRepository {
	return MongoDBOAuthUserRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
	}
}

func (r MongoDBOAuthUserRepository) Get(ctx context.Context, id shared.ID) (user.OAuthUser, error) {
	oauthUsers := r.db.Collection("oauthUsers")
	sr := oauthUsers.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.OAuthUser{}, application.ErrEntityNotFound
		}
		return user.OAuthUser{}, err
	}
	s := mapMToOAuthUser(res)
	return s, nil
}


func (r MongoDBOAuthUserRepository) Save(ctx context.Context, e user.OAuthUser) error {
	id := e.UserID().String()
	oauthUsers := r.db.Collection("oauthUsers")
	var dbEntry = mapOAuthUserToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := oauthUsers.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		return err
	}
	return nil
}


func (r MongoDBOAuthUserRepository) UserWithAccountID(accountID string) (user.OAuthUser, error) {
	users := r.db.Collection("oauthUsers")
	resUser := users.FindOne(context.Background(), bson.M{"userAccountId": accountID})

	var res bson.M
	err := resUser.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.OAuthUser{}, application.ErrEntityNotFound
		}
		return user.OAuthUser{}, err
	}
	u := mapMToOAuthUser(res)
	return u, nil
}


func mapMToOAuthUser(res primitive.M) user.OAuthUser {
	return user.ReconstituteOAuthUser(
		shared.ID(res["_id"].(string)),
		res["userAccountId"].(string),
		user.OAuthProvider(res["oauthProvider"].(int32)))
}

func mapOAuthUserToM(u user.OAuthUser) bson.M {
	return bson.M{
		"_id":           u.UserID().String(),
		"userAccountId": u.UserAccountID(),
		"oauthProvider": u.OAuthProvider(),
	}
}
