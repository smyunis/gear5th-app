package site

import (
	"net/url"
)

type SiteVerificationError struct {
	message string
	inner   error
}

func NewSiteVerificationError(message string, inner error) error {
	return SiteVerificationError{
		message: message,
		inner:   inner,
	}
}

func (e SiteVerificationError) Error() string {
	return e.message
}

func (e SiteVerificationError) Unwrap() error {
	return e.inner
}

type AdsTxtRecord struct {
	AdExchangeDomain string
	PublisherId      string
	Relation         string
	CertAuthTag      string
}

type AdsTxtVerificationService interface {
	VerifyAdsTxt(site *Site, record AdsTxtRecord) error
}

func VerifySiteHostname(site *Site, source url.URL) error {
	siteFqdn := site.url.Hostname()
	sourceFqdn := source.Hostname()

	if siteFqdn != sourceFqdn {
		return NewSiteVerificationError("site hostname mismatch", nil)
	}

	site.Verify()
	return nil
}
