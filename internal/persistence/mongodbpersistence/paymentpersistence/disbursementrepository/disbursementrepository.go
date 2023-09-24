package disbursementrepository

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBDisbursementRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBDisbursementRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBDisbursementRepository {
	return MongoDBDisbursementRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBDisbursementRepository) Get(ctx context.Context, id shared.ID) (disbursement.Disbursement, error) {
	disbursements := r.db.Collection("disbursements")
	sr := disbursements.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return disbursement.Disbursement{}, application.ErrEntityNotFound
		}
		r.logger.Error("disbursements/persistence/get", err)
		return disbursement.Disbursement{}, err
	}
	s := mapMToDisbursement(res)
	return s, nil
}

func (r MongoDBDisbursementRepository) Save(ctx context.Context, e disbursement.Disbursement) error {
	id := e.ID.String()
	disbursements := r.db.Collection("disbursements")
	var dbEntry = mapDisbursementToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := disbursements.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("disbursements/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBDisbursementRepository) DisbursementsForPublisher(publisherID shared.ID, status disbursement.DisbursementStatus) ([]disbursement.Disbursement, error) {
	disbursements := r.db.Collection("disbursements")
	cursor, err := disbursements.Find(context.Background(), bson.D{
		{"publisherId", publisherID.String()}, {"status", status}})
	if err != nil {
		return []disbursement.Disbursement{}, err
	}

	publisherDisbursements := make([]disbursement.Disbursement, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []disbursement.Disbursement{}, application.ErrEntityNotFound
			}
			r.logger.Error("disbursements/persistence/publisher-disbursements", err)
			return []disbursement.Disbursement{}, err
		}
		s := mapMToDisbursement(res)
		publisherDisbursements = append(publisherDisbursements, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("disbursements/persistence/publisher-disbursements", err)
		return []disbursement.Disbursement{}, err
	}
	return publisherDisbursements, nil
}

func mapDisbursementToM(s disbursement.Disbursement) bson.M {
	return bson.M{
		"_id":              s.ID.String(),
		"publisherId":      s.PublisherID.String(),
		"status":           s.Status,
		"amount":           s.Amount,
		"time":             s.Time,
		"start":            s.PeriodStart,
		"end":              s.PeriodEnd,
		"settlementRemark": s.SettlementRemark,
		"paymentProfile": bson.M{
			"paymentMethod": s.PaymentProfile.PaymentMethod,
			"account":       s.PaymentProfile.Account,
			"fullname":      s.PaymentProfile.FullName,
			"phoneNumber":   s.PaymentProfile.PhoneNumber.String(),
		},
	}
}

func mapMToDisbursement(res primitive.M) disbursement.Disbursement {

	profileM := res["paymentProfile"].(primitive.M)
	ph, _ := user.NewPhoneNumber(profileM["phoneNumber"].(string))
	profile := disbursement.PaymentProfile{
		PaymentMethod: profileM["paymentMethod"].(string),
		Account:       profileM["account"].(string),
		FullName:      profileM["fullname"].(string),
		PhoneNumber:   ph,
	}

	return disbursement.Disbursement{
		ID:               shared.ID(res["_id"].(string)),
		PublisherID:      shared.ID(res["publisherId"].(string)),
		Events:           make(shared.Events),
		Amount:           res["amount"].(float64),
		Status:           disbursement.DisbursementStatus(res["status"].(int32)),
		Time:             res["time"].(primitive.DateTime).Time(),
		PeriodStart:      res["start"].(primitive.DateTime).Time(),
		PeriodEnd:        res["end"].(primitive.DateTime).Time(),
		SettlementRemark: res["settlementRemark"].(string),
		PaymentProfile:   profile,
	}
}
