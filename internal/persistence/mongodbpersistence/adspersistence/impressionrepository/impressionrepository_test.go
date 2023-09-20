package impressionrepository_test

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/impression"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/adspersistence/impressionrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var impRepo impression.ImpressionRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	impRepo = impressionrepository.NewMongoDBImpressionRepository(dbStore, logger)
}

func TestSaveImpression(t *testing.T) {
	i := shared.NewID()
	a := impression.NewImpression(i, i, i, i, i)

	err := impRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedimpression(t *testing.T) {
	i := shared.NewID()
	a := impression.NewImpression(i, i, i, i, i)

	err := impRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}

	fa, err := impRepo.Get(context.TODO(), a.ID)

	if fa.ID != a.ID || fa.AdPieceID != a.AdPieceID {
		t.FailNow()
	}
}

func TestImpressionWithintimeperiod(t *testing.T) {
	i := shared.NewID()
	a := impression.NewImpression(i, i, i, i, i)

	err := impRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}
	a2 := impression.NewImpression(i, i, i, i, i)

	err = impRepo.Save(context.TODO(), a2)
	if err != nil {
		t.Fatal(err)
	}
	s := time.Date(2022, 1, 1, 1, 1, 1, 1, time.Local)
	adc, err := impRepo.ImpressionsForPublisher(i, s, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	if !slices.ContainsFunc(adc, func(ac impression.Impression) bool {
		return ac.AdPieceID == i
	}) {
		t.FailNow()
	}
}
