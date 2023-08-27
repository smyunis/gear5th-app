package siteverification

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

type AdsTxtVerificationService struct {
	httpClient infrastructure.HTTPClient
	logger     application.Logger
}

func NewAdsTxtVerificationService(httpClient infrastructure.HTTPClient,
	logger application.Logger) AdsTxtVerificationService {
	return AdsTxtVerificationService{
		httpClient,
		logger,
	}
}

func (a *AdsTxtVerificationService) VerifyAdsTxt(s *site.Site, desiredRecord site.AdsTxtRecord) error {

	adstxtBody, err := a.fetchAdsTxtContent(s)
	if err != nil {
		return site.ErrSiteVerification

	}

	if !hasAdsTxtRecord(adstxtBody, desiredRecord) {
		return site.ErrSiteVerification
	}

	s.Verify()

	return nil
}

func (a *AdsTxtVerificationService) fetchAdsTxtContent(s *site.Site) (string, error) {
	adstxturl, err := a.adsTxtUrl(s)
	if err != nil {
		return "", fmt.Errorf("unable to form ads.txt url: %w", err)
	}

	request, err := http.NewRequest(http.MethodGet, adstxturl, nil)
	if err != nil {
		a.logger.Error("site/adstxtverification", err)
		return "", fmt.Errorf("unable to fetch ads.txt from network: %w", err)
	}

	response, err := a.httpClient.Do(request)
	if err != nil {
		a.logger.Error("site/adstxtverification", err)
		return "", fmt.Errorf("unable to fetch ads.txt from network: %w", err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		a.logger.Error("site/adstxtverification", err)
		return "", fmt.Errorf("unable to parse ads.txt: %w", err)
	}

	return string(body), nil
}

func (*AdsTxtVerificationService) adsTxtUrl(s *site.Site) (string, error) {
	baseUrl := s.URL().Scheme + "://" + s.URL().Host
	adsTxturl, err := url.JoinPath(baseUrl, "ads.txt")
	return adsTxturl, err
}

func hasAdsTxtRecord(body string, record site.AdsTxtRecord) bool {

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

	records := make([]site.AdsTxtRecord, 0)
	for _, rec := range matches {
		fields := strings.Split(string(rec), ",")
		adsTxtRecord := site.AdsTxtRecord{}
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
