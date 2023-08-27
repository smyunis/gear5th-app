package siteadslotservices_test

import (
	"net/url"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/siteadslotservices"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestCanGenerateHTML(t *testing.T) {
	s := newSite()
	slot := s.AddAdSlot("adslotx", adslot.Box)
	_, err := siteadslotservices.GenerateIntegrationHTMLSnippet(s, slot)
	if err != nil {
		t.FailNow()
	}
}

func TestCanServeAdPiece(t *testing.T) {
	s := newSite()
	s.Verify()
	slot := s.AddAdSlot("adslotx", adslot.Box)
	if !siteadslotservices.CanServeAdPiece(s, slot) {
		t.FailNow()
	}
}

func TestCanNotServeAdPieceIfDeactivatedSite(t *testing.T) {
	s := newSite()
	s.Deactivate()
	slot := s.AddAdSlot("adslotx", adslot.Box)
	if siteadslotservices.CanServeAdPiece(s, slot) {
		t.FailNow()
	}
}

func TestCanNotServeAdPieceIfDeactivatedAdSlot(t *testing.T) {
	s := newSite()
	slot := s.AddAdSlot("adslotx", adslot.Box)
	slot.Deactivate()
	if siteadslotservices.CanServeAdPiece(s, slot) {
		t.FailNow()
	}
}


func newSite() site.Site {
	siteUrl, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")
	publisherId := shared.NewID()
	return site.NewSite(publisherId, *siteUrl)
}
