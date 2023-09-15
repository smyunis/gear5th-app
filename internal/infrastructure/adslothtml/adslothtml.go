package adslothtml

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/siteadslotservices"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

var adSlotHTMLTemplate *template.Template

func init() {
	//TODO generate html
	tmpl := `<iframe src="{{.AdServerURL}}" width="{{.AdSlot.AdSlotType.Dimentions.Width}}" 
		height="{{.AdSlot.AdSlotType.Dimentions.Height}}">{{.}}</iframe>`
	adSlotHTMLTemplate = template.Must(template.New("adslot-integration-snippet").Parse(tmpl))
}

type htmlSippetPresenter struct {
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
	var htmlStringBuilder strings.Builder

	adServerURL := a.appURL.JoinPath(fmt.Sprintf("/ads/adserver?slot=%s", strings.ToLower(slot.AdSlotType().String())))

	p := htmlSippetPresenter{
		Site:        s,
		AdSlot:      slot,
		AdServerURL: adServerURL,
	}
	err := adSlotHTMLTemplate.ExecuteTemplate(&htmlStringBuilder, "adslot-integration-snippet", p)
	if err != nil {
		return "", err
	}
	return htmlStringBuilder.String(), nil
}

var x siteadslotservices.AdSlotHTMLSnippetService = AdSlotHTMLSnippetService{}
