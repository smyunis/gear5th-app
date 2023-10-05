package adslothtml

import (
	_ "embed"
	"html/template"
	"net/url"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

var adSlotHTMLTemplate *template.Template

//go:embed adslot.html
var adSlotHTMLTemplateFile string

func init() {
	adSlotHTMLTemplate = template.Must(template.New("adslot-integration-snippet").Parse(adSlotHTMLTemplateFile))
}

type htmlSnippetPresenter struct {
	Site        site.Site
	AdSlot      adslot.AdSlot
	AdServerURL *url.URL
}

type AdSlotHTMLSnippetService struct {
	config infrastructure.ConfigurationProvider
	appURL *url.URL
}

func NewAdSlotHTMLSnippetService(config infrastructure.ConfigurationProvider) AdSlotHTMLSnippetService {

	appurlstr := config.Get("APP_URL", "https://gear5th.com")
	a, err := url.Parse(appurlstr)
	if err != nil {
		panic(err.Error())
	}

	return AdSlotHTMLSnippetService{
		config,
		a,
	}
}

func (a AdSlotHTMLSnippetService) GenerateHTML(s site.Site, slot adslot.AdSlot) (string, error) {

	adServerURL := a.appURL.JoinPath("/ads/adserver")
	q := adServerURL.Query()
	q.Add("slot", strings.ToLower(slot.SlotType.String()))
	q.Add("publisher-id", s.PublisherID.String())
	q.Add("adslot-id", slot.ID.String())
	q.Add("site-id", s.ID.String())
	adServerURL.RawQuery = q.Encode()

	p := htmlSnippetPresenter{
		Site:        s,
		AdSlot:      slot,
		AdServerURL: adServerURL,
	}
	
	var htmlStringBuilder strings.Builder
	err := adSlotHTMLTemplate.ExecuteTemplate(&htmlStringBuilder, "adslot-integration-snippet", p)
	if err != nil {
		return "", err
	}
	return htmlStringBuilder.String(), nil
}
