package authorization

import (
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
)

func CanRemoveSite(actor user.User, s site.Site) bool {
	return actor.HasRole(user.Publisher) && s.PublisherId() == actor.UserID()
}

func CanCreateSite(actor user.User) bool {
	return actor.HasRole(user.Publisher)
}
