//go:build db
package mongodbpersistence_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)




func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

var client *mongo.Client

func setup() {
	testConfig := mongotestdoubles.NewTestEnvConfigurationProvider()
	clientOptions := options.Client().ApplyURI(testConfig.Get("MONGODB_URL", ""))
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		os.Exit(1)
	}
}



func teardown() {
	client.Disconnect(context.TODO())
}

func TestConnectToDb(t *testing.T) {

}

// func TestFindInDatabase(t *testing.T) {
// 	db := client.Database("dbz")

// 	pirate := db.Collection("pirate")

// 	luffyResult := pirate.FindOne(context.TODO(), bson.D{bson.E{
// 		Key:   "_id",
// 		Value: 1,
// 	}})

// 	var luffy bson.M

// 	err := luffyResult.Decode(&luffy)
// 	if errors.Is(err, mongo.ErrNoDocuments) {
// 		t.Fatal("pirate not found")
// 	}

// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}


// }

func TestInsertDoc(t *testing.T) {
	db := client.Database("dbz")
	pirate := db.Collection("pirate")

	robin := bson.M{
		"_id":    5,
		"name":   "Chopper",
		"bounty": 100,
	}

	updateOptions := options.Update().SetUpsert(true)

	_, err := pirate.UpdateByID(context.TODO(), 5, bson.D{{"$set", robin}}, updateOptions)

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFindInserted(t *testing.T) {
	db := client.Database("dbz")
	pirate := db.Collection("pirate")

	robin := bson.M{
		"_id":    5,
		"name":   "Chopper",
		"bounty": 100,
	}

	updateOptions := options.Update().SetUpsert(true)

	_, err := pirate.UpdateByID(context.TODO(), 5, bson.D{{"$set", robin}}, updateOptions)

	if err != nil {
		t.Fatal(err.Error())
	}

	//find

	luffyResult := pirate.FindOne(context.TODO(), bson.D{bson.E{
		Key:   "_id",
		Value: 5,
	}})

	var luffy bson.M

	err = luffyResult.Decode(&luffy)
	if errors.Is(err, mongo.ErrNoDocuments) {
		t.Fatal("pirate not found")
	}

	if err != nil {
		t.Fatal(err.Error())
	}

}

func TestInsertTotalyNew(t *testing.T) {
	db := client.Database("dbz")
	zwar := db.Collection("zwarriors")

	goku := bson.M{
		"_id":       1,
		"name":      "Son Goku",
		"race":      "sayan",
		"trainedBy": []string{"roshi", "korin", "kami", "king kai", "whis"},
	}

	updateOptions := options.Update().SetUpsert(true)

	_, err := zwar.UpdateByID(context.TODO(), 1, bson.D{{"$set", goku}}, updateOptions)

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGetAListOfItems(t *testing.T) {
	db := client.Database("dbz")
	zwar := db.Collection("zwarriors")

	var res bson.M

	sr := zwar.FindOne(context.TODO(), bson.D{{"_id", 1}})
	sr.Decode(&res)

	gokuTrainers := res["trainedBy"].(primitive.A)

	t.Log(len(gokuTrainers))

}
