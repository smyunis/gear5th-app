package siteverification

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
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
		return site.NewSiteVerificationError("unable to access ads.txt from network", err)

	}

	if !hasAdsTxtRecord(adstxtBody, desiredRecord) {
		return site.NewSiteVerificationError("unable to find record in ads.txt", nil)
	}

	return nil
}

func (a *AdsTxtVerificationService) fetchAdsTxtContent(s *site.Site) (string, error) {
	adstxturl, err := a.adsTxtUrl(s)
	if err != nil {
		return "", fmt.Errorf("unable to form ads.txt url: %w", err)
	}

	response, err := a.httpClient.Get(adstxturl)
	if err != nil {
		return "", fmt.Errorf("unable to fetch ads.txt from network: %w", err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("unable to parse ads.txt: %w", err)
	}

	return string(body), nil
}

func (*AdsTxtVerificationService) adsTxtUrl(s *site.Site) (string, error) {
	baseUrl := s.Url().Scheme + "://" + s.Url().Host
	adsTxturl, err := url.JoinPath(baseUrl, "ads.txt")
	return adsTxturl, err
}

func hasAdsTxtRecord(body string, record AdsTxtRecord) bool {

	// Regex to match an ads.txt record
	pattern := `^(?i)\s*(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)+([A-Za-z]|[A-Za-z][A-Za-z0-9\-]*[A-Za-z0-9])\s*,(\s*.*\s*),(\s*(DIRECT|RESELLER)\s*),?(\s*[^,\n]*\s*)$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	body = strings.TrimSpace(body)
	lines := strings.Split(body, "\n")

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
		if reflect.DeepEqual(r, record) {
			return true
		}
	}
	return false
}
