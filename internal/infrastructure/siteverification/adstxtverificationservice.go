package siteverification

// go get github.com/tzafrirben/go-adstxt-crawler/adstxt

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/site"
)

type AdsTxtRecord struct {
	AdExchangeDomain string
	PublisherId      string
	Relation         string
	CertAuthTag      string
}

type AdsTxtVerificationService struct {
	httpClient http.Client
}

func NewAdsTxtVerificationService(httpClient http.Client) AdsTxtVerificationService {
	return AdsTxtVerificationService{
		httpClient: httpClient,
	}
}

func (a *AdsTxtVerificationService) VerifyAdsTxt(s *site.Site, desiredRecord AdsTxtRecord) error {

	adstxtBody, err := a.fetchAdsTxtContent(s)
	if err != nil {
		return err
	}

	if a.hasAdsTxtRecord(adstxtBody, desiredRecord) {
		return nil
	}

	return site.SiteVerificationError{
		Reason: "unable to find record in ads.txt",
		Inner:  err,
	}
}

func (*AdsTxtVerificationService) hasAdsTxtRecord(body string, record AdsTxtRecord) bool {

	body = strings.TrimSpace(body)
	lines := strings.Split(body, "\n")

	pattern := fmt.Sprintf(`^%s,\s.*,\sDIRECT,?[^,]*$`, record.AdExchangeDomain)
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	matches := make([]string, 0)
	for _, line := range lines {
		if regex.MatchString(line) {
			matches = append(matches, line)
		}
	}

	records := make([]AdsTxtRecord, 0)
	for _, rec := range matches {
		fields := strings.Split(string(rec), ",")
		adsTxtRecord := AdsTxtRecord{}
		if len(fields) == 3 || len(fields) == 4 {
			adsTxtRecord.AdExchangeDomain = strings.TrimSpace(fields[0])
			adsTxtRecord.PublisherId = strings.TrimSpace(fields[1])
			adsTxtRecord.Relation = strings.TrimSpace(fields[2])

			if len(fields) == 4 {
				adsTxtRecord.CertAuthTag = strings.TrimSpace(fields[3])
			}
			records = append(records, adsTxtRecord)
		}
	}

	for _, r := range records {
		if r.AdExchangeDomain == record.AdExchangeDomain &&
			r.PublisherId == record.PublisherId {
			return true
		}
	}
	return false

}

func (a *AdsTxtVerificationService) fetchAdsTxtContent(s *site.Site) (string, error) {
	adsTxtUrl, err := a.AdsTxtUrl(s)
	if err != nil {
		return "", site.SiteVerificationError{
			Reason: "unable to form ads.txt url",
			Inner:  err,
		}
	}

	response, err := a.httpClient.Get(adsTxtUrl)
	if err != nil {
		return "", site.SiteVerificationError{
			Reason: "unable to fetch ads.txt from network",
			Inner:  err,
		}
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", site.SiteVerificationError{
			Reason: "unable to parse ads.txt body",
			Inner:  err,
		}
	}

	return string(body), nil
}

func (*AdsTxtVerificationService) AdsTxtUrl(s *site.Site) (string, error) {
	baseUrl := s.Url().Scheme + "://" + s.Url().Host
	adsTxturl, err := url.JoinPath(baseUrl, "ads.txt")
	return adsTxturl, err
}


