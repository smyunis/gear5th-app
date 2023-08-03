package site_test

import (
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


