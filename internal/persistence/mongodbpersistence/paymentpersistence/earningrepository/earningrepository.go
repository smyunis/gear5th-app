package earningrepository

import (
	"context"
	"errors"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBEarningRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBEarningRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBEarningRepository {
	return MongoDBEarningRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBEarningRepository) Get(ctx context.Context, id shared.ID) (earning.Earning, error) {
	earnings := r.db.Collection("earnings")
	sr := earnings.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return earning.Earning{}, application.ErrEntityNotFound
		}
		r.logger.Error("earning/persistence/get", err)
		return earning.Earning{}, err
	}
	s := mapMToEarning(res)
	return s, nil
}

func (r MongoDBEarningRepository) Save(ctx context.Context, e earning.Earning) error {
	id := e.ID.String()
	earnings := r.db.Collection("earnings")
	var dbEntry = mapEarningToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := earnings.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("earning/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBEarningRepository) EarningsForPublisher(publisherID shared.ID, start time.Time, end time.Time) ([]earning.Earning, error) {
	earnings := r.db.Collection("earnings")
	cursor, err := earnings.Find(context.Background(), bson.D{{"publisherId", publisherID.String()},
		{"time", bson.M{"$gt": start}}, {"time", bson.M{"$lt": end}}})
	if err != nil {
		return []earning.Earning{}, err
	}

	pubEarnings := make([]earning.Earning, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []earning.Earning{}, application.ErrEntityNotFound
			}
			r.logger.Error("earnings/persistence/earnings-for-publisher", err)
			return []earning.Earning{}, err
		}
		s := mapMToEarning(res)
		pubEarnings = append(pubEarnings, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("earnings/persistence/earnings-for-publisher", err)
		return []earning.Earning{}, err
	}
	return pubEarnings, nil
}

var x earning.EarningRepository = MongoDBEarningRepository{}

func mapEarningToM(e earning.Earning) bson.M {
	return bson.M{
		"_id":         e.ID.String(),
		"publisherId": e.PublisherID.String(),
		"amount":      e.Amount,
		"reason":      e.Reason,
		"adPieceId":   e.AdPieceID.String(),
		"adSlotId":    e.AdSlotID.String(),
		"siteId":      e.SiteID.String(),
		"time":        e.Time,
	}
}

func mapMToEarning(res primitive.M) earning.Earning {
	return earning.Earning{
		ID:          shared.ID(res["_id"].(string)),
		Events:      make(shared.Events),
		PublisherID: shared.ID(res["publisherId"].(string)),
		AdPieceID:   shared.ID(res["adPieceId"].(string)),
		AdSlotID:    shared.ID(res["adSlotId"].(string)),
		SiteID:      shared.ID(res["siteId"].(string)),
		Reason:      earning.EarningReason(res["reason"].(int32)),
		Amount:      res["amount"].(float64),
		Time:        res["time"].(primitive.DateTime).Time(),
	}
}
