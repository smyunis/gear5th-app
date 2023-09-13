package campaignrepository

import (
	"context"
	"errors"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBCampaignRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBCampaignRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBCampaignRepository {
	return MongoDBCampaignRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBCampaignRepository) Get(ctx context.Context, id shared.ID) (campaign.Campaign, error) {
	campaigns := r.db.Collection("campaigns")
	sr := campaigns.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return campaign.Campaign{}, application.ErrEntityNotFound
		}
		r.logger.Error("campaign/persistence/get", err)
		return campaign.Campaign{}, err
	}
	s := mapMToCampaign(res)
	return s, nil
}

func (r MongoDBCampaignRepository) Save(ctx context.Context, e campaign.Campaign) error {
	id := e.ID.String()
	campaigns := r.db.Collection("campaigns")
	var dbEntry = mapCampaignToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := campaigns.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("campaign/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBCampaignRepository) CampaignsForAdvertiser(advertiserID shared.ID) ([]campaign.Campaign, error) {
	campaigns := r.db.Collection("campaigns")
	cursor, err := campaigns.Find(context.Background(), bson.D{{"advertiserUserId", advertiserID.String()}, {"isPreempted", false}})
	if err != nil {
		return []campaign.Campaign{}, err
	}

	activeCampaigns := make([]campaign.Campaign, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []campaign.Campaign{}, application.ErrEntityNotFound
			}
			r.logger.Error("campaign/persistence/activecampaigns", err)
			return []campaign.Campaign{}, err
		}
		s := mapMToCampaign(res)
		activeCampaigns = append(activeCampaigns, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("campaign/persistence/activecampaigns", err)
		return []campaign.Campaign{}, err
	}
	return activeCampaigns, nil
}

func (r MongoDBCampaignRepository) RunningCampaigns() ([]campaign.Campaign, error) {
	now := time.Now()
	campaigns := r.db.Collection("campaigns")
	cursor, err := campaigns.Find(context.Background(),
		bson.D{{"start", bson.M{"$lt": now}}, {"end", bson.M{"$gt": now}}, {"isPreempted", false}})
	if err != nil {
		return []campaign.Campaign{}, err
	}

	runningCampaings := make([]campaign.Campaign, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []campaign.Campaign{}, application.ErrEntityNotFound
			}
			r.logger.Error("campaign/persistence/runningcampaigns", err)
			return []campaign.Campaign{}, err
		}
		s := mapMToCampaign(res)
		runningCampaings = append(runningCampaings, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("campaign/persistence/runningcampaigns", err)
		return []campaign.Campaign{}, err
	}
	return runningCampaings, nil
}

func mapCampaignToM(s campaign.Campaign) bson.M {
	return bson.M{
		"_id":              s.ID.String(),
		"name":             s.Name,
		"advertiserUserId": s.AdvertiserUserID.String(),
		"start":            s.Start,
		"end":              s.End,
		"runDuration":      s.RunDuration,
		"isPreempted":      s.IsPreempted,
	}
}

func mapMToCampaign(res primitive.M) campaign.Campaign {
	return campaign.Campaign{
		ID:               shared.ID(res["_id"].(string)),
		Events:           make(shared.Events),
		Name:             res["name"].(string),
		AdvertiserUserID: shared.ID(res["advertiserUserId"].(string)),
		Start:            res["start"].(primitive.DateTime).Time(),
		End:              res["end"].(primitive.DateTime).Time(),
		RunDuration:      time.Duration(res["runDuration"].(int64)),
		IsPreempted:      res["isPreempted"].(bool),
	}

}
