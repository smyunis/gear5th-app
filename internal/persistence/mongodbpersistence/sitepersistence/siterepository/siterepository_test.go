//go:build db
package siterepository_test

import (
	"context"
	"net/url"
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/sitepersistence/siterepository"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var siteRepo site.SiteRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	siteRepo = siterepository.NewMongoDBSiteRepository(dbStore,logger)
}

func TestCanSaveSite(t *testing.T) {
	u, _ := url.Parse("https://www.google.com")
	s := site.NewSite(shared.NewID(), *u)

	err := siteRepo.Save(context.TODO(), s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCanGetSavedSite(t *testing.T) {
	u, _ := url.Parse("https://www.google.com")
	s := site.NewSite(shared.NewID(), *u)
	err := siteRepo.Save(context.TODO(), s)
	if err != nil {
		t.Fatal(err)
	}

	fetchedSite, err := siteRepo.Get(context.TODO(), s.SiteID())
	if err != nil {
		t.FailNow()
	}

	if fetchedSite.SiteID() != s.SiteID() {
		t.FailNow()
	}

}

func TestGetActiveSites(t *testing.T) {
	publisherID := shared.NewID()
	u, _ := url.Parse("https://www.google.com")
	s1 := site.NewSite(publisherID, *u)
	s2 := site.NewSite(publisherID, *u)
	s3 := site.NewSite(publisherID, *u)
	s3.Deactivate()
	
	err := siteRepo.Save(context.TODO(), s1)
	if err != nil {
		t.Fatal(err)
	}
	err = siteRepo.Save(context.TODO(), s2)
	if err != nil {
		t.Fatal(err)
	}
	err = siteRepo.Save(context.TODO(), s3)
	if err != nil {
		t.Fatal(err)
	}

	sites,err := siteRepo.ActiveSitesForPublisher(publisherID)
	if err != nil {
		t.Fatal(err)
	}

	if len(sites) != 2 || sites[0].PublisherId() != sites[1].PublisherId()  {
		t.FailNow()
	}
	
}
