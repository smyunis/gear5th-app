package adpiecerepository

import (
	"context"
	"errors"
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBAdPieceRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBAdPieceRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBAdPieceRepository {
	return MongoDBAdPieceRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBAdPieceRepository) Get(ctx context.Context, id shared.ID) (adpiece.AdPiece, error) {
	adPieces := r.db.Collection("adPieces")
	sr := adPieces.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return adpiece.AdPiece{}, application.ErrEntityNotFound
		}
		r.logger.Error("adpiece/persistence/get", err)
		return adpiece.AdPiece{}, err
	}
	s := mapMToAdPiece(res)
	return s, nil
}

func (r MongoDBAdPieceRepository) Save(ctx context.Context, e adpiece.AdPiece) error {
	id := e.ID.String()
	adPieces := r.db.Collection("adPieces")
	var dbEntry = mapAdPieceToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := adPieces.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("adpiece/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBAdPieceRepository) ActiveAdPiecesForCampaign(campaignID shared.ID) ([]adpiece.AdPiece, error) {
	adPieces := r.db.Collection("adPieces")
	cursor, err := adPieces.Find(context.Background(), bson.D{{"campaignId", campaignID.String()}, {"isDeactivated", false}})
	if err != nil {
		return []adpiece.AdPiece{}, err
	}

	activeAdPieces := make([]adpiece.AdPiece, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []adpiece.AdPiece{}, application.ErrEntityNotFound
			}
			r.logger.Error("adpiece/persistence/activecampaigns", err)
			return []adpiece.AdPiece{}, err
		}
		s := mapMToAdPiece(res)
		activeAdPieces = append(activeAdPieces, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("adpiece/persistence/activecampaigns", err)
		return []adpiece.AdPiece{}, err
	}
	return activeAdPieces, nil
}

func mapAdPieceToM(s adpiece.AdPiece) bson.M {
	return bson.M{
		"_id":           s.ID.String(),
		"campaignId":    s.CampaignID.String(),
		"slotType":      s.SlotType,
		"resource":      s.Resource,
		"resourceMIME":  s.ResourceMIMEType,
		"ref":           s.Ref.String(),
		"isDeactivated": s.IsDeactivated,
	}
}

func mapMToAdPiece(res primitive.M) adpiece.AdPiece {
	reflink, err := url.Parse(res["ref"].(string))
	if err != nil {
		reflink = nil
	}
	return adpiece.AdPiece{
		ID:               shared.ID(res["_id"].(string)),
		Events:           make(shared.Events),
		CampaignID:       shared.ID(res["campaignId"].(string)),
		SlotType:         adslot.AdSlotType(res["slotType"].(int32)),
		Resource:         res["resource"].(string),
		ResourceMIMEType: res["resourceMIME"].(string),
		Ref:              reflink,
		IsDeactivated:    res["isDeactivated"].(bool),
	}
}
