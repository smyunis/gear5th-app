package advertiserrepository

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/advertiser"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBAdvertiserRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBAdvertiserRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBAdvertiserRepository {
	return MongoDBAdvertiserRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBAdvertiserRepository) Get(ctx context.Context, id shared.ID) (advertiser.Advertiser, error) {
	advertisers := r.db.Collection("advertisers")
	sr := advertisers.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return advertiser.Advertiser{}, application.ErrEntityNotFound
		}
		r.logger.Error("advertisers/persistence/get", err)
		return advertiser.Advertiser{}, err
	}
	s := mapMToAdvertiser(res)
	return s, nil
}

func (r MongoDBAdvertiserRepository) Save(ctx context.Context, e advertiser.Advertiser) error {
	id := e.UserID.String()
	advertisers := r.db.Collection("advertisers")
	var dbEntry = mapAdvertiserToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := advertisers.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("advertiser/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBAdvertiserRepository) Advertisers() ([]advertiser.Advertiser, error) {
	advertisers := r.db.Collection("advertisers")
	cursor, err := advertisers.Find(context.Background(), bson.D{}, &options.FindOptions{
		Sort: bson.D{{"name", 1}},
	})
	if err != nil {
		return []advertiser.Advertiser{}, err
	}

	allAdvertisers := make([]advertiser.Advertiser, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []advertiser.Advertiser{}, application.ErrEntityNotFound
			}
			r.logger.Error("advertisers/persistence/all", err)
			return []advertiser.Advertiser{}, err
		}
		s := mapMToAdvertiser(res)
		allAdvertisers = append(allAdvertisers, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("advertisers/persistence/all", err)
		return []advertiser.Advertiser{}, err
	}
	return allAdvertisers, nil
}

func mapAdvertiserToM(s advertiser.Advertiser) bson.M {
	return bson.M{
		"_id":  s.UserID.String(),
		"name": s.Name,
		"note": s.Note,
	}
}

func mapMToAdvertiser(res primitive.M) advertiser.Advertiser {
	return advertiser.Advertiser{
		UserID: shared.ID(res["_id"].(string)),
		Name:   res["name"].(string),
		Note:   res["note"].(string),
	}
}
