package authorization

import (
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
)

func CanModifySite(actor user.User, s site.Site) bool {
	return actor.HasRole(user.Publisher) && s.PublisherID == actor.ID
}

func CanCreateSite(actor user.User) bool {
	return actor.HasRole(user.Publisher)
}

func CanModifyAdSlot(actor user.User, s site.Site, slot adslot.AdSlot) bool {
	return actor.HasRole(user.Publisher) && (s.ID == slot.SiteID) &&
		(s.PublisherID == actor.ID)
}
