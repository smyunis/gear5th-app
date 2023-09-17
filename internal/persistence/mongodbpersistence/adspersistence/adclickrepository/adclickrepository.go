package adclickrepository

import (
	"context"
	"errors"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/adclick"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO optimise this repo for write efficency
type MongoDBAdClickRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBAdClickRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBAdClickRepository {
	return MongoDBAdClickRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBAdClickRepository) Get(ctx context.Context, id shared.ID) (adclick.AdClick, error) {
	adClicks := r.db.Collection("adClicks")
	sr := adClicks.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return adclick.AdClick{}, application.ErrEntityNotFound
		}
		r.logger.Error("adclick/persistence/get", err)
		return adclick.AdClick{}, err
	}
	s := mapMToAdClick(res)
	return s, nil
}

func (r MongoDBAdClickRepository) Save(ctx context.Context, e adclick.AdClick) error {
	id := e.ID.String()
	adClicks := r.db.Collection("adClicks")
	var dbEntry = mapAdClickToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := adClicks.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("adclick/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBAdClickRepository) AdClicksForPublisher(publisherID shared.ID, start time.Time, end time.Time) ([]adclick.AdClick, error) {
	adClicks := r.db.Collection("adClicks")
	cursor, err := adClicks.Find(context.Background(), bson.D{{"publisherId", publisherID.String()},
		{"time", bson.M{"$gt": start}}, {"time", bson.M{"$lt": end}}})
	if err != nil {
		return []adclick.AdClick{}, err
	}

	pubAdClicks := make([]adclick.AdClick, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []adclick.AdClick{}, application.ErrEntityNotFound
			}
			r.logger.Error("adclick/persistence/adclicks-for-publisher", err)
			return []adclick.AdClick{}, err
		}
		s := mapMToAdClick(res)
		pubAdClicks = append(pubAdClicks, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("adclick/persistence/adclicks-for-publisher", err)
		return []adclick.AdClick{}, err
	}
	return pubAdClicks, nil
}

func mapAdClickToM(s adclick.AdClick) bson.M {
	return bson.M{
		"_id":         s.ID.String(),
		"viewId":      s.ViewID.String(),
		"slotId":      s.OriginAdSlotID.String(),
		"adPieceId":   s.ClickedAdPieceID.String(),
		"siteId":      s.OriginSiteID.String(),
		"publisherId": s.OriginPublisherID.String(),
		"time":        s.Time,
	}
}

func mapMToAdClick(res primitive.M) adclick.AdClick {
	return adclick.AdClick{
		ID:                shared.ID(res["_id"].(string)),
		Events:            make(shared.Events),
		ViewID:            shared.ID(res["viewId"].(string)),
		ClickedAdPieceID:  shared.ID(res["adPieceId"].(string)),
		OriginSiteID:      shared.ID(res["siteId"].(string)),
		OriginAdSlotID:    shared.ID(res["slotId"].(string)),
		OriginPublisherID: shared.ID(res["publisherId"].(string)),
		Time:              res["time"].(primitive.DateTime).Time(),
	}
}
