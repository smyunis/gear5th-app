package authorization_test

import (
	"net/url"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/authorization"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestCanNotRemoveSiteThePublisherDoesNotOwn(t *testing.T) {
	pubID := shared.NewID()
	siteURL, _ := url.Parse("https://www.google.com")
	s := site.NewSite(pubID, *siteURL)

	email, _ := user.NewEmail("x@y.com")
	u := user.NewUser(email)

	if authorization.CanModifySite(u, s) {
		t.FailNow()
	}
}

func TestUserCanCreateSite(t *testing.T) {
	email, _ := user.NewEmail("x@y.com")
	u := user.NewUser(email)
	u.SignUpPublisher()

	if !authorization.CanCreateSite(u) {
		t.FailNow()
	}
}

func TestUserCanModifyAdSlotOnlyIfUserOwnsSiteItBelongsTo(t *testing.T) {
	email, _ := user.NewEmail("x@y.com")
	u := user.NewUser(email)
	u.SignUpPublisher()
	siteURL, _ := url.Parse("https://www.google.com")
	s := site.NewSite(u.UserID(), *siteURL)
	slot := adslot.NewAdSlot(s.ID(), "adslot-x", adslot.Horizontal)

	if !authorization.CanModifyAdSlot(u, s, slot) {
		t.FailNow()
	}
}
