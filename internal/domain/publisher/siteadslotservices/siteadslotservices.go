package siteadslotservices

import (
	"html/template"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
)

var adSlotHTMLTemplate *template.Template

type htmlSippetPresenter struct {
	userID   string
	siteID   string
	adSlotID string
}

func init() {
	//TODO generate html
	tmpl := `<iframe>{{.}}</iframe>`
	adSlotHTMLTemplate = template.Must(template.New("adslot-integration-snippet").Parse(tmpl))
}

func GenerateIntegrationHTMLSnippet(s site.Site, slot adslot.AdSlot) (string, error) {
	var htmlStringBuilder strings.Builder
	p := htmlSippetPresenter{
		userID:   s.PublisherId().String(),
		siteID:   s.SiteID().String(),
		adSlotID: slot.AdSlotID().String(),
	}
	err := adSlotHTMLTemplate.ExecuteTemplate(&htmlStringBuilder, "adslot-integration-snippet", p)
	if err != nil {
		return "", err
	}
	return htmlStringBuilder.String(), nil
}


func CanServeAdPiece(s site.Site, slot adslot.AdSlot) bool {
	return s.CanServeAdPiece() && !slot.IsDeactivated()
}