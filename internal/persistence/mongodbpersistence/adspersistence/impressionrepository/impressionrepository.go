package impressionrepository

import (
	"context"
	"errors"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/impression"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBImpressionRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBImpressionRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBImpressionRepository {

	// TODO optimise this repo for write efficency
	dbIndexes := dbStore.Database().Collection("impressions").Indexes()
	dbIndexes.DropAll(context.Background())

	return MongoDBImpressionRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBImpressionRepository) Get(ctx context.Context, id shared.ID) (impression.Impression, error) {
	impressions := r.db.Collection("impressions")
	sr := impressions.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return impression.Impression{}, application.ErrEntityNotFound
		}
		r.logger.Error("impression/persistence/get", err)
		return impression.Impression{}, err
	}
	s := mapMToImpression(res)
	return s, nil
}

func (r MongoDBImpressionRepository) Save(ctx context.Context, e impression.Impression) error {
	id := e.ID.String()
	impressions := r.db.Collection("impressions")
	var dbEntry = mapImpressionToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := impressions.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("impression/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBImpressionRepository) ImpressionsForPublisher(publisherID shared.ID, start time.Time, end time.Time) ([]impression.Impression, error) {
	impressions := r.db.Collection("impressions")
	cursor, err := impressions.Find(context.Background(), bson.D{{"publisherId", publisherID.String()},
		{"time", bson.M{"$gt": start}}, {"time", bson.M{"$lt": end}}})
	if err != nil {
		return []impression.Impression{}, err
	}

	pubImpressions := make([]impression.Impression, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []impression.Impression{}, application.ErrEntityNotFound
			}
			r.logger.Error("impression/persistence/impressions-for-publisher", err)
			return []impression.Impression{}, err
		}
		s := mapMToImpression(res)
		pubImpressions = append(pubImpressions, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("impression/persistence/impressions-for-publisher", err)
		return []impression.Impression{}, err
	}
	return pubImpressions, nil
}

func (r MongoDBImpressionRepository) DailyImpressionCount(day time.Time) (int, error) {
	impressions := r.db.Collection("impressions")

	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	count, err := impressions.CountDocuments(context.Background(),
		bson.D{{"time", bson.M{"$gt": dayStart}}, {"time", bson.M{"$lt": dayEnd}}})
	if err != nil {
		r.logger.Error("impression/persistence/count-daily-impression", err)
		return 0, err
	}

	return int(count), nil
}

func (r MongoDBImpressionRepository) ImpressionsCountForPublisher(publisherID shared.ID, start time.Time, end time.Time) (int, error) {
	impressions := r.db.Collection("impressions")
	count, err := impressions.CountDocuments(context.Background(),
		bson.D{{"publisherId", publisherID.String()},
			{"time", bson.M{"$gt": start}}, {"time", bson.M{"$lt": end}}})
	if err != nil {
		r.logger.Error("impression/persistence/count-impressions-for-publisher", err)
		return 0, err
	}
	return int(count), nil
}

func mapImpressionToM(s impression.Impression) bson.M {
	return bson.M{
		"_id":         s.ID.String(),
		"slotId":      s.OriginAdSlotID.String(),
		"adPieceId":   s.AdPieceID.String(),
		"siteId":      s.OriginSiteID.String(),
		"publisherId": s.OriginPublisherID.String(),
		"time":        s.Time,
	}
}

func mapMToImpression(res primitive.M) impression.Impression {
	return impression.Impression{
		ID:                shared.ID(res["_id"].(string)),
		Events:            make(shared.Events),
		AdPieceID:         shared.ID(res["adPieceId"].(string)),
		OriginSiteID:      shared.ID(res["siteId"].(string)),
		OriginAdSlotID:    shared.ID(res["slotId"].(string)),
		OriginPublisherID: shared.ID(res["publisherId"].(string)),
		Time:              res["time"].(primitive.DateTime).Time(),
	}
}
