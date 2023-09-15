package siteadslotservices

import (
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
)

type AdSlotHTMLSnippetService interface {
	GenerateHTML(s site.Site, slot adslot.AdSlot) (string, error)
}

func CanServeAdPiece(s site.Site, slot adslot.AdSlot) bool {
	return s.CanServeAdPiece() && !slot.IsDeactivated()
}
