package identitypersistence

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
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
	return MongoDBUserRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
	}
}

func (r MongoDBUserRepository) Get(id shared.ID) (user.User, error) {

	u := &user.User{}
	users := r.db.Collection("users")
	resUser := users.FindOne(context.Background(), bson.M{"_id": id.String()})

	var res bson.M
	err := resUser.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return *u, shared.NewEntityNotFoundError(id.String(), "user")
		}
		return *u, err
	}

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

	return usr, nil
}

func (r MongoDBUserRepository) Save(u user.User) error {
	id := u.UserID().String()
	users := r.db.Collection("users")

	var dbEntry = bson.M{
		"_id":                  id,
		"email":                u.Email().Email(),
		"phoneNumber":          u.PhoneNumber().PhoneNumber(),
		"isEmailVerified":      u.IsEmailVerified(),
		"roles":                u.Roles(),
		"authenticationMethod": u.AuthenticationMethod(),
		"signUpDate":           u.SignUpDate(),
	}

	updateOptions := options.Update().SetUpsert(true)
	_, err := users.UpdateByID(context.Background(), id, bson.D{{"$set", dbEntry}}, updateOptions)

	if err != nil {
		return err
	}

	return nil
}

func (r MongoDBUserRepository) UserWithEmail(email user.Email) (user.User, error) {

	return user.User{}, shared.NewEntityNotFoundError(email.Email(), "user")
}
