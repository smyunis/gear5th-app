/// go:embed db

package campaignrepository_test

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/advertiserpersistence/campaignrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var campaignRepo campaign.CampaignRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	campaignRepo = campaignrepository.NewMongoDBCampaignRepository(dbStore, logger)
}

func TestCanSaveCampaign(t *testing.T) {
	c := campaign.NewCampaign("coca-cola summer",
		shared.NewID(), time.Now(), time.Now().Add(6*time.Hour))

	err := campaignRepo.Save(context.TODO(), c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCanGetSavedCampaign(t *testing.T) {
	c := campaign.NewCampaign("coca-cola summer",
		shared.NewID(), time.Now(), time.Now().Add(6*time.Hour))

	err := campaignRepo.Save(context.TODO(), c)
	if err != nil {
		t.Fatal(err)
	}

	fc, err := campaignRepo.Get(context.TODO(), c.ID)
	if err != nil {
		t.Fatal(err)
	}

	if fc.ID != c.ID && fc.AdvertiserUserID != c.AdvertiserUserID {
		t.FailNow()
	}
}

func TestCanGetCampaignsForAdvertiser(t *testing.T) {
	advertiserID := shared.NewID()
	c1 := campaign.NewCampaign("coca-cola summer",
		advertiserID, time.Now(), time.Now().Add(6*time.Hour))

	err := campaignRepo.Save(context.TODO(), c1)
	if err != nil {
		t.Fatal(err)
	}
	c2 := campaign.NewCampaign("pepsi summer",
		advertiserID, time.Now(), time.Now().Add(7*time.Hour))

	err = campaignRepo.Save(context.TODO(), c2)
	if err != nil {
		t.Fatal(err)
	}

	c3 := campaign.NewCampaign("fanta summer",
		advertiserID, time.Now(), time.Now().Add(4*time.Hour))
	c3.Quit()
	err = campaignRepo.Save(context.TODO(), c3)
	if err != nil {
		t.Fatal(err)
	}

	camps, err := campaignRepo.CampaignsForAdvertiser(advertiserID)
	if err != nil {
		t.Fatal(err)
	}

	if len(camps) != 2 {
		t.FailNow()
	}

}

func TestRunningCampaings(t *testing.T) {
	advertiserID := shared.NewID()
	c1 := campaign.NewCampaign("bobo king chips",
		advertiserID,
		time.Date(2012, 1, 1, 1, 1, 1, 1, time.UTC),
		time.Date(2016, 1, 1, 1, 1, 1, 1, time.UTC))

	err := campaignRepo.Save(context.TODO(), c1)
	if err != nil {
		t.Fatal(err)
	}

	c2 := campaign.NewCampaign("sun li drinks",
		advertiserID, time.Now(), time.Now().Add(7*time.Hour))

	err = campaignRepo.Save(context.TODO(), c2)
	if err != nil {
		t.Fatal(err)
	}

	rc, err := campaignRepo.RunningCampaigns()
	if err != nil {
		t.FailNow()
	}

	if slices.ContainsFunc(rc, func(c campaign.Campaign) bool {
		return c.Name == "bobo king chips"
	}) {
		t.FailNow()
	}

	if !slices.ContainsFunc(rc, func(c campaign.Campaign) bool {
		return c.Name == "sun li drinks"
	}) {
		t.FailNow()
	}
}
