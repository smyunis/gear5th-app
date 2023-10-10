package userrepository

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

type MongoDBUserRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
}

func NewMongoDBUserRepository(dbStore mongodbpersistence.MongoDBStore) MongoDBUserRepository {
	
	dbIndexes := dbStore.Database().Collection("users").Indexes()
	dbIndexes.CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"email", 1}},
	})

	return MongoDBUserRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
	}
}

func (r MongoDBUserRepository) Get(ctx context.Context, id shared.ID) (user.User, error) {

	users := r.db.Collection("users")
	resUser := users.FindOne(ctx, bson.M{"_id": id.String()})

	var res bson.M
	err := resUser.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, application.ErrEntityNotFound
		}
		return user.User{}, err
	}

	u := mapMToUser(res)

	return u, nil
}

func (r MongoDBUserRepository) Save(ctx context.Context, u user.User) error {
	id := u.ID.String()
	users := r.db.Collection("users")

	var dbEntry = mapUserToM(u)

	updateOptions := options.Update().SetUpsert(true)
	_, err := users.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)

	if err != nil {
		return err
	}

	return nil
}

func (r MongoDBUserRepository) UserWithEmail(ctx context.Context, email user.Email) (user.User, error) {

	users := r.db.Collection("users")
	resUser := users.FindOne(context.Background(), bson.M{"email": email.String()})

	var res bson.M
	err := resUser.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, application.ErrEntityNotFound
		}
		return user.User{}, err
	}

	u := mapMToUser(res)

	return u, nil
}

func mapMToUser(res primitive.M) user.User {
	userId := shared.ID(res["_id"].(string))
	email, _ := user.NewEmail(res["email"].(string))
	phoneNumber, _ := user.NewPhoneNumber(res["phoneNumber"].(string))
	authNMethod := user.AuthenticationMethod(res["authenticationMethod"].(int32))
	isEmailVerified := res["isEmailVerified"].(bool)

	rolesRes := res["roles"].(primitive.A)
	roles := make([]user.UserRole, 0)
	for _, role := range rolesRes {
		roles = append(roles, user.UserRole(role.(int32)))
	}

	signUpDate := res["signUpDate"].(primitive.DateTime)

	usr := user.ReconstituteUser(userId, email, phoneNumber,
		isEmailVerified, roles, authNMethod, signUpDate.Time())
	return usr
}

func mapUserToM(u user.User) bson.M {
	var dbEntry = bson.M{
		"_id":                  u.ID.String(),
		"email":                u.Email.String(),
		"phoneNumber":          u.PhoneNumber.String(),
		"isEmailVerified":      u.IsEmailVerified,
		"roles":                u.Roles,
		"authenticationMethod": u.AuthenticationMethod,
		"signUpDate":           u.SignUpDate,
	}
	return dbEntry
}
