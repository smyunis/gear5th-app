package site_test

import (
	"fmt"
	"net/url"
	"testing"

	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

func TestSiteHostnameVerification(t *testing.T) {
	siteUrl, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")
	s := site.NewSite(shared.NewId(), *siteUrl)
	source, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")

	err := site.VerifySiteHostname(&s, *source)

	if err != nil {
		t.FailNow()
	}
}

func TestGetAdsTxtRecord(t *testing.T) {
	siteUrl, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")
	publisherId := shared.NewId()
	s := site.NewSite(publisherId, *siteUrl)
	expectedRecord := fmt.Sprintf("gear5th.com, %s, DIRECT", publisherId.String())

	record := site.GetAdsTxtRecord(s)

	if expectedRecord != record {
		t.Fatalf("expected: %s, given: %s", expectedRecord, record)
	}
}
