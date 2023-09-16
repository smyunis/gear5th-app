package userrepository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
)

type userCache struct {
	ID                   string
	Email                string
	PhoneNumber          string
	IsEmailVerified      bool
	Roles                []user.UserRole
	AuthenticationMethod user.AuthenticationMethod
	SignUpDate           time.Time
}

type MongoDBUserRepositoryCached struct {
	delegate   user.UserRepository
	cacheStore application.KeyValueStore
}

func NewMongoDBUserRepositoryCached(dbStore mongodbpersistence.MongoDBStore,
	cacheStore application.KeyValueStore) MongoDBUserRepositoryCached {
	return MongoDBUserRepositoryCached{
		delegate: MongoDBUserRepository{
			dbStore: dbStore,
			db:      dbStore.Database(),
		},
		cacheStore: cacheStore,
	}
}

func (r MongoDBUserRepositoryCached) Get(ctx context.Context, id shared.ID) (user.User, error) {
	cachedUser := &userCache{}
	uStr, err := r.cacheStore.Get(userWithIDCacheKey(id))
	jsonErr := json.Unmarshal([]byte(uStr), cachedUser)
	if err != nil || jsonErr != nil {
		u, err := r.delegate.Get(ctx, id)
		if err != nil {
			return user.User{}, err
		}
		j, err := json.Marshal(userCache{
			u.ID.String(),
			u.Email.String(),
			u.PhoneNumber.String(),
			u.IsEmailVerified,
			u.Roles,
			u.AuthenticationMethod,
			u.SignUpDate,
		})
		if err != nil {
			return u, nil
		}
		r.cacheStore.Save(userWithIDCacheKey(id), string(j), 3*time.Hour)
		return u, nil
	}
	email, _ := user.NewEmail(cachedUser.Email)
	pn, _ := user.NewPhoneNumber(cachedUser.PhoneNumber)
	return user.User{
		shared.ID(cachedUser.ID),
		email,
		pn,
		cachedUser.IsEmailVerified,
		cachedUser.Roles,
		cachedUser.AuthenticationMethod,
		cachedUser.SignUpDate,
		make(shared.Events),
	}, nil
}

func (r MongoDBUserRepositoryCached) UserWithEmail(ctx context.Context, email user.Email) (user.User, error) {
	return r.delegate.UserWithEmail(ctx, email)
}

func (r MongoDBUserRepositoryCached) Save(ctx context.Context, e user.User) error {
	return r.delegate.Save(ctx, e)
}

func userWithIDCacheKey(id shared.ID) string {
	return fmt.Sprintf("user:%s", id.String())
}

func userWithEmailCacheKey(email user.Email) string {
	return fmt.Sprintf("user:%s", email.String())
}

var x user.UserRepository = MongoDBUserRepositoryCached{}
