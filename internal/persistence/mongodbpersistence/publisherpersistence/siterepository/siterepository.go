package siterepository

import (
	"context"
	"errors"
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBSiteRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBSiteRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBSiteRepository {
	return MongoDBSiteRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBSiteRepository) Get(ctx context.Context, id shared.ID) (site.Site, error) {
	sites := r.db.Collection("sites")
	sr := sites.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return site.Site{}, application.ErrEntityNotFound
		}
		return site.Site{}, err
	}
	s := mapMToSite(res)
	return s, nil
}

func (r MongoDBSiteRepository) Save(ctx context.Context, e site.Site) error {
	id := e.ID.String()
	sites := r.db.Collection("sites")

	var dbEntry = mapSiteToM(e)

	updateOptions := options.Update().SetUpsert(true)
	_, err := sites.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		return err
	}
	return nil
}

func (r MongoDBSiteRepository) ActiveSitesForPublisher(publisherID shared.ID) ([]site.Site, error) {
	sites := r.db.Collection("sites")
	cursor, err := sites.Find(context.Background(), bson.D{{"publisherId", publisherID.String()}, {"isDeactivated", false}})
	if err != nil {
		return []site.Site{}, err
	}

	activeSites := make([]site.Site, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []site.Site{}, application.ErrEntityNotFound
			}
			r.logger.Error("site/persistence/activesites", err)
			return []site.Site{}, err
		}
		s := mapMToSite(res)
		activeSites = append(activeSites, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("site/persistence/activesites", err)
		return []site.Site{}, err
	}
	return activeSites, nil
}

func mapSiteToM(s site.Site) bson.M {

	monetizationHistory := bson.A{}
	for _, history := range s.MonetizationStatusHistory() {
		monetizationHistory = append(monetizationHistory, bson.M{
			"canMonetize": history.CanMonetize,
			"time":        history.Time,
		})
	}

	siteURL := s.URL
	var dbEntry = bson.M{
		"_id":                 s.ID.String(),
		"url":                 siteURL.String(),
		"isVerified":          s.IsVerified,
		"publisherId":         s.PublisherID.String(),
		"isDeactivated":       s.IsDeactivated,
		"monetizationHistory": monetizationHistory,
	}
	return dbEntry
}

func mapMToSite(res primitive.M) site.Site {
	mh := res["monetizationHistory"].(primitive.A)
	monetiaztionStatusHistory := make([]site.MonetizationStatus, 0)
	for _, historyEntry := range mh {
		h := historyEntry.(bson.M)

		canMonetize := h["canMonetize"].(bool)
		atTime := h["time"].(primitive.DateTime).Time()

		monetiaztionStatusHistory = append(monetiaztionStatusHistory,
			site.MonetizationStatus{
				CanMonetize: canMonetize,
				Time:        atTime,
			})
	}
	siteURL, _ := url.Parse(res["url"].(string))

	s := site.ReconstituteSite(
		shared.ID(res["_id"].(string)),
		shared.ID(res["publisherId"].(string)),
		*siteURL,
		res["isVerified"].(bool),
		res["isDeactivated"].(bool), monetiaztionStatusHistory)
	return s
}
