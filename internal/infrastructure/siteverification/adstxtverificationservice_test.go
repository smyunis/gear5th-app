package siteverification_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/siteverification"
)

func TestVerifyAdsTxt(t *testing.T) {
	publisherId := "9929"
	su, _ := url.Parse("https://www.wikihow.com/")
	s := site.NewSite(shared.Id(publisherId), *su)

	service := siteverification.NewAdsTxtVerificationService(http.Client{
		Timeout: time.Minute * 1,
	})

    // appnexus.com, 9929, DIRECT, f5ab79cb980f11d1
	record := siteverification.AdsTxtRecord{
		AdExchangeDomain: "appnexus.com",
		PublisherId: "9929",
		Relation: "DIRECT",
		CertAuthTag: "f5ab79cb980f11d1",
	}

	err := service.VerifyAdsTxt(&s,record)

	if err != nil {
		t.Fatal(err.Error())
	}

}
