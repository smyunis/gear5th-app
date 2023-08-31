package adslotrepository

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBAdSlotRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBAdSlotRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBAdSlotRepository {
	return MongoDBAdSlotRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBAdSlotRepository) Get(ctx context.Context, id shared.ID) (adslot.AdSlot, error) {
	adSlots := r.db.Collection("adSlots")
	sr := adSlots.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return adslot.AdSlot{}, application.ErrEntityNotFound
		}
		r.logger.Error("adslotrepository/get", err)
		return adslot.AdSlot{}, err
	}
	s := mapMToAdSlot(res)
	return s, nil
}

func (r MongoDBAdSlotRepository) Save(ctx context.Context, e adslot.AdSlot) error {
	id := e.ID().String()
	adSlots := r.db.Collection("adSlots")
	var dbEntry = mapAdSlotToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := adSlots.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("adslotrepository/save", err)
		return err
	}
	return nil
}

func (r MongoDBAdSlotRepository) ActiveAdSlotsForSite(siteID shared.ID) ([]adslot.AdSlot, error) {
	adSlots := r.db.Collection("adSlots")
	cursor, err := adSlots.Find(context.Background(),
		bson.D{{"siteId", siteID.String()}, {"isDeactivated", false}})
	if err != nil {
		return []adslot.AdSlot{}, err
	}

	activeAdSlots := make([]adslot.AdSlot, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []adslot.AdSlot{}, application.ErrEntityNotFound
			}
			r.logger.Error("adslot/persistence/activeadslotsforsite", err)
			return []adslot.AdSlot{}, err
		}
		s := mapMToAdSlot(res)
		activeAdSlots = append(activeAdSlots, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("adslot/persistence/activeadslotsforsite", err)
		return []adslot.AdSlot{}, err
	}
	return activeAdSlots, nil
}

func mapMToAdSlot(res primitive.M) adslot.AdSlot {
	return adslot.ReconstituteAdSlot(
		shared.ID(res["_id"].(string)),
		shared.ID(res["siteId"].(string)),
		res["name"].(string),
		adslot.AdSlotType(res["adSlotType"].(int32)),
		res["isDeactivated"].(bool))
}

func mapAdSlotToM(s adslot.AdSlot) bson.M {
	return bson.M{
		"_id":           s.ID().String(),
		"siteId":        s.SiteID().String(),
		"name":          s.Name(),
		"adSlotType":    s.Type(),
		"isDeactivated": s.IsDeactivated(),
	}
}
