package publisherrepository

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBPublisherRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
}

func NewMongoDBPublisherRepository(dbStore mongodbpersistence.MongoDBStore) MongoDBPublisherRepository {
	return MongoDBPublisherRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
	}
}

func (r MongoDBPublisherRepository) Get(ctx context.Context, id shared.ID) (publisher.Publisher, error) {

	publishers := r.db.Collection("publishers")

	sr := publishers.FindOne(ctx, bson.D{{"_id", id.String()}})

	var res bson.M
	err := sr.Decode(&res)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return publisher.Publisher{}, application.ErrEntityNotFound
		}
		return publisher.Publisher{}, err
	}

	p := mapMToPublisher(res)

	return p, nil
}

func (r MongoDBPublisherRepository) Save(ctx context.Context, pub publisher.Publisher) error {
	publishers := r.db.Collection("publishers")

	p := mapPublisherToM(pub)

	updateOptions := options.Update().SetUpsert(true)
	_, err := publishers.UpdateByID(ctx, pub.UserID.String(),
		bson.D{{"$set", p}}, updateOptions)

	if err != nil {
		return err
	}

	return nil
}

func mapPublisherToM(pub publisher.Publisher) primitive.M {

	p := bson.M{
		"_id":              pub.UserID.String(),
		"lastDisbursement": pub.LastDisbursement,
	}
	return p
}

func mapMToPublisher(res primitive.M) publisher.Publisher {

	p := publisher.Publisher{
		shared.ID(res["_id"].(string)),
		res["lastDisbursement"].(primitive.DateTime).Time(),
	}
	return p
}
