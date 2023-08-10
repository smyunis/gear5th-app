package publisherrepository

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
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

func (r MongoDBPublisherRepository) Get(ctx context.Context,id shared.ID) (publisher.Publisher, error) {

	publishers := r.db.Collection("publishers")

	sr := publishers.FindOne(context.Background(), bson.D{{"_id", id.String()}})

	var res bson.M
	err := sr.Decode(&res)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return publisher.Publisher{}, shared.NewEntityNotFoundError(id.String(), "publishers")
		}
		return publisher.Publisher{}, err
	}

	p := mapMToPublisher(res)

	return p, nil
}

func (r MongoDBPublisherRepository) Save(ctx context.Context,pub publisher.Publisher) error {
	publishers := r.db.Collection("publishers")

	p := mapPublisherToM(pub)

	updateOptions := options.Update().SetUpsert(true)
	_, err := publishers.UpdateByID(context.Background(), pub.UserID().String(),
		bson.D{{"$set", p}}, updateOptions)

	if err != nil {
		return err
	}

	return nil
}

func mapPublisherToM(pub publisher.Publisher) primitive.M {
	notifications := bson.A{}
	for _, n := range pub.UnacknowledgedNotifications() {
		notifications = append(notifications, bson.M{
			"message": n.Message(),
			"time":    n.Time(),
		})
	}

	p := bson.M{
		"_id":                         pub.UserID().String(),
		"unacknowledgedNotifications": notifications,
	}
	return p
}

func mapMToPublisher(res primitive.M) publisher.Publisher {
	nr := res["unacknowledgedNotifications"].(primitive.A)
	notifications := make([]shared.Notification, 0)
	for _, n := range nr {
		nm := n.(bson.M)
		nmessage := nm["message"].(string)
		ntime := nm["time"].(primitive.DateTime).Time()
		notifications = append(notifications,
			shared.ReconstituteNotification(nmessage, ntime))
	}

	p := publisher.ReconstitutePublisher(
		shared.ID(res["_id"].(string)),
		notifications)
	return p
}
