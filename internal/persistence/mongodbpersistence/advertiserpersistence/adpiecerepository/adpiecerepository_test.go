//go:build db
package adpiecerepository_test

import (
	"context"
	"net/url"
	"os"
	"reflect"
	"slices"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/advertiserpersistence/adpiecerepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var adPieceRepo adpiece.AdPieceRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	adPieceRepo = adpiecerepository.NewMongoDBAdPieceRepository(dbStore, logger)
}

func TestSaveAdPiece(t *testing.T) {
	r, _ := url.Parse("https://www.youtube.com/watch?v=GpE37iwg-qI")
	a := adpiece.NewAdPiece(shared.NewID(), adslot.Box, r, "res-11", "image/jpeg")

	err := adPieceRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedAdPiece(t *testing.T) {
	r, _ := url.Parse("https://www.youtube.com/watch?v=GpE37iwg-qI")
	a := adpiece.NewAdPiece(shared.NewID(), adslot.Box, r, "res-11", "image/jpeg")

	err := adPieceRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}

	fa, err := adPieceRepo.Get(context.TODO(), a.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fa, a) {
		t.FailNow()
	}
}

func TestActiveAdpiecesForCampaign(t *testing.T) {
	campaignID := shared.NewID()
	r, _ := url.Parse("https://www.youtube.com/watch?v=GpE37iwg-qI")
	a := adpiece.NewAdPiece(campaignID, adslot.Box, r, "res-11", "image/jpeg")

	err := adPieceRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}
	r2, _ := url.Parse("https://www.youtube.com/watch?v=eNlRqyUV4ss")
	a2 := adpiece.NewAdPiece(campaignID, adslot.Vertical, r2, "res-22", "image/jpeg")
	a2.Deactivate()
	err = adPieceRepo.Save(context.TODO(), a2)
	if err != nil {
		t.Fatal(err)
	}

	rc,err := adPieceRepo.ActiveAdPiecesForCampaign(campaignID)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.ContainsFunc(rc, func(a adpiece.AdPiece) bool {
		return a.Resource == "res-11"
	}) {
		t.FailNow()
	}

	if slices.ContainsFunc(rc, func(a adpiece.AdPiece) bool {
		return a.Resource == "res-22"
	}) {
		t.FailNow()
	}
}
