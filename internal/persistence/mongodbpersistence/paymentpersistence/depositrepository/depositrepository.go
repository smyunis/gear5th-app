package depositrepository

import (
	"context"
	"errors"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/deposit"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBDepositRepository struct {
	dbStore mongodbpersistence.MongoDBStore
	db      *mongo.Database
	logger  application.Logger
}

func NewMongoDBDepositRepository(dbStore mongodbpersistence.MongoDBStore,
	logger application.Logger) MongoDBDepositRepository {
	return MongoDBDepositRepository{
		dbStore: dbStore,
		db:      dbStore.Database(),
		logger:  logger,
	}
}

func (r MongoDBDepositRepository) Get(ctx context.Context, id shared.ID) (deposit.Deposit, error) {
	deposits := r.db.Collection("deposits")
	sr := deposits.FindOne(ctx, bson.D{{"_id", id.String()}})
	var res bson.M
	err := sr.Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return deposit.Deposit{}, application.ErrEntityNotFound
		}
		r.logger.Error("deposit/persistence/get", err)
		return deposit.Deposit{}, err
	}
	s := mapMToDeposit(res)
	return s, nil
}

func (r MongoDBDepositRepository) Save(ctx context.Context, e deposit.Deposit) error {
	id := e.ID.String()
	deposits := r.db.Collection("deposits")
	var dbEntry = mapDepositToM(e)
	updateOptions := options.Update().SetUpsert(true)
	_, err := deposits.UpdateByID(ctx, id, bson.D{{"$set", dbEntry}}, updateOptions)
	if err != nil {
		r.logger.Error("deposit/persistence/save", err)
		return err
	}
	return nil
}

func (r MongoDBDepositRepository) DailyDisposits(day time.Time) ([]deposit.Deposit, error) {
	deposits := r.db.Collection("deposits")
	cursor, err := deposits.Find(context.Background(), bson.D{
		{"start", bson.M{"$lt": day}}, {"end", bson.M{"$gt": day}}})
	if err != nil {
		return []deposit.Deposit{}, err
	}

	todayDeposits := make([]deposit.Deposit, 0)
	for cursor.Next(context.Background()) {
		var res bson.M
		err := cursor.Decode(&res)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return []deposit.Deposit{}, application.ErrEntityNotFound
			}
			r.logger.Error("deposits/persistence/deposits-available-for-day", err)
			return []deposit.Deposit{}, err
		}
		s := mapMToDeposit(res)
		todayDeposits = append(todayDeposits, s)
	}
	if err := cursor.Err(); err != nil {
		r.logger.Error("deposits/persistence/deposits-available-for-day", err)
		return []deposit.Deposit{}, err
	}
	return todayDeposits, nil
}

func mapDepositToM(s deposit.Deposit) bson.M {
	return bson.M{
		"_id":          s.ID.String(),
		"advertiserId": s.AdvertiserID.String(),
		"amount":       s.Amount,
		"depositTime":  s.DepositTime,
		"start":        s.Start,
		"end":          s.End,
	}
}

func mapMToDeposit(res primitive.M) deposit.Deposit {
	return deposit.Deposit{
		ID:           shared.ID(res["_id"].(string)),
		Events:       make(shared.Events),
		AdvertiserID: shared.ID(res["advertiserId"].(string)),
		Amount:       res["amount"].(float64),
		DepositTime:  res["depositTime"].(primitive.DateTime).Time(),
		Start:        res["start"].(primitive.DateTime).Time(),
		End:          res["end"].(primitive.DateTime).Time(),
	}
}

