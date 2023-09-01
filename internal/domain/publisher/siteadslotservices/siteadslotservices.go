package siteadslotservices

import (
	"html/template"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
)

var adSlotHTMLTemplate *template.Template

type htmlSippetPresenter struct {
	UserID   string
	SiteID   string
	AdSlotID string
	Width    int
	Height   int
}

func init() {
	//TODO generate html
	tmpl := `<iframe width="{{.Width}}" height="{{.Height}}">{{.}}</iframe>`
	adSlotHTMLTemplate = template.Must(template.New("adslot-integration-snippet").Parse(tmpl))
}

func GenerateIntegrationHTMLSnippet(s site.Site, slot adslot.AdSlot) (string, error) {
	var htmlStringBuilder strings.Builder
	p := htmlSippetPresenter{
		UserID:   s.PublisherId().String(),
		SiteID:   s.ID().String(),
		AdSlotID: slot.ID().String(),
		Width: slot.AdSlotType().Dimentions().Width,
		Height: slot.AdSlotType().Dimentions().Height,
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
